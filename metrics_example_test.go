package metrics_test

import (
	"fmt"
	"sync"

	"github.com/bradenaw/metrics"
)

func ExampleBucketedGaugeGroup() {
	// --- metrics.go ------------------------------------------------------------------------------

	var (
		streamsPerConnectionDef = metrics.NewGaugeDef1[string](
			"connections_by_active_streams",
			"The number of open TCP connections by the number of active streams on each.",
			metrics.Unit(""),
			[...]string{"bucket"}, // key
		)
	)

	// --- server.go -------------------------------------------------------------------------------

	type Connection struct {
		// ...
	}
	type Stream struct {
		// ...
	}

	type Server struct {
		// Should be called from Server.Close().
		stopEmittingMetrics func()

		m           sync.Mutex
		connections map[*Connection][]*Stream
	}

	// Imagine this is func NewServer():
	_ = func(mtr *metrics.Metrics) *Server {
		s := &Server{ /* ... */ }
		streamsPerConnection := metrics.NewBucketedGaugeGroup(
			mtr,
			streamsPerConnectionDef,
			// These are the boundaries between buckets. This causes these gauges to be made:
			//   connections_by_active_streams   bucket:lt_1
			//   connections_by_active_streams   bucket:gte_1_lt_10
			//   connections_by_active_streams   bucket:gte_10_lt_100
			//   connections_by_active_streams   bucket:gte_100
			[]float64{1, 10, 100},
		)

		// EveryFlush will be called at least once per 10-second metrics interval. It is called
		// infrequently so it's acceptable to do slightly-expensive metrics computation.
		s.stopEmittingMetrics = mtr.EveryFlush(func() {
			s.m.Lock()
			defer s.m.Unlock()
			for _, streams := range s.connections {
				streamsPerConnection.Observe(float64(len(streams)))
			}
			streamsPerConnection.Emit()
		})

		return s
	}
}

func ExampleExponentialBuckets() {
	buckets := metrics.ExponentialBuckets(100, 10, 3)
	fmt.Println("ExponentialBuckets(100, 10, 3) ->", buckets)

	buckets = metrics.ExponentialBuckets(100, 2, 5)
	fmt.Println("ExponentialBuckets(100, 2, 5) ->", buckets)

	// Output:
	// ExponentialBuckets(100, 10, 3) -> [100 1000 10000]
	// ExponentialBuckets(100, 2, 5) -> [100 200 400 800 1600]
}

func ExampleLinearBuckets() {
	buckets := metrics.LinearBuckets(100, 50, 3)
	fmt.Println("LinearBuckets(100, 50, 3) ->", buckets)

	buckets = metrics.LinearBuckets(100, 75, 4)
	fmt.Println("LinearBuckets(100, 75, 4) ->", buckets)

	// Output:
	// LinearBuckets(100, 50, 3) -> [100 150 200]
	// LinearBuckets(100, 75, 4) -> [100 175 250 325]
}
