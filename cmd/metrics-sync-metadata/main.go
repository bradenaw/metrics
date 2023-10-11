package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/bradenaw/juniper/xmaps"
	"github.com/bradenaw/metrics"
)

func main() {
	err := main2()
	if err != nil {
		panic(err)
	}
}

var usage = `
metrics-sync-metadata syncs metric metadata to DataDog. All of the methods for creating metric
definitions in github.com/bradenaw/metrics also require entering units and descriptions so that
these live in source alongside the code that emits them.  This program takes that metadata and sends
it to DataDog so that it appears the same way in the UI, including units on graphs.

This program parses metric metadata from stdin in the same format that it's emitted from
metrics.DumpDefs(). Normally, your binary that emits metrics should accept an additional flag that
makes it call metrics.DumpDefs() and exit.  After building for production, run the binary with that
flag, and then pass its output to this program.

For example:

==== in main.go ====================================================================================

	package main

	import (
		"flag"

		"github.com/bradenaw/metrics"
	)

	func main() {
		dumpDefs = flag.Bool(
			"dump-metric-defs",
			false,
			"If set, dumps metric definitions to stdout on startup and then exits.",
		)
		flag.Parse()
		if *dumpDefs {
			err := metrics.DumpDefs()
			if err != nil {
				panic(err)
			}
			return
		}

		// ... The rest of your program that emits metrics.
	}

==== after building ./foo ==========================================================================

	./foo --dump-metric-defs | DD_API_KEY="<DD_API_KEY>" DD_APP_KEY="<DD_APP_KEY>" metrics-sync-metadata
`

func main2() error {
	metricPrefix := flag.String(
		"metric-prefix",
		"", // default
		"If supplied, places this plus a dot before each metric name.",
	)
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprintf(out, "Usage of %s:\n", os.Args[0])
		fmt.Fprint(out, usage)
		fmt.Fprint(out, "\n\n")
		fmt.Fprint(out, "Flags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	var defs map[string]metrics.Metadata
	err = json.Unmarshal(b, &defs)
	if err != nil {
		return err
	}

	ctx := datadog.NewDefaultContext(context.Background())
	configuration := datadog.NewConfiguration()
	// Since we're making a bunch of calls to update metrics, we may end up getting 429'd.
	configuration.RetryConfiguration.EnableRetry = true
	configuration.RetryConfiguration.MaxRetries = 1_000_000
	apiClient := datadog.NewAPIClient(configuration)
	api := datadogV1.NewMetricsApi(apiClient)

	// Check against ListActiveMetrics to avoid accidentally polluting if, for example,
	// --metric-prefix is missing or incorrect.
	resp, _, err := api.ListActiveMetrics(
		ctx,
		time.Now().Add(-30*24*time.Hour).Unix(),
		*datadogV1.NewListActiveMetricsOptionalParameters(),
	)
	if err != nil {
		return err
	}

	metricNames := xmaps.SetFromSlice(resp.Metrics)

	for name, metadata := range defs {
		fullName := name
		if len(*metricPrefix) > 0 {
			fullName = *metricPrefix + "." + name
		}

		if !metricNames.Contains(fullName) {
			fmt.Fprintf(
				os.Stderr,
				"%s in defs, but has not been emitted to datadog in the last 30d, so not updating "+
					"metadata\n",
				fullName,
			)
			continue
		}

		_, r, err := api.UpdateMetricMetadata(ctx, fullName, datadogV1.MetricMetadata{
			Unit:        datadog.PtrString(string(metadata.Unit)),
			Description: &metadata.Description,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when calling `MetricsApi.UpdateMetricMetadata`: %v\n", err)
			fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
			return err
		}
	}
	return nil
}
