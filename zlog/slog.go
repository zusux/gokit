package zlog

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type sLog struct {
	prefix    []any
	hasValuer bool
	ctx       context.Context
}

func NewSlog() *sLog {
	return &sLog{
		ctx:    context.Background(),
		prefix: make([]any, 0),
	}
}

func (s *sLog) Debugf(format string, args ...interface{}) {
	zap.S().With().Debugf(format, args...)
}

func (s *sLog) Printf(format string, args ...interface{}) {
	Printf(format, args...)
}

func (s *sLog) Infof(format string, args ...interface{}) {
	zap.S().Infof(format, args...)
}
func (s *sLog) Warnf(format string, args ...interface{}) {
	zap.S().Warnf(format, args...)
}
func (s *sLog) Errorf(format string, args ...interface{}) {
	zap.S().Errorf(format, args...)
}
func (s *sLog) Panicf(format string, args ...interface{}) {
	zap.S().Panicf(format, args...)
}
func (s *sLog) Fatalf(format string, args ...interface{}) {
	zap.S().Fatalf(format, args...)
}
func (s *sLog) Log(level log.Level, keyvals ...any) error {
	zap.S().Log(zapcore.Level(level), keyvals...)
	return nil
}

// Context returns a shallow copy of l with its context changed
// to ctx. The provided ctx must be non-nil.
func Context(ctx context.Context) *sLog {
	return &sLog{
		ctx:    ctx,
		prefix: make([]any, 0),
	}
}
