package runtime_metrics

import (
	"fmt"
	"math"
	gometrics "runtime/metrics"

	"github.com/bradenaw/metrics"
)

func Emit(m *metrics.Metrics) {
	descriptions := gometrics.All()
	samples := make([]gometrics.Sample, len(descriptions))
	gauges := make([][]*metrics.Gauge, len(descriptions))

	for i, description := range descriptions {
		samples[i].Name = description.Name
	}
	gometrics.Read(samples)
	for i, description := range descriptions {
		if description.Kind == gometrics.KindFloat64Histogram {
			bucketNames := makeBucketNames(samples[i].Value.Float64Histogram().Buckets)
			gauges[i] = make([]*metrics.Gauge, 0, len(bucketNames))
			for _, bucketName := range bucketNames {
				gauges[i] = append(gauges[i], m.Gauge(bucketedGaugeDefs[i].Values(bucketName)))
			}
		} else {
			gauges[i] = []*metrics.Gauge{m.Gauge(gaugeDefs[i])}
		}
	}

	m.EveryFlush(func() {
		gometrics.Read(samples)

		for i, sample := range samples {
			description := descriptions[i]
			switch description.Kind {
			case gometrics.KindFloat64:
				v := sample.Value.Float64()
				gauges[i][0].Set(v)
			case gometrics.KindUint64:
				v := sample.Value.Uint64()
				gauges[i][0].Set(float64(v))
			case gometrics.KindFloat64Histogram:
				v := sample.Value.Float64Histogram()
				for j := range v.Counts {
					gauges[i][j].Set(float64(v.Counts[i]))
				}
			}
		}
	})
}

// similar to metrics.bucketNames, but slightly different. runtime/metrics's buckets are all lower
// bounds, and they include -Inf and +Inf. Attempt to make the same style of names as
// BucketedCounter/BucketedGaugeGroup.
func makeBucketNames(lowerBounds []float64) []string {
	names := make([]string, len(lowerBounds))
	for i := range lowerBounds {
		unboundedLower := math.IsInf(lowerBounds[0], -1)
		unboundedUpper := i == len(lowerBounds)-1 || math.IsInf(lowerBounds[i+1], 1)
		if unboundedLower && unboundedUpper {
			names[i] = ""
		} else if unboundedLower {
			names[i] = fmt.Sprintf("lt_%g", lowerBounds[i+1])
		} else if unboundedUpper {
			names[i] = fmt.Sprintf("gte_%g", lowerBounds[i])
		} else {
			names[i] = fmt.Sprintf("gte_%g_lt_%g", lowerBounds[i], lowerBounds[i+1])
		}
	}
	return names
}
