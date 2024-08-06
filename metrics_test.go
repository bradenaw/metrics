package metrics

import (
	"strings"
	"sync"
	"testing"
)

func TestBucketedGaugeGroup(t *testing.T) {
	d := GaugeDef1[string]{
		name:          "test_bucketed_gauge_group",
		keys:          [...]string{"bucket"},
		allComparable: true,
		ok:            true,
	}

	gg := NewBucketedGaugeGroup(NoOpMetrics, d, []float64{1, 10, 100, 1000})

	// lt_1
	gg.Observe(0)
	gg.Observe(0)

	// gte_1_lt_10
	gg.Observe(1)
	gg.Observe(2)
	gg.Observe(9)

	// gte_10_lt_100
	gg.Observe(10)
	gg.Observe(11)
	gg.Observe(12)
	gg.Observe(99)

	// gte_100_lt_1000
	gg.Observe(100)
	gg.Observe(101)
	gg.Observe(102)
	gg.Observe(103)
	gg.Observe(999)

	// gte_1000
	gg.Observe(1000)
	gg.Observe(1001)
	gg.Observe(1002)
	gg.Observe(1003)
	gg.Observe(1004)
	gg.Observe(1005)

	gg.Emit()

	type expected struct {
		bucket string
		count  int
	}
	for _, e := range []expected{
		{"lt_1", 2},
		{"gte_1_lt_10", 3},
		{"gte_10_lt_100", 4},
		{"gte_100_lt_1000", 5},
		{"gte_1000", 6},
	} {
		g, ok := NoOpMetrics.gauges.Load(
			newMetricKey(
				"test_bucketed_gauge_group",
				1, // nTags
				[maxTags]any{e.bucket},
				true,
			),
		)
		if !ok {
			t.Fatalf("bucket %s didn't get created", e.bucket)
		}
		if g.value() != float64(e.count) {
			t.Fatalf("bucket %s has value %f, expected %d", e.bucket, g.value(), e.count)
		}
	}
}

func TestBucketedCounter(t *testing.T) {
	d := CounterDef1[string]{
		name:          "test_bucketed_counter",
		keys:          [...]string{"bucket"},
		allComparable: true,
		ok:            true,
	}

	bc := NewBucketedCounter(NoOpMetrics, d, []float64{1, 10, 100, 1000})

	// lt_1
	bc.Observe(0)
	bc.Observe(0)

	// gte_1_lt_10
	bc.Observe(1)
	bc.Observe(2)
	bc.Observe(9)

	// gte_10_lt_100
	bc.Observe(10)
	bc.Observe(11)
	bc.Observe(12)
	bc.Observe(99)

	// gte_100_lt_1000
	bc.Observe(100)
	bc.Observe(101)
	bc.Observe(102)
	bc.Observe(103)
	bc.Observe(999)

	// gte_1000
	bc.Observe(1000)
	bc.Observe(1001)
	bc.Observe(1002)
	bc.Observe(1003)
	bc.Observe(1004)
	bc.Observe(1005)

	type expected struct {
		bucket string
		count  int
	}
	for _, e := range []expected{
		{"lt_1", 2},
		{"gte_1_lt_10", 3},
		{"gte_10_lt_100", 4},
		{"gte_100_lt_1000", 5},
		{"gte_1000", 6},
	} {
		c, ok := NoOpMetrics.counters.Load(
			newMetricKey(
				"test_bucketed_counter",
				1, // nTags
				[maxTags]any{e.bucket},
				true, // allComparable
			),
		)
		if !ok {
			t.Fatalf("bucket %s didn't get created", e.bucket)
		}
		if c.v.Load() != int64(e.count) {
			t.Fatalf("bucket %s has value %d, expected %d", e.bucket, c.v.Load(), e.count)
		}
	}
}

type capturingPublisher struct {
	mu       sync.Mutex
	counters map[string]int64
}

func (p *capturingPublisher) Gauge(name string, value float64, tags []string, rate float64) error {
	return nil
}
func (p *capturingPublisher) Distribution(name string, value float64, tags []string, rate float64) error {
	return nil
}
func (p *capturingPublisher) Set(name string, value string, tags []string, rate float64) error {
	return nil
}
func (p *capturingPublisher) Count(name string, value int64, tags []string, rate float64) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.counters[p.makeKey(name, tags)] += value
	return nil
}
func (p *capturingPublisher) makeKey(name string, tags []string) string {
	return name + ":" + strings.Join(tags, ",")
}
func (p *capturingPublisher) countSeen(name string, tags []string) int64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.counters[p.makeKey(name, tags)]
}

func TestMetrics(t *testing.T) {
	p := capturingPublisher{counters: make(map[string]int64)}
	m := New(&p)

	def := CounterDef3[string, int, bool]{
		name: "test_metrics_counter",
		keys: [...]string{"string", "int", "bool"},
		ok:   true,
	}

	a := m.Counter(def.Prefix1("foo").Values(123, false))
	b := m.Counter(def.Prefix2("foo", 123).Values(false))
	c := m.Counter(def.Values("foo", 123, false))
	d := m.Counter(def.Values("bar", 456, true))

	a.Add(1)
	b.Add(2)
	c.Add(3)
	d.Add(7)

	m.Flush()

	abcSeen := p.countSeen(def.name, []string{"string:foo", "int:123", "bool:false"})
	if abcSeen != 6 {
		t.Fatalf("expected 6, got %d", abcSeen)
	}
	dSeen := p.countSeen(def.name, []string{"string:bar", "int:456", "bool:true"})
	if dSeen != 7 {
		t.Fatalf("expected 7, got %d", dSeen)
	}
}

func TestTagValueSanitize(t *testing.T) {
	check := func(
		s string,
		expected string,
	) {
		actual := tagValueSanitize(s)
		if actual != expected {
			t.Errorf("tagValueSanitize(%q) -> %q, expected %q", s, actual, expected)
		}
	}

	check("abcdefghijklmnopqrstuvwxyz0123456789-:./", "abcdefghijklmnopqrstuvwxyz0123456789-:./")
	check("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-:./", "abcdefghijklmnopqrstuvwxyz0123456789-:./")
	check("abc?123", "abc_123")
	check(",", "_")
	check("|@", "__")
	// NOTE: ɱ̊ is two runes
	check("2o4?uASfd$j⁛1℘aℵ]ɱ̊Mę14\nq", "2o4_uasfd_j_1_a____m_14_q")
}

func BenchmarkTagValueSanitize(b *testing.B) {
	const alreadyValid = "abcdefghijklmnopqrstuvwxyz0123456789-:./"
	b.Run("AlreadyValid", func(b *testing.B) {
		b.ReportAllocs()

		total := 0
		for i := 0; i < b.N; i++ {
			total += len(tagValueSanitize(alreadyValid))
		}
		b.Log(total)
	})

	const withInvalid = "2o43uasfdaj⁛℘ℵɱ̊ę1230"
	b.Run("WithInvalid", func(b *testing.B) {
		b.ReportAllocs()

		total := 0
		for i := 0; i < b.N; i++ {
			total += len(tagValueSanitize(withInvalid))
		}
		b.Log(total)
	})
}

//func BenchmarkReport(b *testing.B) {
//	m := New(noOpPublisher{})
//
//	counter := m.Counter(CounterDef{name: "benchmark_report_counter", ok: true})
//	b.Run("Counter", func(b *testing.B) {
//		b.ReportAllocs()
//		for i := 0; i < b.N; i++ {
//			counter.Add(1)
//		}
//	})
//
//	gauge := m.Gauge(GaugeDef{name: "benchmark_report_gauge", ok: true})
//	b.Run("Gauge", func(b *testing.B) {
//		b.ReportAllocs()
//		for i := 0; i < b.N; i++ {
//			gauge.Set(float64(i))
//		}
//	})
//
//	distribution := m.Distribution(DistributionDef{name: "benchmark_report_distribution", ok: true})
//	b.Run("Distribution", func(b *testing.B) {
//		b.ReportAllocs()
//		for i := 0; i < b.N; i++ {
//			distribution.Observe(float64(i))
//		}
//	})
//
//	set := m.Set(SetDef{name: "benchmark_report_distribution", ok: true})
//	b.Run("Set", func(b *testing.B) {
//		b.ReportAllocs()
//		for i := 0; i < b.N; i++ {
//			set.Observe("asdf")
//		}
//	})
//}
//
//func BenchmarkMetricLookup(b *testing.B) {
//	m := New(noOpPublisher{})
//	b.ReportAllocs()
//	d := CounterDef3[string, int, bool]{
//		name:          "benchmark_metric_lookup",
//		keys:          [...]string{"", "", ""},
//		allComparable: true,
//		ok:            true,
//	}
//
//	foo := "foo"
//
//	for i := 0; i < b.N; i++ {
//		withValues := d.Values(foo, 1, false)
//		m.Counter(withValues)
//	}
//}
