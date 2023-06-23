package main

import (
	"github.com/bradenaw/metrics"
)

var (
	runCount = metrics.NewCounterDef(
		"runs",
		"logged every time this process is started",
		metrics.UnitRun,
	)

	runningGauge = metrics.NewGaugeDef(
		"running",
		"a gauge set to 1 while this process is running, in aggregate shows the total number that "+
			"are running",
		metrics.UnitProcess,
	)

	// Defines a counter with one tag, whose key is "name" and whose value is a string.
	functionCallCount = metrics.NewCounterDef1[string](
		"function_calls",
		"counts the number of time each function is called",
		metrics.UnitEvent,
		"name",
	)
)
