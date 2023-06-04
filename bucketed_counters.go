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
	c *Counter,
	boundaries []float64,
) *BucketedCounter {
	var counters []*Counter

	if len(boundaries) == 0 {
		counters = []*Counter{c}
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
		// [-./] all have meanings in numbers, colons already used for tag key:values, so that leaves alphanum and _

		addCounter := func(tag string) *Counter {
			return c.m.counter(c.name, append(c.tags[0:len(c.tags):len(c.tags)], tag))
		}

		counters = make([]*Counter, len(boundaries)+1)
		counters[0] = addCounter(fmt.Sprintf("bucket:lt_%f", boundaries[0]))
		for i := 1; i < len(boundaries); i++ {
			counters[i] = addCounter(fmt.Sprintf("bucket:in_%f_%f", boundaries[i-1], boundaries[i]))
		}
		counters[len(boundaries)] = addCounter(fmt.Sprintf("bucket:gt_%f", boundaries[len(boundaries)-1]))
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
