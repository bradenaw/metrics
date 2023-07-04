package metrics_test

import (
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
