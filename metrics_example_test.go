package metrics_test

import (
	"fmt"

	"github.com/bradenaw/metrics"
)

func ExampleExponentialBuckets() {
	buckets := metrics.ExponentialBuckets(100, 10, 3)
	fmt.Println("ExponentialBuckets(100, 10, 3) ->", buckets)

	buckets = metrics.ExponentialBuckets(100, 2, 5)
	fmt.Println("ExponentialBuckets(100, 2, 5) ->", buckets)

	// Output:
	// ExponentialBuckets(100, 10, 3) -> [100 1000 10000]
	// ExponentialBuckets(100, 2, 5) -> [100 200 400 800 1600]
}

func ExampleLinearBuckets() {
	buckets := metrics.LinearBuckets(100, 50, 3)
	fmt.Println("LinearBuckets(100, 50, 3) ->", buckets)

	buckets = metrics.LinearBuckets(100, 75, 4)
	fmt.Println("LinearBuckets(100, 75, 4) ->", buckets)

	// Output:
	// LinearBuckets(100, 50, 3) -> [100 150 200]
	// LinearBuckets(100, 75, 4) -> [100 175 250 325]
}
