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
)
