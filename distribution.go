package metrics

import (
	"math"
	"sync/atomic"

	"github.com/bradenaw/juniper/xsync"
)

const distributionErrorBound = 0.03

type concurrentSketch struct {
	positive xsync.Map[uint16, *atomic.Uint32]
	negative xsync.Map[uint16, *atomic.Uint32]
}

func (c *concurrentSketch) Observe(v float64) {
	target := &c.positive
	if v < 0 {
		v = -v
		target = &c.negative
	}
	bucket := uint16(math.Trunc(math.Log(v) / math.Log(1+distributionErrorBound)))
	counter, ok := target.Load(bucket)
	if !ok {
		counter, _ = target.LoadOrStore(bucket, new(atomic.Uint32))
	}
	counter.Add(1)
}

type sketch struct {
	total    int
	positive map[uint16]int
	negative map[uint16]int
}

func (s *sketch) cloneFrom(c *concurrentSketch) {
	for bucket := range s.positive {
		delete(s.positive, bucket)
	}
	for bucket := range s.negative {
		delete(s.negative, bucket)
	}

	s.total = 0

	c.positive.Range(func(bucket uint16, counter *atomic.Uint32) bool {
		if s.positive == nil {
			s.positive = make(map[uint16]int)
		}
		count := int(counter.Load())
		s.positive[bucket] = count
		s.total += count
		return true
	})

	c.negative.Range(func(bucket uint16, counter *atomic.Uint32) bool {
		if s.negative == nil {
			s.negative = make(map[uint16]int)
		}
		count := int(counter.Load())
		s.negative[bucket] = count
		s.total += count
		return true
	})
}

func (s *sketch) diffIter(other *sketch, f func(value float64, count int) bool) {
	// NOTE: Doesn't count buckets that only appear in other, but because of our usage of this, that
	// never happens because we only ever diff running totals against each other where s is the
	// newer.
	for bucket := range s.positive {
		diff := s.positive[bucket] - other.positive[bucket]
		if diff == 0 {
			continue
		}
		value := math.Pow(1+distributionErrorBound, float64(bucket))
		if !f(value, diff) {
			return
		}
	}
	for bucket := range s.negative {
		diff := s.negative[bucket] - other.negative[bucket]
		if diff == 0 {
			continue
		}
		value := -math.Pow(1+distributionErrorBound, float64(bucket))
		if !f(value, diff) {
			return
		}
	}
}

func (s *sketch) values(f func(value float64, count int) bool) {
	for bucket, count := range s.positive {
		value := math.Pow(1+distributionErrorBound, float64(bucket))
		if !f(value, count) {
			break
		}
	}
	for bucket, count := range s.negative {
		value := -math.Pow(1+distributionErrorBound, float64(bucket))
		if !f(value, count) {
			break
		}
	}
}
