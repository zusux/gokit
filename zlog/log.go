package zlog

import "go.uber.org/zap"

type Field = zap.Field

// Debug
func Debug(msg string, fields ...Field) {
	zap.L().Debug(msg, fields...)
}
func Debugf(format string, args ...interface{}) {
	zap.S().Debugf(format, args...)
}

// Info
func Info(msg string, fields ...Field) {
	zap.L().Info(msg, fields...)
}
func Infof(format string, args ...interface{}) {
	zap.S().Infof(format, args...)
}

// Warn
func Warn(msg string, fields ...Field) {
	zap.L().Warn(msg, fields...)
}
func Warnf(format string, args ...interface{}) {
	zap.S().Warnf(format, args...)
}

// Error
func Error(msg string, fields ...Field) {
	zap.L().Error(msg, fields...)
}
func Errorf(format string, args ...interface{}) {
	zap.S().Errorf(format, args...)
}

// DPanic
func DPanic(msg string, fields ...Field) {
	zap.L().DPanic(msg, fields...)
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func Panic(msg string, fields ...Field) {
	zap.L().Panic(msg, fields...)
}
func Panicf(format string, args ...interface{}) {
	zap.S().Panicf(format, args...)
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func Fatal(msg string, fields ...Field) {
	zap.L().Fatal(msg, fields...)
}
func Fatalf(format string, args ...interface{}) {
	zap.S().Fatalf(format, args...)
}

// Sync calls the underlying Core's Sync method, flushing any buffered log
// entries. Applications should take care to call Sync before exiting.
func Sync() error {
	return zap.L().Sync()
}

func Print(format string, fields ...Field) {
	zap.L().Info(format, fields...)
}
func Printf(format string, args ...interface{}) {
	zap.S().Infof(format, args...)
}
