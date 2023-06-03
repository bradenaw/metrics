package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/bradenaw/metrics"
)

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

	m := metrics.New(&statsd.NoOpClient{})
	defer m.Flush()

	runCount.Bind(m).Add(1)
	runningGauge.Bind(m).Set(1)
}
