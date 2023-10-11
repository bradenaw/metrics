package runtime_metrics

import (
	"github.com/bradenaw/metrics"
)

var (
	gaugeDef = metrics.NewGaugeDef1[string](
		"go.runtime.metrics.num",
		"Export of Go's runtime/metrics package.",
		metrics.NoUnits,
		[...]string{"metric_name"},
	)
	bucketedGaugeDef = metrics.NewGaugeDef2[string, string](
		"go.runtime.gomaxprocs",
		"Export of Go's runtime.GOMAXPROCS().",
		metrics.NoUnits,
		[...]string{"metric_name", "bucket"},
	)
)
