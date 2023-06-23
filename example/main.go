package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bradenaw/metrics"
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
	var showMetricsMetadata bool
	flag.BoolVar(
		&showMetricsMetadata,
		"metrics-metadata-json",
		false,
		"If set, prints the metadata for the metrics reported by this process in the format "+
			"accepted by https://docs.datadoghq.com/api/latest/metrics/#edit-metric-metadata and "+
			"then exits.",
	)
	flag.Parse()

	if showMetricsMetadata {
		fmt.Println(metrics.FormatMetadataJSON())
		os.Exit(0)
	}

	// In reality, would be
	//   metrics.New(statsd.New("addr-here", statsd.WithoutClientSideAggregation()))
	m := metrics.NoOpMetrics
	defer m.Flush()

	runCount.Bind(m).Add(1)
	runningGauge.Bind(m).Set(1)

	f := &frobber{
		// logs as function_calls with tag name:Foo
		fooCalls: functionCallCount.Values("Foo").Bind(m),
		// logs as function_calls with tag name:Bar
		barCalls: functionCallCount.Values("Bar").Bind(m),
	}
	f.Foo()
	f.Bar()
	f.Bar()
}
