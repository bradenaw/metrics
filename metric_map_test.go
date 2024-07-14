package metrics

import (
	"maps"
	"math/rand/v2"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestMetricMapConcurrent(t *testing.T) {
	const (
		concurrency = 10
		duration    = 3 * time.Second
	)

	start := time.Now()

	var m metricMap[*bool]

	var (
		totalLoads         atomic.Int64
		totalStoreAttempts atomic.Int64
		totalStores        atomic.Int64
		totalRanges        atomic.Int64
	)

	var wg sync.WaitGroup
	expecteds := make([]map[metricKey]*bool, concurrency)
	for i := 0; i < concurrency; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()

			myKeys := make(map[metricKey]*bool)
			expecteds[i] = myKeys
			checkKey := func(k metricKey, existing *bool) {
				expected, mine := myKeys[k]
				if mine && existing != expected {
					t.Fatal("previously stored a value but it got overwritten")
				}
			}

			loads := 0
			storeAttempts := 0
			ranges := 0

			for time.Since(start) < duration {
				k := newMetricKey(
					"", // name
					1,  // n
					[maxTags]any{
						// Quantize the keys a little to increase the chance of collision between
						// goroutines.
						int(time.Since(start).Round(time.Millisecond)) + rand.N(concurrency),
					},
					true, // allComparable
				)

				loads++
				existing, ok := m.Load(k)
				if ok {
					checkKey(k, existing)
				} else {
					v := new(bool)
					storeAttempts++
					v2, loaded := m.LoadOrStore(k, v)
					if !loaded {
						if v != v2 {
							t.Fatal("stored a value but got back something different")
						}
						myKeys[k] = v
					}
				}

				if rand.N(1000) == 0 {
					ranges++
					seen := 0
					m.Range(func(k metricKey, v *bool) bool {
						expected, mine := myKeys[k]
						if !mine {
							return true
						}
						if expected != v {
							t.Fatal("Range saw a different value than the last LoadOrStore")
						}
						seen++
						return true
					})
					if seen != len(myKeys) {
						t.Fatal("Range didn't see a key it should have")
					}
				}
			}

			totalLoads.Add(int64(loads))
			totalStoreAttempts.Add(int64(storeAttempts))
			totalStores.Add(int64(len(myKeys)))
			totalRanges.Add(int64(ranges))
		}()
	}

	wg.Wait()

	expected := make(map[metricKey]*bool)
	for _, oneExpected := range expecteds {
		for k, v := range oneExpected {
			expected[k] = v
		}
	}
	actual := make(map[metricKey]*bool)
	m.Range(func(k metricKey, v *bool) bool {
		actual[k] = v
		return true
	})

	if !maps.Equal(expected, actual) {
		t.Fatal("ending maps did not match")
	}

	t.Logf("totalLoads         = %d", totalLoads.Load())
	t.Logf("totalStoreAttempts = %d", totalStoreAttempts.Load())
	t.Logf("totalStores        = %d", totalStores.Load())
	t.Logf("totalRanges        = %d", totalRanges.Load())
}
