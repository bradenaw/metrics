package runtime_metrics

import (
	"regexp"
	gometrics "runtime/metrics"
	"strings"

	"github.com/bradenaw/metrics"
)

var (
	// Lifted from https://pkg.go.dev/runtime/metrics@go1.21.3#Description.
	nameRegexp = regexp.MustCompile("^(?P<name>/[^:]+):(?P<unit>[^:*/]+(?:[*/][^:*/]+)*)$")

	// Slightly-manual reproduction of the list in
	// https://pkg.go.dev/runtime/metrics@go1.21.3#hdr-Supported_metrics
	unitsStrToUnit = map[string]metrics.Unit{
		"bytes":   metrics.UnitByte,
		"threads": metrics.UnitThread,
		"seconds": metrics.UnitSecond,
		"events":  metrics.UnitEvent,
		// gc-cycle is a counter incremented every GC, gc-cycles is a number of garbage collections.
		"gc-cycles":   metrics.UnitGarbageCollection,
		"objects":     metrics.UnitObject,
		"percent":     metrics.UnitPercent,
		"cpu-seconds": metrics.UnitSecond,
		"calls":       metrics.UnitEvent,
	}

	descriptions []gometrics.Description
	// With the same indexes as descriptions, KindFloat64 and KindUint64 are populated.
	gaugeDefs []metrics.GaugeDef
	// With the same indexes as descriptions, KindFloat64Histogram are populated.
	bucketedGaugeDefs []metrics.GaugeDef1[string]
)

func init() {
	descriptions = gometrics.All()
	gaugeDefs = make([]metrics.GaugeDef, len(descriptions))
	bucketedGaugeDefs = make([]metrics.GaugeDef1[string], len(descriptions))

	for i, description := range descriptions {
		submatches := nameRegexp.FindStringSubmatch(description.Name)
		cleanName := submatches[1]
		cleanName = strings.Trim(cleanName, "/")
		unitsStr := submatches[2]
		// Include units in the name because there are duplicates of just the name part, e.g.
		//   gc/heap/allocs:bytes
		//   gc/heap/allocs:objects
		name := sanitizeName("go.runtime." + cleanName + "_" + unitsStr)
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

var disallowedNameChars = regexp.MustCompile("[^a-zA-Z0-9_.]")

// The actual allowed regex is ^[a-z][a-zA-Z0-9_.]{0,199}$
func sanitizeName(s string) string {
	if len(s) == 0 {
		return s
	}
	s = disallowedNameChars.ReplaceAllString(s, "_")
	if !(s[0] >= 'a' && s[0] <= 'z') {
		s = "a" + s
	}
	if len(s) > 199 {
		s = s[:199]
	}
	return s
}
