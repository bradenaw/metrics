package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bradenaw/metrics"
	"github.com/bradenaw/metrics/runtime_metrics"
)

type frobber struct {
	fooCalls *metrics.Counter
	barCalls *metrics.Counter
}

func (f *frobber) Foo() {
	f.fooCalls.Add(1)
}

func (f *frobber) Bar() {
	f.barCalls.Add(1)
}

func main() {
	var showMetricNames bool
	flag.BoolVar(
		&showMetricNames,
		"show-metrics",
		false,
		"If set, dumps information about all of the metrics defined in this binary and exits.",
	)
	flag.Parse()

	if showMetricNames {
		err := metrics.DumpDefs()
		if err != nil {
			fmt.Println("couldn't dump metrics: ", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// In reality, would be
	//   metrics.New(statsd.New("addr-here", statsd.WithoutClientSideAggregation()))
	m := metrics.NoOpMetrics
	defer m.Flush()
	runtime_metrics.Emit(m)

	m.Counter(runCounterDef).Add(1)
	m.Gauge(runningGaugeDef).Set(1)

	f := &frobber{
		// logs as function_calls with tag name:Foo
		fooCalls: m.Counter(functionCallCounterDef.Values("Foo")),
		// logs as function_calls with tag name:Bar
		barCalls: m.Counter(functionCallCounterDef.Values("Bar")),
	}
	f.Foo()
	f.Bar()
	f.Bar()
}
