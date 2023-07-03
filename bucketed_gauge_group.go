package metrics

import (
	"github.com/bradenaw/juniper/xsort"
)

type BucketedGaugeGroup struct {
	boundaries []float64
	gauges     []*Gauge
	pending    []float64
}

// NewBucketedGaugeGroup returns a group of gauges that will emit the number of observations in each
// bucket.
//
// By convention, the key for d is "bucket."
func NewBucketedGaugeGroup(
	m *Metrics,
	d GaugeDef1[string],
	boundaries []float64,
) *BucketedGaugeGroup {
	names := bucketNames(boundaries)
	gauges := make([]*Gauge, len(names))
	for i, name := range names {
		gauges[i] = m.Gauge(d.Values(name))
	}
	return &BucketedGaugeGroup{
		boundaries: boundaries,
		pending:    make([]float64, len(gauges)),
	}
}

func (gg *BucketedGaugeGroup) Observe(v float64) {
	idx := xsort.Search(gg.boundaries, xsort.OrderedLess[float64], v)
	gg.pending[idx]++
}

func (gg *BucketedGaugeGroup) Emit() {
	for i := range gg.pending {
		gg.gauges[i].Set(gg.pending[i])
		gg.pending[i] = 0
	}
}
