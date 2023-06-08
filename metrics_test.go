package metrics

import (
	"testing"
)

var (
	testGauge = NewGaugeDef2[string, string](
		"my.metric.name",
		"this is my description",
		UnitByte,
		[...]string{"foo", "bar"},
	)
)

func TestNothing(t *testing.T) {
	testGauge.Bind(NoOpMetrics, "baz", "qux").Set(1)
}
