package main

import (
	"github.com/bradenaw/metrics"
)

var (
	runCount = metrics.NewCounterDef(
		"runs",
		metrics.UnitRun,
		"logged every time this process is started",
	)
	runningGauge = metrics.NewGaugeDef(
		"running",
		metrics.UnitProcess,
		"a gauge set to 1 while this process is running, in aggregate shows the total number that "+
			"are running",
	)
	functionCallCount = metrics.NewCounterDef1[string](
		"function_calls",
		metrics.UnitEvent,
		"name", // the key of the tag
		"counts the number of time each function is called",
	)
)
