package metrics

import (
	"fmt"
	"math"
	"sort"

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
//
// Boundaries must be in increasing order.
func NewBucketedCounter(
	m *Metrics,
	d CounterDef1[string],
	boundaries []float64,
) *BucketedCounter {
	if !boundariesSortedAndUnique(boundaries) {
		boundaries = nil
	}
	names := bucketNames(boundaries)
	counters := make([]*Counter, len(names))
	for i, name := range names {
		counters[i] = m.Counter(d.Values(name))
	}
	return &BucketedCounter{
		boundaries: boundaries,
		counters:   counters,
	}
}

func (b *BucketedCounter) Observe(v float64) {
	idx := xsort.Search(b.boundaries, xsort.OrderedLess[float64], v)
	if idx < len(b.boundaries) && v == b.boundaries[idx] {
		idx++
	}
	b.counters[idx].Add(1)
}

func bucketNames(boundaries []float64) []string {
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

	if len(boundaries) == 0 {
		return []string{""}
	}

	results := make([]string, len(boundaries)+1)
	results[0] = fmt.Sprintf("lt_%g", boundaries[0])
	for i := 1; i < len(boundaries); i++ {
		results[i] = fmt.Sprintf(
			"gte_%g_lt_%g",
			boundaries[i-1],
			boundaries[i],
		)
	}
	results[len(boundaries)] = fmt.Sprintf(
		"gte_%g",
		boundaries[len(boundaries)-1],
	)
	return results
}

// ExponentialBuckets returns exponentially-increasing bucket boundaries for use with
// NewBucketedCounter and NewBucketedGaugeGroup.
//
// start is the first value, base is the base of the exponent, and n is the number of boundaries.
func ExponentialBuckets(start float64, base float64, n int) []float64 {
	boundaries := make([]float64, n)
	for i := range boundaries {
		boundaries[i] = start * math.Pow(base, float64(i))
	}
	return boundaries
}

// LinearBuckets returns linearly-increasing bucket boundaries for use with NewBucketedCounter and
// NewBucketedGaugeGroup.
//
// start is the first value, step is the distance between values, and n is the number of boundaries.
func LinearBuckets(start float64, step float64, n int) []float64 {
	boundaries := make([]float64, n)
	for i := range boundaries {
		boundaries[i] = start + step*float64(i)
	}
	return boundaries
}

func boundariesSortedAndUnique(s []float64) bool {
	if !sort.Float64sAreSorted(s) {
		return false
	}
	for i := 0; i < len(s)-1; i++ {
		if s[i] == s[i+1] {
			return false
		}
	}
	return true
}
