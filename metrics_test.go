package metrics

import (
	"testing"
)

var (
	testDef = NewGaugeDef2[string, string](
		"my.metric.name",
		"this is my description",
		UnitByte,
		[...]string{"foo", "bar"},
	)
)

func TestNothing(t *testing.T) {
	NoOpMetrics.Gauge(testDef.Values("baz", "qux")).Set(1)
}
