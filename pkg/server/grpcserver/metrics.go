package grpcserver

import (
	"context"
	"time"

	"github.com/definancialbr/golang-container-kit/pkg/metrics"
	"github.com/definancialbr/golang-container-kit/pkg/metrics/prometheus"
	"google.golang.org/grpc"
)

const (
	grpcNamespace              = "grpc"
	grpcUnaryRequestsCount     = "unary_requests_total"
	grpcUnaryRequestsDuration  = "unary_request_duration_seconds"
	grpcStreamRequestsCount    = "stream_requests_total"
	grpcStreamRequestsDuration = "stream_request_duration_seconds"
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

// MetricsMiddleware returns an unary grpc interceptor for instrumentation.
func MetricsInterceptor(metrics metrics.MetricService) grpc.UnaryServerInterceptor {

	grpcRequests := metrics.Counter(
		prometheus.WithName(grpcUnaryRequestsCount),
		prometheus.WithHelp("Number of unary operations"),
		prometheus.WithLabels([]string{"status", "method"}),
	)

	grpcDuration := metrics.Histogram(
		prometheus.WithNamespace(grpcNamespace),
		prometheus.WithName(grpcUnaryRequestsDuration),
		prometheus.WithBuckets(durationBuckets),
		prometheus.WithHelp("Spend time by processing an unary request"),
		prometheus.WithLabels([]string{"method"}),
	)

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		path := info.FullMethod
		start := time.Now()
		requestDuration := grpcDuration.WithLabelValues(path)
		resp, err := handler(ctx, req)
		requestDuration.Observe(time.Since(start).Seconds())
		status := "ok"
		if err != nil {
			status = "err"
		}

		grpcRequests.WithLabelValues(status, path).Add(1)
		return resp, err
	}
}

func MetricsStreamInterceptor(metrics metrics.MetricService) grpc.StreamServerInterceptor {

	grpcRequests := metrics.Counter(
		prometheus.WithName(grpcStreamRequestsCount),
		prometheus.WithHelp("Number of stream operations"),
		prometheus.WithLabels([]string{"status", "method"}),
	)

	grpcDuration := metrics.Histogram(
		prometheus.WithNamespace(grpcNamespace),
		prometheus.WithName(grpcStreamRequestsDuration),
		prometheus.WithBuckets(durationBuckets),
		prometheus.WithHelp("Spend time by processing an stream request"),
		prometheus.WithLabels([]string{"method"}),
	)

	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		path := info.FullMethod
		start := time.Now()
		requestDuration := grpcDuration.WithLabelValues(path)
		err := handler(srv, ss)
		requestDuration.Observe(time.Since(start).Seconds())
		status := "ok"
		if err != nil {
			status = "err"
		}

		grpcRequests.WithLabelValues(status, path).Add(1)
		return err
	}
}
