package runtime_metrics

import (
	"regexp"
	gometrics "runtime/metrics"
	"strings"

	"github.com/bradenaw/metrics"
)

var (
	// Lifted from https://pkg.go.dev/runtime/metrics@go1.21.3#Description.
	nameRegexp     = regexp.MustCompile("^(?P<name>/[^:]+):(?P<unit>[^:*/]+(?:[*/][^:*/]+)*)$")
	unitsStrToUnit = map[string]metrics.Unit{
		"bytes":   metrics.UnitByte,
		"threads": metrics.UnitThread,
		"seconds": metrics.UnitSecond,
		"events":  metrics.UnitEvent,
		// gc-cyle is a counter incremented every GC, gc-cycles is a number of garbage collections.
		"gc-cycles":   metrics.UnitGarbageCollection,
		"objects":     metrics.UnitObject,
		"percent":     metrics.UnitPercent,
		"cpu-seconds": metrics.UnitSecond,
		"calls":       metrics.UnitEvent,
	}

	gaugeDefs         []metrics.GaugeDef
	bucketedGaugeDefs []metrics.GaugeDef1[string]
)

func init() {
	descriptions := gometrics.All()
	samples := make([]gometrics.Sample, len(descriptions))
	gaugeDefs = make([]metrics.GaugeDef, len(descriptions))
	bucketedGaugeDefs = make([]metrics.GaugeDef1[string], len(descriptions))

	for i, description := range descriptions {
		samples[i].Name = description.Name
	}
	gometrics.Read(samples)
	for i, description := range descriptions {
		submatches := nameRegexp.FindStringSubmatch(description.Name)
		name := "go.runtime." + strings.ReplaceAll(strings.Trim(submatches[0], "/"), "/", ".")
		unitsStr := submatches[1]
		unit := unitsStrToUnit[unitsStr]

		if description.Kind == gometrics.KindFloat64Histogram {
			bucketedGaugeDefs[i] = metrics.NewGaugeDef1[string](
				name,
				description.Description,
				unit,
				[...]string{"bucket"},
			)
		} else {
			gaugeDefs[i] = metrics.NewGaugeDef(
				name,
				description.Description,
				unit,
			)
		}
	}
}
