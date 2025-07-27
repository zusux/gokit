package zlog

import "go.uber.org/zap"

type Field = zap.Field

// Debug
func Debug(msg string, fields ...Field) {
	zap.L().WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).Debug(msg, fields...)
}
func Debugf(format string, args ...interface{}) {
	zap.S().WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).Debugf(format, args...)
}

// Info
func Info(msg string, fields ...Field) {
	zap.L().WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).Info(msg, fields...)
}
func Infof(format string, args ...interface{}) {
	zap.S().WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).Infof(format, args...)
}

// Warn
func Warn(msg string, fields ...Field) {
	zap.L().WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).Warn(msg, fields...)
}
func Warnf(format string, args ...interface{}) {
	zap.S().WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).Warnf(format, args...)
}

// Error
func Error(msg string, fields ...Field) {
	zap.L().WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).Error(msg, fields...)
}
func Errorf(format string, args ...interface{}) {
	zap.S().WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).Errorf(format, args...)
}

// DPanic
func DPanic(msg string, fields ...Field) {
	zap.L().WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).DPanic(msg, fields...)
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func Panic(msg string, fields ...Field) {
	zap.L().WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).Panic(msg, fields...)
}
func Panicf(format string, args ...interface{}) {
	zap.S().WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).Panicf(format, args...)
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func Fatal(msg string, fields ...Field) {
	zap.L().WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).Fatal(msg, fields...)
}
func Fatalf(format string, args ...interface{}) {
	zap.S().WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).Fatalf(format, args...)
}

// Sync calls the underlying Core's Sync method, flushing any buffered log
// entries. Applications should take care to call Sync before exiting.
func Sync() error {
	return zap.L().WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).Sync()
}

func Print(format string, fields ...Field) {
	zap.L().WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).Info(format, fields...)
}
func Printf(format string, args ...interface{}) {
	zap.S().WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).Infof(format, args...)
}
