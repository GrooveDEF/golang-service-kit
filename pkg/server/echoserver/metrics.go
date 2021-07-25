package echoserver

import (
	"reflect"
	"strconv"
	"time"

	"github.com/GrooveDEF/golang-container-kit/pkg/metrics"
	"github.com/GrooveDEF/golang-container-kit/pkg/metrics/prometheus"
	echo "github.com/labstack/echo/v4"
)

const (
	httpSubsystem        = "http"
	httpNamespace        = "echo"
	httpRequestsCount    = "requests_total"
	httpRequestsDuration = "request_duration_seconds"
	notFoundPath         = "/not-found"
)

var (
	durationBuckets = []float64{
		0.0005,
		0.001, // 1ms
		0.002,
		0.005,
		0.01, // 10ms
		0.02,
		0.05,
		0.1, // 100 ms
		0.2,
		0.5,
		1.0, // 1s
		2.0,
		5.0,
		10.0, // 10s
		15.0,
		20.0,
		30.0,
	}
)

func isNotFoundHandler(handler echo.HandlerFunc) bool {
	return reflect.ValueOf(handler).Pointer() == reflect.ValueOf(echo.NotFoundHandler).Pointer()
}

// MetricsMiddleware returns an echo middleware for instrumentation.
func MetricsMiddleware(metrics metrics.MetricService) echo.MiddlewareFunc {

	httpRequests := metrics.Counter(
		prometheus.WithNamespace(httpNamespace),
		prometheus.WithName(httpRequestsCount),
		prometheus.WithSubsystem(httpSubsystem),
		prometheus.WithHelp("Number of HTTP operations"),
		prometheus.WithLabels([]string{"status", "method", "handler"}),
	)

	httpDuration := metrics.Histogram(
		prometheus.WithNamespace(httpNamespace),
		prometheus.WithName(httpRequestsDuration),
		prometheus.WithSubsystem(httpSubsystem),
		prometheus.WithBuckets(durationBuckets),
		prometheus.WithHelp("Spend time by processing a route"),
		prometheus.WithLabels([]string{"method", "handler"}),
	)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			req := ctx.Request()
			path := ctx.Path()

			// to avoid attack high cardinality of 404
			if isNotFoundHandler(ctx.Handler()) {
				path = notFoundPath
			}

			start := time.Now()

			requestDuration := httpDuration.WithLabelValues(req.Method, path)
			err := next(ctx)
			requestDuration.Observe(time.Since(start).Seconds())

			if err != nil {
				ctx.Error(err)
			}

			status := strconv.Itoa(ctx.Response().Status)

			httpRequests.WithLabelValues(status, req.Method, path).Add(1)

			return err
		}
	}
}
