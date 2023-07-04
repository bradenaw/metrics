// package metrics wraps github.com/DataDog/datadog-go/v5/statsd with a more ergonomic and
// computationally cheaper interface.
//
// This is done by separating tags from logging metrics so that for frequently-logged gauges and
// counters logging is just a single atomic operation.
//
// For each metric type of Gauge, Count, Histogram, Distribution, and Set, there are a set of
// NewMDefY methods where M is the metric type and Y is the number of tags. Calls to NewMDefY must
// be done at init-time (ideally in a top-level var block) of a metrics.go file with names as full
// literals so that metrics are easily greppable. Metrics not defined this way will cause the
// process to panic if still at init-time, meaning before any code in main() has run, otherwise will
// produce non-functional stats and produce to a gauge stat called metrics.bad_metric_definitions.
// It's a good idea to put an alert on this stat so that if it starts logging during a deploy, you
// know your other metrics may not be trustworthy.
//
// See the example folder for an example of usage.
//
// Generally you will have:
//  1. metrics.NewMDefY calls in a top-level var block of metrics.go in packages that log metrics,
//     named ending in "Def".
//  2. metrics.New() in main(), passed down through constructors.
//  3. At the point of logging metrics, e.g.
//     m.Counter(myCounterDef.Values(tag1, tag2)).Add(1)
package metrics

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"runtime"
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
	Histogram(name string, value float64, tags []string, rate float64) error
	Distribution(name string, value float64, tags []string, rate float64) error
	Set(name string, value string, tags []string, rate float64) error
}

// TagValue is the value of a key:value pair in a metric tag. They are formatted the same as
// fmt.Sprint.
type TagValue any

type Metrics struct {
	p        Publisher
	bg       *xsync.Group
	flushNow func()

	gauges        xsync.Map[metricKey, *Gauge]
	counters      xsync.Map[metricKey, *Counter]
	histograms    xsync.Map[metricKey, *Histogram]
	distributions xsync.Map[metricKey, *Distribution]
	sets          xsync.Map[metricKey, *Set]

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
func (p noOpPublisher) Histogram(name string, value float64, tags []string, rate float64) error {
	return nil
}
func (p noOpPublisher) Distribution(name string, value float64, tags []string, rate float64) error {
	return nil
}
func (p noOpPublisher) Set(name string, value string, tags []string, rate float64) error { return nil }

var (
	NoOpMetrics = &Metrics{
		p:       noOpPublisher{},
		bg:      xsync.NewGroup(context.Background()),
		flushed: make(chan struct{}),
		polls:   make(map[int]func()),
	}

	// Used to return from Metrics.Metric() methods when the definition is invalid and the stat
	// can't be logged.
	noOpCounter      = &Counter{m: NoOpMetrics}
	noOpGauge        = &Gauge{m: NoOpMetrics}
	noOpHistogram    = &Histogram{m: NoOpMetrics}
	noOpDistribution = &Distribution{m: NoOpMetrics}
	noOpSet          = &Set{m: NoOpMetrics}

	badDefsDef = NewGaugeDef1[string](
		"metrics.bad_metric_definitions",
		"The number of calls to NewMGaugeY that are invalid for some reason. These definitions "+
			"will not be able to log metrics at all.",
		UnitItem,
		[...]string{"reason"},
	)
)

func New(p Publisher) *Metrics {
	m := &Metrics{
		p:       p,
		bg:      xsync.NewGroup(context.Background()),
		flushed: make(chan struct{}),
		polls:   make(map[int]func()),
	}

	badDefsCallersFramesGauge := m.Gauge(badDefsDef.Values("runtime_caller_failed"))
	badDefsNotAtInitGauge := m.Gauge(badDefsDef.Values("not_at_init_time"))

	m.flushNow = m.bg.PeriodicOrTrigger(flushInterval, 0 /*jitter*/, func(ctx context.Context) {
		m.m.Lock()
		polls := maps.Values(m.polls)
		m.m.Unlock()
		for _, poll := range polls {
			poll()
		}

		badDefsCallersFramesGauge.Set(float64(badDefsCallersFrames.Load()))
		badDefsNotAtInitGauge.Set(float64(badDefsNotAtInit.Load()))

		m.gauges.Range(func(_ metricKey, g *Gauge) bool {
			g.publish()
			return true
		})
		m.counters.Range(func(_ metricKey, c *Counter) bool {
			c.publish()
			return true
		})

		m.m.Lock()
		close(m.flushed)
		m.flushed = make(chan struct{})
		m.m.Unlock()
	})

	return m
}

// Counter returns the Counter for the given CounterDef. For metrics with tags (e.g. CounterDef2),
// the CounterDef can be made by calling Values(), for example:
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

	k := newMetricKey(d.name, d.tags)
	c, ok := m.counters.Load(k)
	if !ok {
		c = &Counter{
			m:    m,
			name: d.name,
			tags: d.tags,
		}
		c, _ = m.counters.LoadOrStore(k, c)
	}
	return c
}

func (m *Metrics) Gauge(d GaugeDef) *Gauge {
	if !d.ok {
		return noOpGauge
	}

	k := newMetricKey(d.name, d.tags)
	g, ok := m.gauges.Load(k)
	if !ok {
		g = &Gauge{
			m:    m,
			name: d.name,
			tags: d.tags,
		}
		g.v.Store(math.Float64bits(math.NaN()))
		g, _ = m.gauges.LoadOrStore(k, g)
	}
	return g
}

func (m *Metrics) Histogram(d HistogramDef) *Histogram {
	if !d.ok {
		return noOpHistogram
	}

	k := newMetricKey(d.name, d.tags)
	c, ok := m.histograms.Load(k)
	if !ok {
		c = &Histogram{
			m:          m,
			name:       d.name,
			tags:       d.tags,
			sampleRate: d.sampleRate,
		}
		c, _ = m.histograms.LoadOrStore(k, c)
	}
	return c
}

func (m *Metrics) Distribution(d DistributionDef) *Distribution {
	if !d.ok {
		return noOpDistribution
	}

	k := newMetricKey(d.name, d.tags)
	c, ok := m.distributions.Load(k)
	if !ok {
		c = &Distribution{
			m:          m,
			name:       d.name,
			tags:       d.tags,
			sampleRate: d.sampleRate,
		}
		c, _ = m.distributions.LoadOrStore(k, c)
	}
	return c
}

func (m *Metrics) Set(d SetDef) *Set {
	if !d.ok {
		return noOpSet
	}

	k := newMetricKey(d.name, d.tags)
	c, ok := m.sets.Load(k)
	if !ok {
		c = &Set{
			m:          m,
			name:       d.name,
			tags:       d.tags,
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

	id := m.nextID
	m.nextID++
	m.polls[id] = f

	return func() {
		m.m.Lock()
		defer m.m.Unlock()
		delete(m.polls, id)
	}
}

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

func (m *Metrics) Close() {
	m.bg.StopAndWait()
}

// Gauge is a metric that reports the last value that it was set to.
//
// Unlike Datadog's native gauge in the statsd client, Gauges report this value until the end of the
// process or until explicitly Unset().
type Gauge struct {
	m    *Metrics
	name string
	tags []string
	v    atomic.Uint64
}

func (g *Gauge) Set(v float64) {
	g.v.Store(math.Float64bits(v))
}

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
	g.m.p.Gauge(g.name, v, g.tags, 1 /*samplingRate*/)
}

// Counter is a metric that keeps track of the number of events that happen per time interval.
type Counter struct {
	m    *Metrics
	name string
	tags []string
	v    atomic.Int64
}

func (c *Counter) Add(n int64) {
	c.v.Add(n)
}

func (c *Counter) publish() {
	v := c.v.Swap(0)
	if v > 0 {
		c.m.p.Count(c.name, v, c.tags, 1)
	}
}

type Histogram struct {
	m          *Metrics
	name       string
	tags       []string
	sampleRate float64
}

func (h *Histogram) Observe(value float64) {
	h.m.p.Histogram(h.name, value, h.tags, h.sampleRate)
}

type Distribution struct {
	m          *Metrics
	name       string
	tags       []string
	sampleRate float64
}

func (h *Distribution) Observe(value float64) {
	h.m.p.Distribution(h.name, value, h.tags, h.sampleRate)
}

type Set struct {
	m          *Metrics
	name       string
	tags       []string
	sampleRate float64
}

func (s *Set) Observe(value string) {
	s.m.p.Set(s.name, value, s.tags, s.sampleRate)
}

// metricKey is used to dedupe metrics so that multiple calls on a def result in the same metric. It
// contains the name and tags.
type metricKey string

func newMetricKey(name string, tags []string) metricKey {
	if len(tags) == 0 {
		return metricKey(name)
	}

	n := len(name) + 1
	for _, tag := range tags {
		n += len(tag) + 1
	}

	var sb strings.Builder
	sb.Grow(n)
	_, _ = sb.WriteString(name)
	_, _ = sb.WriteString(":")
	for i, tag := range tags {
		if i != 0 {
			_, _ = sb.WriteString(",")
		}
		_, _ = sb.WriteString(tag)
	}
	return metricKey(sb.String())
}

func makeTag(key string, value TagValue) string {
	if len(key) == 0 {
		return tagValueString(value)
	}
	return key + ":" + tagValueString(value)
}

func tagValueString(v TagValue) string {
	return fmt.Sprint(v)
}

type MetricType string

const (
	CounterType      MetricType = "counter"
	GaugeType        MetricType = "gauge"
	HistogramType    MetricType = "histogram"
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
var nameRegexp = regexp.MustCompile("^[a-z][a-zA-Z0-9_.]{0,199}")

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
			"metric definition's name doesn't match required %s\n\n"+
				"metric %s defined at %s:%d",
			nameRegexp, name, file, line,
		))
	}
	if !strings.HasSuffix(file, "/metrics.go") {
		panic(fmt.Sprintf(
			"metric definitions must be defined in init() or a top-level var block of a "+
				"file named metrics.go\n\n"+
				"metric %s defined at %s:%d",
			name, file, line,
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

	return true
}

// Defs returns metadata about all of the metric definitions in this binary. Since metrics are
// registered during init-time, this should be called only after main() has already begun.
func Defs() map[string]Metadata {
	result := make(map[string]Metadata)
	defs.Range(func(name string, m *Metadata) bool {
		result[name] = *m
		return true
	})
	return result
}

func joinStrings(a []string, b []string) []string {
	if len(a) == 0 {
		return b
	}
	return append(a, b...)
}
