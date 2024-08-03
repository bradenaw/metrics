package metrics

import (
	"math"
	"sync/atomic"

	"github.com/bradenaw/juniper/xsync"
)

const distributionErrorBound = 0.03

type concurrentSketch struct {
	positive     xsync.Map[int16, *atomic.Uint32]
	negative     xsync.Map[int16, *atomic.Uint32]
	prevPositive map[int16]int
	prevNegative map[int16]int
}

func (c *concurrentSketch) Observe(v float64) {
	target := &c.positive
	if v < 0 {
		v = -v
		target = &c.negative
	}
	// We don't have to worry about overflow here because the largest bucket number is 2^16-1 and
	// 1.03^(2^16-1) is larger than the max float64.
	bucket := int16(math.Trunc(math.Log(v) / math.Log(1+distributionErrorBound)))
	counter, ok := target.Load(bucket)
	if !ok {
		counter, _ = target.LoadOrStore(bucket, new(atomic.Uint32))
	}
	counter.Add(1)
}

// Calls f for all of the Observe()s that happened since the last call to newObservations.
func (s *concurrentSketch) newObservations(f func(value float64, count int) bool) {
	// TODO: track the number of unchanged buckets, and shrink the sketch if it's a high enough
	// percentage

	done := false
	s.positive.Range(func(bucket int16, counter *atomic.Uint32) bool {
		count := int(counter.Load())
		diff := count - s.prevPositive[bucket]
		if diff == 0 {
			return true
		}

		value := math.Pow(1+distributionErrorBound, float64(bucket))
		if s.prevPositive == nil {
			s.prevPositive = make(map[int16]int)
		}
		s.prevPositive[bucket] = count
		if !done {
			if !f(value, diff) {
				done = true
			}
		}
		return true
	})
	s.negative.Range(func(bucket int16, counter *atomic.Uint32) bool {
		count := int(counter.Load())
		diff := count - s.prevNegative[bucket]
		if diff == 0 {
			return true
		}

		value := -math.Pow(1+distributionErrorBound, float64(bucket))
		if s.prevNegative == nil {
			s.prevNegative = make(map[int16]int)
		}
		s.prevNegative[bucket] = count
		if !done {
			if !f(value, diff) {
				done = true
			}
		}
		return true
	})
}
