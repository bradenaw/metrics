package metrics

import (
	"github.com/bradenaw/juniper/xsort"
)

// BucketedGaugeGroup is a group of gauges sectioned by buckets of observed values. This is useful
// for measuring the state of the system in the number of items that fall into each bucket.
//
// For example, if your process has a set of buffers that grow and shrink, you might be curious how
// many of each exist and what size each of them are. A good summary would be the number of buffers
// that are 1KB or less, 1KB-10KB, and 10KB-100KB, and 100KB+. This can be achieved by using a
// BucketedGaugeGroup with bucket boundaries [1024, 10*1024, 100*1024], which will produce four
// gauges with the following tag values:
//
//	lt_1024              which counts Observe()s since the last Emit() with v < 1024
//	gte_1024_lt_10240    which counts Observe()s since the last Emit() with 1024 <= v < 10240
//	gte_10240_lt_102400  which counts Observe()s since the last Emit() with 10240 <= v < 102400
//	gte_102400           which counts Observe()s since the last Emit() with 102400 <= v
//
// BucketedGaugeGroups are usually emitted to using Metrics.EveryFlush.
type BucketedGaugeGroup struct {
	boundaries []float64
	gauges     []*Gauge
	pending    []float64
}

// NewBucketedGaugeGroup returns a group of gauges that will emit the number of observations in each
// bucket.
//
// By convention, the key for d is "bucket."
//
// Boundaries must be in sorted order.
func NewBucketedGaugeGroup(
	m *Metrics,
	d GaugeDef1[string],
	boundaries []float64,
) *BucketedGaugeGroup {
	if !boundariesSortedAndUnique(boundaries) {
		boundaries = nil
	}
	names := bucketNames(boundaries)
	gauges := make([]*Gauge, len(names))
	for i, name := range names {
		gauges[i] = m.Gauge(d.Values(name))
	}
	return &BucketedGaugeGroup{
		boundaries: boundaries,
		gauges:     gauges,
		pending:    make([]float64, len(gauges)),
	}
}

// Observe adds one to the bucket that v falls in.
func (gg *BucketedGaugeGroup) Observe(v float64) {
	idx := xsort.Search(gg.boundaries, xsort.OrderedLess[float64], v)
	if idx < len(gg.boundaries) && v == gg.boundaries[idx] {
		idx++
	}
	gg.pending[idx]++
}

// Emit emits the observations passed to Observe() since the last Emit() as gauges.
func (gg *BucketedGaugeGroup) Emit() {
	for i := range gg.pending {
		gg.gauges[i].Set(gg.pending[i])
		gg.pending[i] = 0
	}
}
