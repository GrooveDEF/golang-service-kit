package grpcserver

import (
	"github.com/definancialbr/golang-container-kit/pkg/container"
	"google.golang.org/grpc"
)

func NewGRPCServer(ctn *container.Container) *grpc.Server {
	return grpc.NewServer(grpc.UnaryInterceptor(MetricsInterceptor(ctn.Metrics)))
}
