package grpc

import (
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/zusux/gokit/gserver/invoke/api"
	"github.com/zusux/gokit/gserver/invoke/isrv"
	"github.com/zusux/gokit/middleware"
	"time"
)

func GServer(addr string, timeout time.Duration, extOpts ...grpc.ServerOption) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			middleware.Tracing(),
		),
	}

	if addr != "" {
		opts = append(opts, grpc.Address(addr))
	}
	if timeout != 0 {
		opts = append(opts, grpc.Timeout(timeout))
	}
	for _, v := range extOpts {
		opts = append(opts, v)
	}
	srv := grpc.NewServer(opts...)
	api.RegisterInternalRouterServer(srv, isrv.NewInvokeService())
	return srv
}
