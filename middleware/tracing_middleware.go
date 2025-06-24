package middleware

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/zusux/gokit/utils/tracing"
)

func Tracing() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			newCtx, tid := tracing.InjectTrace(ctx)
			// 设置 trace_id 到响应头
			if tr, ok := transport.FromServerContext(ctx); ok {
				tr.ReplyHeader().Set("trace_id", tid.String())
			}
			return handler(newCtx, req)
		}
	}
}
