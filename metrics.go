// package metrics wraps github.com/DataDog/datadog-go/v5/statsd with a more ergonomic and
// computationally cheaper interface.
//
// This is done by separating tags from logging metrics so that for frequently-logged gauges and
// counters logging is just a single atomic operation.
//
// For each metric type of Gauge, Count, Distribution, and Set, there are a set of NewMDefY methods
// where M is the metric type and Y is the number of tags. Calls to NewMDefY must be done at
// init-time (ideally in a top-level var block) of a metrics.go file with names as full literals so
// that metrics are easily greppable. Metrics not defined this way will cause the process to panic
// if still at init-time, meaning before any code in main() has run, otherwise will produce
// non-functional stats and produce to a gauge stat called metrics.bad_metric_definitions.  It's a
// good idea to put an alert on this stat so that if it starts logging during a deploy, you know
// your other metrics may not be trustworthy.
//
// See the example folder for an example of usage.
//
// Generally you will have:
//  1. metrics.NewMDefY calls in a top-level var block of metrics.go in packages that log metrics,
//     named ending in "Def".
//  2. metrics.New() in main(), passed down through constructors.
//  3. At the point of logging metrics, e.g.
//     m.Counter(myCounterDef.Values(tag1, tag2)).Add(1)
//
// # bad_metric_definitions reasons
//
//   - not_at_init_time: The call to NewMDefY did not happen at init time (in either a top-level var
//     block or func init()).
//   - runtime_caller_failed: [runtime.Caller] returned false in its final return trying to evaluate
//     the above.
//   - observe_duration_bad_units: [Distribution.ObserveDuration] was used on a def that did not
//     have compatible units. See the comment on [Distribution.ObserveDuration].
package metrics

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bradenaw/juniper/xslices"
	"github.com/bradenaw/juniper/xsync"
	"golang.org/x/exp/maps"
)

const (
	// Datadog's flush interval is 10 seconds, so we need to use something that
	// is at least that fast. For counters, we need to use something that divides Datadog's flush
	// interval so that every flush has the same number of counter aggregates put into it.
	//
	// https://docs.datadoghq.com/developers/dogstatsd/?tab=hostagent
	flushInterval = 2 * time.Second
)

// Publisher is the subset of github.com/DataDog/datadog-go/v5/statsd.ClientInterface used by this
// package.
//
// Publisher should _not_ have client-side aggregation enabled because this package also does
// aggregation. It is enabled by default in datadog-go/v5, so should be disabled with
// statsd.WithoutClientSideAggregation().
type Publisher interface {
	Gauge(name string, value float64, tags []string, rate float64) error
	Count(name string, value int64, tags []string, rate float64) error
	Distribution(name string, value float64, tags []string, rate float64) error
	Set(name string, value string, tags []string, rate float64) error
}

// TagValue is the value of a key:value pair in a metric tag. They are formatted the same as
// fmt.Sprint unless the type implements TagValuer, in which case MetricTagValue() is used instead.
//
// TagValues that produce the same string are considered the same.
type TagValue any

// See the comment on type TagValue.
type TagValuer interface {
	MetricTagValue() string
}

type Metrics struct {
	sender   *newlineDelimPacketSender
	bg       *xsync.Group
	flushNow func()

	gauges        metricMap[*Gauge]
	counters      metricMap[*Counter]
	distributions metricMap[*Distribution]
	sets          metricMap[*Set]

	m       sync.Mutex
	flushed chan struct{}
	nextID  int
	polls   map[int]func()
}

type noOpPublisher struct{}

func (p noOpPublisher) Gauge(name string, value float64, tags []string, rate float64) error {
	return nil
}
func (p noOpPublisher) Count(name string, value int64, tags []string, rate float64) error {
	return nil
}
func (p noOpPublisher) Distribution(name string, value float64, tags []string, rate float64) error {
	return nil
}
func (p noOpPublisher) Set(name string, value string, tags []string, rate float64) error { return nil }

var (
	NoOpMetrics = &Metrics{
		bg:      xsync.NewGroup(context.Background()),
		flushed: make(chan struct{}),
		polls:   make(map[int]func()),
	}

	// Used to return from Metrics.Metric() methods when the definition is invalid and the stat
	// can't be logged.
	noOpCounter      = &Counter{m: NoOpMetrics}
	noOpGauge        = &Gauge{m: NoOpMetrics}
	noOpDistribution = &Distribution{m: NoOpMetrics}
	noOpSet          = &Set{m: NoOpMetrics}

	badDefsDef = NewGaugeDef1[string](
		"metrics.bad_metric_definitions",
		"The number of calls to NewMDefY that are invalid for some reason. These definitions "+
			"will not be able to log metrics at all.",
		UnitItem,
		[...]string{"reason"},
	)
)

func New() *Metrics {
	m := &Metrics{
		bg:      xsync.NewGroup(context.Background()),
		flushed: make(chan struct{}),
		polls:   make(map[int]func()),
	}

	badDefsCallersFramesGauge := m.Gauge(badDefsDef.Values("runtime_caller_failed"))
	badDefsNotAtInitGauge := m.Gauge(badDefsDef.Values("not_at_init_time"))
	badDefsObserveDurationBadUnitsGauge := m.Gauge(badDefsDef.Values("observe_duration_bad_units"))

	m.flushNow = m.bg.PeriodicOrTrigger(flushInterval, 0 /*jitter*/, func(ctx context.Context) {
		m.m.Lock()
		polls := maps.Values(m.polls)
		m.m.Unlock()
		for _, poll := range polls {
			poll()
		}

		badDefsCallersFramesGauge.Set(float64(badDefsCallersFrames.Load()))
		badDefsNotAtInitGauge.Set(float64(badDefsNotAtInit.Load()))
		badDefsObserveDurationBadUnitsGauge.Set(float64(badObserveDurations.Load()))

		m.gauges.Range(func(_ metricKey, g *Gauge) bool {
			g.publish()
			return true
		})
		m.counters.Range(func(_ metricKey, c *Counter) bool {
			c.publish()
			return true
		})
		m.distributions.Range(func(_ metricKey, d *Distribution) bool {
			d.publish()
			return true
		})

		m.m.Lock()
		close(m.flushed)
		m.flushed = make(chan struct{})
		m.m.Unlock()
	})

	return m
}

// Counter returns the Counter for the given CounterDef. For the same CounterDef, including one
// produced from CounterDefY.Values() with the same values, this will return the same *Counter.
//
// For metrics with tags (e.g. CounterDef2), the CounterDef can be made by calling Values(), for
// example:
//
//	// ---- metrics.go -----------------------------------------------------------------------------
//	rpcResponseDef = metrics.NewCounterDef2[string, string](
//		"rpc_responses",
//		"Counts responses to each RPC by method and status.",
//		[...]string{"method", "status"},
//		metrics.UnitResponse,
//	)
//
//	// ---- at the point of logging the metric -----------------------------------------------------
//	m.Counter(rpcResponseDef.Values(methodName, status)).Add(1)
//
// Metrics.Counter is relatively expensive relative to Counter.Add, so very high-throughput logging
// should cache the result of this function:
//
//	// ---- at creation of RPC server --------------------------------------------------------------
//	// rpc_responses method:get status:ok
//	s.getOKCounter = m.Counter(rpcResponseDef.Values("get", "ok"))
//
//	// rpc_responses method:get status:error
//	s.getErrorCounter = m.Counter(rpcResponseDef.Values("get", "error"))
//
//	// ---- inside the Get() RPC handler -----------------------------------------------------------
//	if err == nil {
//		s.getOKCounter.Add(1)
//	} else {
//		s.getErrorCounter.Add(1)
//	}
func (m *Metrics) Counter(d CounterDef) *Counter {
	if !d.ok {
		return noOpCounter
	}

	k := newMetricKey(d.name, d.tags.n, d.tags.values, d.allComparable)
	c, ok := m.counters.Load(k)
	if !ok {
		c = &Counter{
			m:    m,
			name: d.name,
			tags: makeTags(d.tags.keys[:d.tags.n], d.tags.values[:d.tags.n]),
		}
		c, _ = m.counters.LoadOrStore(k, c)
	}
	return c
}

// Gauge returns the Gauge for the given GaugeDef. For the same GaugeDef, including one produced
// from GaugeDefY.Values() with the same values, this will return the same *Gauge.
func (m *Metrics) Gauge(d GaugeDef) *Gauge {
	if !d.ok {
		return noOpGauge
	}

	k := newMetricKey(d.name, d.tags.n, d.tags.values, d.allComparable)
	g, ok := m.gauges.Load(k)
	if !ok {
		g = &Gauge{
			m:    m,
			name: d.name,
			tags: makeTags(d.tags.keys[:d.tags.n], d.tags.values[:d.tags.n]),
		}
		g.v.Store(math.Float64bits(math.NaN()))
		g, _ = m.gauges.LoadOrStore(k, g)
	}
	return g
}

// Distribution returns the Distribution for the given DistributionDef. For the same
// DistributionDef, including one produced from DistributionDefY.Values() with the same values, this
// will return the same *Distribution.
func (m *Metrics) Distribution(d DistributionDef) *Distribution {
	if !d.ok {
		return noOpDistribution
	}

	k := newMetricKey(d.name, d.tags.n, d.tags.values, d.allComparable)
	c, ok := m.distributions.Load(k)
	if !ok {
		c = &Distribution{
			m:    m,
			name: d.name,
			unit: d.unit,
			tags: makeTags(d.tags.keys[:d.tags.n], d.tags.values[:d.tags.n]),
		}
		c, _ = m.distributions.LoadOrStore(k, c)
	}
	return c
}

// Set measures the cardinality of values passed to Observe for each time bucket, that is, it
// estimates how many _unique_ values have been passed to it.
type Set struct {
	m          *Metrics
	name       string
	tags       string
	sampleRate float64
}

func (s *Set) Observe(value string) {
	panic("todo")
}

// Set returns the Set for the given SetDef. For the same SetDef, including one produced from
// SetDefY.Values() with the same values, this will return the same *Set.
func (m *Metrics) Set(d SetDef) *Set {
	if !d.ok {
		return noOpSet
	}

	k := newMetricKey(d.name, d.tags.n, d.tags.values, d.allComparable)
	c, ok := m.sets.Load(k)
	if !ok {
		c = &Set{
			m:          m,
			name:       d.name,
			tags:       makeTags(d.tags.keys[:d.tags.n], d.tags.values[:d.tags.n]),
			sampleRate: d.sampleRate,
		}
		c, _ = m.sets.LoadOrStore(k, c)
	}
	return c
}

// EveryFlush calls f once before each aggregate metric flush. This is useful for e.g. gauges that
// need to be periodically computed.
//
// f happens on the same goroutine that flushes metrics, so it should not be too expensive or it can
// interfere with metrics being sent.
func (m *Metrics) EveryFlush(f func()) func() {
	m.m.Lock()
	defer m.m.Unlock()

	// Held during f so that the returned stop-reporting can be sure that f will not be in progress
	// nor called again once it returns.
	var fm sync.Mutex
	done := false

	id := m.nextID
	m.nextID++
	m.polls[id] = func() {
		fm.Lock()
		defer fm.Unlock()
		if done {
			return
		}
		f()
	}

	return func() {
		fm.Lock()
		defer fm.Unlock()
		m.m.Lock()
		defer m.m.Unlock()
		done = true
		delete(m.polls, id)
	}
}

// Flush immediately sends pending metric data to the Publisher given to m in New() and blocks
// until complete.
func (m *Metrics) Flush() {
	if m.flushNow == nil {
		return
	}
	m.m.Lock()
	flushed := m.flushed
	m.m.Unlock()
	m.flushNow()
	<-flushed
}

// Close frees resources associated with m. After Close, m should not be used.
func (m *Metrics) Close() {
	m.bg.StopAndWait()
}

// Gauge is a metric that reports the last value that it was set to.
//
// Unlike Datadog's native gauge in the statsd client, Gauges report this value until the end of the
// process or until explicitly Unset().
//
// Gauges are good for measuring states, for example the number of open connections or the size of a
// buffer.
type Gauge struct {
	m    *Metrics
	name string
	tags string
	v    atomic.Uint64
}

// Set sets the value of the gauge. The gauge will continue to have this value until the next Set or
// Unset, or the end of the process.
func (g *Gauge) Set(v float64) {
	g.v.Store(math.Float64bits(v))
}

// Unset unsets the value of the gauge. If the Gauge remains unset, it will have no value for time
// buckets after this.
func (g *Gauge) Unset() {
	g.v.Store(math.Float64bits(math.NaN()))
}

func (g *Gauge) value() float64 {
	return math.Float64frombits(g.v.Load())
}

func (g *Gauge) publish() {
	v := g.value()
	if math.IsNaN(v) {
		return
	}

	g.m.sender.Write([]byte(g.name))
	g.m.sender.Write([]byte(":"))
	var b [64]byte
	g.m.sender.Write(strconv.AppendFloat(b[:0], v, 'f', -1, 64))
	g.m.sender.Write([]byte("|g|@1|#"))
	g.m.sender.Write([]byte(g.tags))
	g.m.sender.WriteNewline()
}

// Counter is a metric that keeps track of the number of events that happen per time interval.
//
// Counters are good for measuring the rate of events, for example requests per second, or measuring
// the ratio between events by using tags, such as error rate.
type Counter struct {
	m    *Metrics
	name string
	tags string
	v    atomic.Int64
}

func (c *Counter) Add(n int64) {
	c.v.Add(n)
}

func (c *Counter) publish() {
	v := c.v.Swap(0)
	if v == 0 {
		return
	}

	c.m.sender.Write([]byte(c.name))
	c.m.sender.Write([]byte(":"))
	var b [20]byte // 20 digits fits 2^64
	c.m.sender.Write(strconv.AppendInt(b[:0], v, 10 /*base*/))
	c.m.sender.Write([]byte("|c|@1|#"))
	c.m.sender.Write([]byte(c.tags))
	c.m.sender.WriteNewline()
}

// Distribution produces quantile metrics, e.g. 50th, 90th, 99th percentiles of the values passed to
// Observe for each time bucket.
type Distribution struct {
	m    *Metrics
	name string
	unit Unit
	tags string

	curr concurrentSketch
	prev sketch
}

func (d *Distribution) Observe(value float64) {
	d.curr.Observe(value)
}

func (d *Distribution) publish() {
	// https://docs.datadoghq.com/developers/dogstatsd/datagram_shell?tab=metrics
	//
	// <METRIC_NAME>:<VALUE1>:<VALUE2>:<VALUE3>|<TYPE>|@<SAMPLE_RATE>|#<TAG_KEY_1>:<TAG_VALUE_1>,<TAG_2>
	//
	// 'd' is the type for distribution

	// Formatted with 'f',-1, values are all about twenty bytes long which means we can only fit a
	// few packed values together.
	maxPerLine := (1400 - len(d.name) - len("|d|@1|#") - len(d.tags) - 1) / 20
	valuesThisLine := 0
	d.prev.newObservationsSince(&d.curr, func(value float64, count int) bool {
		if count < 10 {
			// TODO: uh, what if maxPerLine is super tiny because of a ton of tags? just let the
			// packet tear?
			if valuesThisLine == 0 {
				d.writeStart()
			}
			if valuesThisLine+count > maxPerLine {
				d.writeEndNoSample()
				d.writeStart()
				valuesThisLine = 0
			}
			d.writeValues(value, count)
			valuesThisLine += count
		} else {
			if valuesThisLine > 0 {
				d.writeEndNoSample()
				valuesThisLine = 0
			}
			d.writeStart()
			d.writeValues(value, 1)
			d.writeEnd(1 / float64(count))
		}
		return true
	})
	if valuesThisLine > 0 {
		d.writeEndNoSample()
	}
}

func (d *Distribution) writeStart() {
	d.m.sender.Write([]byte(d.name))
}

func (d *Distribution) writeValues(value float64, count int) {
	var b [64]byte
	// TODO: since value is coming out of a distribution bucket, it's actually a pretty constrained
	// set and we could just have a lookup table for common values instead
	vBytes := strconv.AppendFloat(b[:0], value, 'f', -1, 64)
	for i := 0; i < count; i++ {
		d.m.sender.Write([]byte(":"))
		d.m.sender.Write(vBytes)
	}
}

func (d *Distribution) writeEndNoSample() {
	d.m.sender.Write([]byte("|d|@1|#"))
	d.m.sender.Write([]byte(d.tags))
	d.m.sender.WriteNewline()
}

func (d *Distribution) writeEnd(sampleRate float64) {
	d.m.sender.Write([]byte("|d|@"))
	var b [64]byte
	// TODO: also consider a lookup table for, say, the first thousand
	d.m.sender.Write(strconv.AppendFloat(b[:0], sampleRate, 'f', -1, 64))
	d.m.sender.Write([]byte("|#"))
	d.m.sender.Write([]byte(d.tags))
	d.m.sender.WriteNewline()
}

var (
	badObserveDurationsSet = xsync.Map[string, struct{}]{}
	badObserveDurations    atomic.Uint64
)

// As long as d's units are in nanoseconds, microseconds, milliseconds, seconds, minutes, or hours,
// records the given duration in the correct units.
//
// Days and weeks are not supported, because not all days are the same length - in parts of the
// world that observe daylight savings, one day of the year is 25 hours and another is 23. As such,
// a time.Duration is not enough information to know days nor weeks so those must be recorded
// differently.
//
// Other units will record nothing, but will emit a metrics.bad_metrics_definitions with
// reason:observe_duration_bad_units.
func (d *Distribution) ObserveDuration(value time.Duration) {
	switch d.unit {
	case UnitNanosecond:
		d.Observe(float64(value.Nanoseconds()))
	case UnitMicrosecond:
		d.Observe(value.Seconds() * 1_000_000)
	case UnitMillisecond:
		d.Observe(value.Seconds() * 1_000)
	case UnitSecond:
		d.Observe(value.Seconds())
	case UnitMinute:
		d.Observe(value.Seconds() / 60)
	case UnitHour:
		d.Observe(value.Seconds() / 3600)
	default:
		_, loaded := badObserveDurationsSet.LoadOrStore(d.name, struct{}{})
		if !loaded {
			badObserveDurations.Add(1)
		}
	}
}

// metricKey is used to dedupe metrics so that multiple calls on a def result in the same metric. It
// contains the name and tag values.
type metricKey struct {
	name   string
	values [maxTags]any
}

func newMetricKey(name string, n int, values [maxTags]any, allComparable bool) metricKey {
	if allComparable {
		// Fast path - avoid reflection when all of the tag values are comparable.
		return metricKey{
			name:   name,
			values: values,
		}
	}

	k := metricKey{name: name}

	for i := range values {
		if reflect.ValueOf(values[i]).Comparable() {
			k.values[i] = values[i]
		} else {
			k.values[i] = tagValueString(values[i])
		}
	}

	return k
}

func makeTags(keys []string, values []any) string {
	var sb strings.Builder
	for i := range keys {
		if i != 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(makeTag(keys[i], values[i]))
	}
	return sb.String()
}

func makeTag(key string, value any) string {
	if len(key) == 0 {
		return tagValueString(value)
	}
	return key + ":" + tagValueString(value)
}

var validTagCharacters = func() [256]bool {
	var b [256]bool

	// https://docs.datadoghq.com/getting_started/tagging/#define-tags
	//
	// Allows:
	// - Alphanumerics
	// - Underscores
	// - Minuses
	// - Colons
	// - Periods
	// - Slashes

	// Note that uppercase have to get lowercased so we'll consider them 'invalid' and let the below
	// stuff sort it out.

	for r := 'a'; r <= 'z'; r++ {
		b[int(r)] = true
	}
	for r := '0'; r <= '9'; r++ {
		b[int(r)] = true
	}
	b['-'] = true
	b[':'] = true
	b['.'] = true
	b['/'] = true

	return b
}()

// Tag constraints are here:
// https://docs.datadoghq.com/getting_started/tagging/#define-tags
//
// There's some peculiarity here because tags that don't match this definition are automatically
// converted via lowercasing and replacing invalid characters with underscores, but:
//
// - This happens in the ddagent, after we've already sent it over the network.
// - The protocol to the ddagent uses some of these unsupported characters as delimiters.
// - The datadog client does _not_ do this conversion nor escape invalid characters.
//
// Which means that very strange results can come out. For example, a comma in a tag value will get
// interpreted as two separate tags, |@ will be interpreted as a sampling rate, \n will be
// interpreted as another metric altogether, etc.
//
// So we'll repeat the conversion here to save grief.
func tagValueSanitize(s string) string {
	ok := true

	// Most strings should be valid already so let's make that cheap and require no allocs.
	for _, b := range []byte(s) {
		if !validTagCharacters[b] {
			ok = false
			break
		}
	}

	if ok {
		return s
	}

	var sb strings.Builder
	// This'll never produce a string longer than s because we replace upper with lowercase ASCII
	// (both always one byte) and replace other characters including multi-byte with _ which is also
	// one byte.
	//
	// We could do the math but the caller should really just be handing us better strings.
	sb.Grow(len(s))
	for _, r := range s {
		if int(r) < 256 && validTagCharacters[int(r)] {
			sb.WriteRune(r)
		} else if r >= 'A' && r <= 'Z' {
			sb.WriteRune(r + ('a' - 'A'))
		} else {
			sb.WriteRune('_')
		}
	}
	return sb.String()
}

func tagValueString(v any) string {
	switch v := v.(type) {
	case string:
		return tagValueSanitize(v)
	case bool:
		return strconv.FormatBool(v)
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(int64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(uint64(v), 10)
	case TagValuer:
		return tagValueSanitize(v.MetricTagValue())
	case fmt.Stringer:
		return tagValueSanitize(v.String())
	default:
		return tagValueSanitize(fmt.Sprint(v))
	}
}

type MetricType string

const (
	CounterType      MetricType = "counter"
	GaugeType        MetricType = "gauge"
	DistributionType MetricType = "distribution"
	SetType          MetricType = "set"
)

type Metadata struct {
	MetricType  MetricType     `json:"metricType"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Unit        Unit           `json:"unit"`
	Keys        []string       `json:"keys"`
	ValueTypes  []reflect.Type `json:"-"`
	File        string         `json:"file"`
	Line        int            `json:"line"`
}

var defs xsync.Map[string, *Metadata]
var badDefsCallersFrames atomic.Int64
var badDefsNotAtInit atomic.Int64

// https://docs.datadoghq.com/metrics/custom_metrics/#naming-custom-metrics
var nameRegexp = regexp.MustCompile("^[a-z][a-zA-Z0-9_.]{0,199}$")

// https://docs.datadoghq.com/getting_started/tagging/
//
// : is allowed in tags, but because the first : is also used to mark the end of the key and the
// beginning of the value, we don't allow them here.
//
// Also, empty string is accepted, which makes the tag entirely into whatever the value is.
var tagKeyRegexp = regexp.MustCompile("^(|[a-z][a-zA-Z0-9_./-]{0,199})$")

// From https://docs.datadoghq.com/getting_started/tagging/
var reservedTagKeys = map[string]struct{}{
	"host":    {},
	"device":  {},
	"source":  {},
	"service": {},
	"env":     {},
	"version": {},
}

// Returns false if the metric definition is invalid, and so should not emit.
func registerDef(
	metricType MetricType,
	name string,
	description string,
	unit Unit,
	keys []string,
	valueTypes []reflect.Type,
) bool {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		badDefsCallersFrames.Add(1)
		return false
	}
	fn := runtime.FuncForPC(pc)
	if !strings.HasSuffix(fn.Name(), ".init") {
		badDefsNotAtInit.Add(1)
		return false
	}

	// Now we know it's init-time, which means it's safe to panic.

	if !nameRegexp.MatchString(name) {
		panic(fmt.Sprintf(
			"metric definition's name %q doesn't match required %s (see "+
				"https://docs.datadoghq.com/metrics/custom_metrics/#naming-custom-metrics)\n\n"+
				"metric %s defined at %s:%d",
			name,
			nameRegexp,
			name,
			file,
			line,
		))
	}
	if !(strings.HasSuffix(file, "/metrics.go") || strings.HasSuffix(file, "_example_test.go")) {
		panic(fmt.Sprintf(
			"metric definitions must be defined in init() or a top-level var block of a "+
				"file named metrics.go\n\n"+
				"metric %s defined at %s:%d",
			name, file, line,
		))
	}
	if len(description) > 400 {
		panic(fmt.Sprintf(
			"metric descriptions cannot be more than 400 characters, this one is %d\n\n"+
				"metric %s defined at %s:%d",
			len(description), name, file, line,
		))
	}

	d, loaded := defs.LoadOrStore(name, &Metadata{
		MetricType:  metricType,
		Name:        name,
		Description: description,
		Unit:        unit,
		Keys:        xslices.Clone(keys),
		ValueTypes:  valueTypes,
		File:        file,
		Line:        line,
	})
	if loaded {
		panic(fmt.Sprintf(
			"multiple definitions for metric %s:\n"+
				"\t%s:%d\n"+
				"\t%s:%d",
			name,
			d.File, d.Line,
			file, line,
		))
	}

	seenKeys := make(map[string]bool, len(keys))
	for _, key := range keys {
		_, ok := reservedTagKeys[key]
		if ok {
			panic(fmt.Sprintf(
				"metric used reserved tag key %q (see "+
					"https://docs.datadoghq.com/getting_started/tagging/#overview)\n\n"+
					"metric %s defined at %s:%d",
				key,
				name,
				file,
				line,
			))
		}
		if !tagKeyRegexp.MatchString(key) {
			panic(fmt.Sprintf(
				"metric tag key %q doesn't match %s (see "+
					"https://docs.datadoghq.com/getting_started/tagging/#define-tags)\n\n"+
					"metric %s defined at %s:%d",
				key,
				tagKeyRegexp,
				name,
				file,
				line,
			))
		}
		if key != "" && seenKeys[key] {
			panic(fmt.Sprintf(
				"duplicate tag key %q\n\n"+
					"metric %s defined at %s:%d",
				key,
				name,
				file,
				line,
			))
		}
		seenKeys[key] = true
	}

	return true
}

// Defs returns metadata about all of the metric definitions in this process.
//
// Since metrics are registered during init-time, this should be called only after main() has
// already begun.
func Defs() map[string]Metadata {
	result := make(map[string]Metadata)
	defs.Range(func(name string, m *Metadata) bool {
		result[name] = *m
		return true
	})
	return result
}

// DumpDefs prints JSON-formatted metadata about all of the metrics defined in this process to
// stdout.
//
// Since metrics are registered during init-time, this should be called only after main() has
// already begun.
func DumpDefs() error {
	b, err := json.MarshalIndent(Defs(), "" /*prefix*/, "  " /*indent*/)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func join[E any](a []E, b []E) []E {
	if len(a) == 0 {
		return b
	}
	return append(a, b...)
}
