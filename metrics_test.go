package metrics

import (
	"testing"
)

var (
	testGauge = NewGaugeDef2[string, string](
		"my.metric.name",
		UnitByte,
		[...]string{"foo", "bar"},
		"this is my description",
	)
)

func TestNothing(t *testing.T) {
	testGauge.Bind(nil, "baz", "qux").Set(1)
}
