package metrics

import (
	"fmt"
	"math"

	"github.com/bradenaw/juniper/xsort"
)

// BucketedCounter is a counter sectioned into buckets.
//
// For example, with boundaries []float64{100, 200, 400}, BucketedCounter will produce four counters
// with the following tag values:
//
//	lt_100             which counts Observe()s with v < 100
//	gte_100_lt_200     which counts Observe()s with 100 <= v < 200
//	gte_200_lt_400     which counts Observe()s with 200 <= v < 400
//	gte_400            which counts Observe()s with 400 <= v
type BucketedCounter struct {
	boundaries []float64
	counters   []*Counter
}

// NewBucketedCounter returns a counter that keeps track of observed values between the given
// boundaries.
//
// By convention, the key for d is "bucket."
func NewBucketedCounter(
	m *Metrics,
	d *CounterDef1[string],
	boundaries []float64,
) *BucketedCounter {
	if len(boundaries) == 0 {
		return &BucketedCounter{
			boundaries: nil,
			counters:   []*Counter{d.Bind(m, "")},
		}
	}

	// https://docs.datadoghq.com/getting_started/tagging/
	//
	// > Tags must start with a letter and after that may contain the characters listed below:
	// >
	// > - Alphanumerics
	// > - Underscores
	// > - Minuses
	// > - Colons
	// > - Periods
	// > - Slashes
	//
	// [-./] all have meanings in numbers, colons already used for tag key:values, so that leaves
	// alphanum and _

	counters := make([]*Counter, len(boundaries)+1)
	counters[0] = d.Bind(m, fmt.Sprintf(
		"lt_%f",
		boundaries[0],
	))
	for i := 1; i < len(boundaries); i++ {
		counters[i] = d.Bind(m, fmt.Sprintf(
			"gte_%f_lt_%f",
			boundaries[i-1],
			boundaries[i],
		))
	}
	counters[len(boundaries)] = d.Bind(m, fmt.Sprintf(
		"gte_%f",
		boundaries[len(boundaries)-1],
	))

	return &BucketedCounter{
		boundaries: boundaries,
		counters:   counters,
	}
}

func (b *BucketedCounter) Observe(v float64) {
	idx := xsort.Search(b.boundaries, xsort.OrderedLess[float64], v)
	b.counters[idx].Add(1)
}

func ExponentialBuckets(start float64, base float64, n int) []float64 {
	buckets := make([]float64, n)
	for i := range buckets {
		buckets[i] = start * math.Pow(base, float64(i))
	}
	return buckets
}
