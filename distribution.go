package metrics

import (
	"math"
	"strconv"
	"sync/atomic"

	"github.com/bradenaw/juniper/xsync"
)

// We'll round each Distribution value by at most this much, as a ratio. e.g. 0.03 means at most 3%
// error for each observation. The tradeoff is for resources: more error bound means fewer buckets,
// which means we can both store the sketch more efficiently and communicate about it more
// efficiently.
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
func (s *concurrentSketch) newObservations(f func(value distributionBucket, count int) bool) {
	// TODO: track the number of unchanged buckets, and shrink the sketch if it's a high enough
	// percentage
	// shrinking requires another rwmutex so we can swap out the maps...

	done := false
	s.positive.Range(func(bucket int16, counter *atomic.Uint32) bool {
		count := int(counter.Load())
		diff := count - s.prevPositive[bucket]
		if diff == 0 {
			return true
		}

		if s.prevPositive == nil {
			s.prevPositive = make(map[int16]int)
		}
		s.prevPositive[bucket] = count
		if !done {
			if !f(distributionBucket{bucket: bucket, positive: true}, diff) {
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

		if s.prevNegative == nil {
			s.prevNegative = make(map[int16]int)
		}
		s.prevNegative[bucket] = count
		if !done {
			if !f(distributionBucket{bucket: bucket, positive: false}, diff) {
				done = true
			}
		}
		return true
	})
}

type distributionBucket struct {
	bucket   int16
	positive bool
}

func (bucket distributionBucket) Value() float64 {
	v := math.Pow(1+distributionErrorBound, float64(bucket.bucket))
	if !bucket.positive {
		return -v
	}
	return v
}

func (bucket distributionBucket) AppendString(b []byte) []byte {
	if bucket.bucket >= minLookupBucketStr && bucket.bucket <= maxLookupBucketStr {
		if bucket.positive {
			return append(b, positiveBucketStrs[bucket.bucket-minLookupBucketStr]...)
		} else {
			return append(b, negativeBucketStrs[bucket.bucket-minLookupBucketStr]...)
		}
	}
	return bucket.appendStringNoLookup(b)
}

func (bucket distributionBucket) appendStringNoLookup(b []byte) []byte {
	// 4 is plenty of sigfigs because ddagent is going to put them in a sketch that rounds by 3%
	// anyway.
	//
	// 'g' will put very large/small values in scientific format. The docs don't promise that
	// they'll parse this, but the docs don't really say anything (they just say that values are "An
	// integer or float."), so depending only on documented behavior is obviously nonsense. But
	// ddagent does indeed call strconv.ParseFloat[1], so whatever.
	//
	// [1]: https://github.com/DataDog/datadog-agent/blob/486de8dc3f9462f403a7135374e7306ef05861e5/comp/dogstatsd/server/parse.go#L167
	return strconv.AppendFloat(b, bucket.Value(), 'g', 4, 64)
}

const (
	minLookupBucketStr = -512
	maxLookupBucketStr = 512
)

var (
	positiveBucketStrs [maxLookupBucketStr - minLookupBucketStr + 1][]byte
	negativeBucketStrs [maxLookupBucketStr - minLookupBucketStr + 1][]byte
)

func init() {
	for i := int16(minLookupBucketStr); i <= maxLookupBucketStr; i++ {
		positiveBucketStrs[i-minLookupBucketStr] = distributionBucket{bucket: i, positive: true}.appendStringNoLookup(nil)
		negativeBucketStrs[i-minLookupBucketStr] = distributionBucket{bucket: i, positive: false}.appendStringNoLookup(nil)
	}
}
