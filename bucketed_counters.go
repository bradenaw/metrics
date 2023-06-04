package metrics

import (
	"fmt"

	"github.com/bradenaw/juniper/xsort"
)

type BucketedCounter struct {
	boundaries []float64
	counters   []*Counter
}

func NewBucketedCounter(
	m *Metrics,
	def *CounterDef1[string],
	boundaries []float64,
) *BucketedCounter {
	var counters []*Counter

	if len(boundaries) == 0 {
		counters = []*Counter{def.Bind(m, "")}
	} else {
		// https://docs.datadoghq.com/getting_started/tagging/
		//
		// Tags must start with a letter and after that may contain the characters listed below:
		//
		// - Alphanumerics
		// - Underscores
		// - Minuses
		// - Colons
		// - Periods
		// - Slashes
		//
		// [-./] all have meanings in numbers, so that leaves alphanum and [:_]
		//
		// Colons already used for tag key:values.
		counters = make([]*Counter, len(boundaries)+1)
		counters[0] = def.Bind(m, fmt.Sprintf("lt_%f", boundaries[0]))
		for i := 1; i < len(boundaries); i++ {
			counters[i] = def.Bind(m, fmt.Sprintf("in_%f_%f", boundaries[i-1], boundaries[i]))
		}
		counters[len(boundaries)] = def.Bind(m, fmt.Sprintf("gt_%f", boundaries[len(boundaries)-1]))
	}

	return &BucketedCounter{
		boundaries: boundaries,
		counters:   counters,
	}
}

func (b *BucketedCounter) Observe(v float64) {
	idx := xsort.Search(b.boundaries, xsort.OrderedLess[float64], v)
	b.counters[idx].Add(1)
}
