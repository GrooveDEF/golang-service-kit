package grpcserver

import (
	"github.com/definancialbr/golang-container-kit/pkg/container"
	"google.golang.org/grpc"
)

func NewGRPCServer(ctn *container.Container, serverOptions ...grpc.ServerOption) *grpc.Server {
	options := append(serverOptions, grpc.UnaryInterceptor(MetricsInterceptor(ctn.Metrics)))
	return grpc.NewServer(options...)
}
