package travel

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	v1 "github.com/webmin7761/go-school/homework/final/api/travel/v1"
	"github.com/webmin7761/go-school/homework/final/internal/conf"
	"github.com/webmin7761/go-school/homework/final/internal/service"
	"go.opentelemetry.io/otel/trace"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, tracer trace.TracerProvider, travel *service.TravelService) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			middleware.Chain(
				tracing.Server(tracing.WithTracerProvider(tracer)),
				logging.Server(log.DefaultLogger),
				recovery.Recovery(),
			),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterTravelServiceServer(srv, travel)
	return srv
}
