package zlog

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/zusux/gokit/utils/tracing"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type sLog struct {
	log       *zap.Logger
	prefix    []any
	hasValuer bool
	ctx       context.Context
}

func NewSlog(log *zap.Logger, callerSkip int) *sLog {
	return &sLog{
		log:    log.WithOptions(zap.AddCaller(), zap.AddCallerSkip(callerSkip)),
		ctx:    context.Background(),
		prefix: make([]any, 0),
	}
}
func (s *sLog) Debugf(format string, args ...interface{}) {
	s.log.Sugar().With(s.prefix...).Debugf(format, args...)
}

func (s *sLog) Printf(format string, args ...interface{}) {
	s.log.Sugar().Infof(format, args...)
}

func (s *sLog) Infof(format string, args ...interface{}) {
	s.log.Sugar().Infof(format, args...)
}
func (s *sLog) Warnf(format string, args ...interface{}) {
	s.log.Sugar().Warnf(format, args...)
}
func (s *sLog) Errorf(format string, args ...interface{}) {
	s.log.Sugar().Errorf(format, args...)
}
func (s *sLog) Panicf(format string, args ...interface{}) {
	s.log.Sugar().Panicf(format, args...)
}
func (s *sLog) Fatalf(format string, args ...interface{}) {
	s.log.Sugar().Fatalf(format, args...)
}
func (s *sLog) Log(level log.Level, keyvals ...any) error {
	s.log.Sugar().Log(zapcore.Level(level), keyvals...)
	return nil
}

// Context returns a shallow copy of l with its context changed
// to ctx. The provided ctx must be non-nil.
func Context(ctx context.Context) *sLog {
	name := tracing.TraceIdName
	traceId := tracing.GetTraceIDFromCtx(ctx)
	return &sLog{
		log:    zap.L().WithOptions(zap.AddCaller(), zap.AddCallerSkip(2)),
		ctx:    ctx,
		prefix: []any{name, traceId},
	}
}
