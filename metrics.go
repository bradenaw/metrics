// package metrics wraps github.com/DataDog/datadog-go/v5/statsd with a more ergonomic and
// computationally cheaper interface.
//
// This is done by separating tags from logging metrics so that for frequently-logged gauges and
// counters logging is just a single atomic operation.
//
// For each metric type of Gauge, Count, Histogram, Distribution, and Set, there are a set of
// NewMDefY methods where M is the metric type and Y is the number of tags. A Def can be bound to a
// set of tags with Bind(), producing a metric that can be logged to. It's intended for high
// throughput metrics to hold onto the metric produced by Bind(). By convention, calls to NewMDefY
// should be done at init time, ideally in a var block of a metrics.go file with names as full
// literals so that metrics are easily greppable.
package metrics

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bradenaw/juniper/xsort"
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

	gauges   xsync.Map[metricKey, *Gauge]
	counters xsync.Map[metricKey, *Counter]

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
)

func New(p Publisher) *Metrics {
	m := &Metrics{
		p:       p,
		bg:      xsync.NewGroup(context.Background()),
		flushed: make(chan struct{}),
		polls:   make(map[int]func()),
	}

	m.flushNow = m.bg.PeriodicOrTrigger(flushInterval, 0 /*jitter*/, func(ctx context.Context) {
		m.m.Lock()
		polls := maps.Values(m.polls)
		m.m.Unlock()
		for _, poll := range polls {
			poll()
		}

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

func (m *Metrics) gauge(name string, tags []string) *Gauge {
	k := newMetricKey(name, tags)
	g, ok := m.gauges.Load(k)
	if !ok {
		g = &Gauge{
			m:    m,
			name: name,
			tags: tags,
		}
		g.v.Store(math.Float64bits(math.NaN()))
		g, _ = m.gauges.LoadOrStore(k, g)
	}
	return g
}

func (m *Metrics) counter(name string, tags []string) *Counter {
	k := newMetricKey(name, tags)
	c, ok := m.counters.Load(k)
	if !ok {
		c = &Counter{
			m:    m,
			name: name,
			tags: tags,
		}
		c, _ = m.counters.LoadOrStore(k, c)
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
	old := math.Float64frombits(g.v.Swap(math.Float64bits(v)))
	if math.IsNaN(old) && !math.IsNaN(v) {
		g.publishValue(v)
	}
}

func (g *Gauge) Unset() {
	g.v.Store(math.Float64bits(math.NaN()))
}

func (g *Gauge) publish() {
	v := math.Float64frombits(g.v.Load())
	if math.IsNaN(v) {
		return
	}
	g.publishValue(v)
}

func (g *Gauge) publishValue(v float64) {
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

type Set[K any] struct {
	m          *Metrics
	name       string
	tags       []string
	sampleRate float64
}

func (s *Set[K]) Observe(value K) {
	s.m.p.Set(s.name, fmt.Sprint(value), s.tags, s.sampleRate)
}

// metricKey is used to dedupe metrics so that multiple Bind() calls on a def result in the same
// metric. It contains the name and tags.
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

type metricType int

const (
	counterType metricType = iota + 1
	gaugeType
	histogramType
	distributionType
	setType
)

type metadata struct {
	metricType   metricType
	name         string
	unit         Unit
	description  string
	multipleDefs atomic.Bool
}

var defs xsync.Map[string, *metadata]

func registerDef(metricType metricType, name string, unit Unit, description string) {
	d, loaded := defs.LoadOrStore(name, &metadata{
		metricType:  metricType,
		name:        name,
		unit:        unit,
		description: description,
	})
	if loaded {
		d.multipleDefs.Store(true)
	}
}

// Prints the metrics defined by this process in the format accepted by Datadog's API for metric
// metadata.
//
// https://docs.datadoghq.com/api/latest/metrics/#edit-metric-metadata
func FormatMetadataJSON() string {
	ms := metadatasByName()
	var sb strings.Builder

	writeMetadataJSON := func(name string, unit Unit, description string, multipleDefs bool) {
		type metadataJSON struct {
			Unit        string `json:"unit"`
			Description string `json:"description"`
		}

		_, _ = sb.WriteString(name)
		_, _ = sb.WriteString(" ")
		if multipleDefs {
			description += "\n\nWARNING: multiple defs in code, possibly conflicting"
		}
		b, err := json.Marshal(&metadataJSON{
			Unit:        string(unit),
			Description: description,
		})
		if err != nil {
			panic(err)
		}
		_, _ = sb.Write(b)
		_, _ = sb.WriteString("\n")
	}

	// https://docs.datadoghq.com/metrics/types

	for _, m := range ms {
		switch m.metricType {
		case counterType, gaugeType, setType:
			writeMetadataJSON(m.name, m.unit, m.description, m.multipleDefs.Load())
		case histogramType:
			for _, suffix := range [...]string{"avg", "median", "95percentile", "max"} {
				writeMetadataJSON(m.name+"."+suffix, m.unit, m.description, m.multipleDefs.Load())
			}
			writeMetadataJSON(m.name+".count", UnitEvent, m.description, m.multipleDefs.Load())
		case distributionType:
			for _, prefix := range [...]string{
				"avg",
				"max",
				"min",
				"sum",
				"p50",
				"p75",
				"p90",
				"p95",
				"p99",
			} {
				writeMetadataJSON(prefix+":"+m.name, m.unit, m.description, m.multipleDefs.Load())
			}
			writeMetadataJSON("count:"+m.name, UnitEvent, m.description, m.multipleDefs.Load())
		}
	}

	return sb.String()
}

func metadatasByName() []*metadata {
	var result []*metadata
	defs.Range(func(_ string, m *metadata) bool {
		result = append(result, m)
		return true
	})
	xsort.Slice(result, func(a, b *metadata) bool {
		return a.name < b.name
	})
	return result
}

func joinStrings(a []string, b []string) []string {
	if len(a) == 0 {
		return b
	}
	return append(a, b...)
}
