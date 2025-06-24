package tracing

import (
	"context"
	cRand "crypto/rand"

	"math/rand"
	"time"

	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport"
	"go.opentelemetry.io/otel/trace"
)

var globalRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func InjectTrace(ctx context.Context) (context.Context, trace.TraceID) {
	// 1. 优先从 ctx 获取
	if tid, ok := tracing.TraceID()(ctx).(trace.TraceID); ok && tid.IsValid() {
		return ctx, tid
	}

	// 2. 检查外部请求 trace_id
	var traceID string
	if tr, ok := transport.FromServerContext(ctx); ok {
		traceID = tr.RequestHeader().Get("trace_id")
	}

	// 3. 校验透传 trace_id
	tid, err := trace.TraceIDFromHex(traceID)
	if err != nil || !tid.IsValid() {
		tid = trace.TraceID(generateSecureBytes(16))
	}

	// 4. 新 spanID
	spanID := trace.SpanID(generateSecureBytes(8))

	// 5. 注入新上下文
	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: tid,
		SpanID:  spanID,
	})
	newCtx := trace.ContextWithSpanContext(ctx, sc)

	return newCtx, tid
}

func generateSecureBytes(length int) []byte {
	b := make([]byte, length)
	if _, err := cRand.Read(b); err != nil {
		for i := 0; i < length; i++ {
			b[i] = byte(globalRand.Intn(256))
		}
	}
	return b
}
