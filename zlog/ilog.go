package zlog

import "go.uber.org/zap"

type Field = zap.Field

type ILogger interface {
	// Debug
	Debug(msg string, fields ...Field)
	Debugf(format string, args ...interface{})
	// Info
	Info(msg string, fields ...Field)
	Infof(format string, args ...interface{})

	//Warn
	Warn(msg string, fields ...Field)
	Warnf(format string, args ...interface{})
	//Error
	Error(msg string, fields ...Field)
	Errorf(format string, args ...interface{})
	// DPanic
	DPanic(msg string, fields ...Field)

	// Panic logs a message at PanicLevel. The message includes any fields passed
	// at the log site, as well as any fields accumulated on the logger.
	//
	// The logger then panics, even if logging at PanicLevel is disabled.
	Panic(msg string, fields ...Field)
	Panicf(format string, args ...interface{})
	// Fatal logs a message at FatalLevel. The message includes any fields passed
	// at the log site, as well as any fields accumulated on the logger.
	//
	// The logger then calls os.Exit(1), even if logging at FatalLevel is
	// disabled.
	Fatal(msg string, fields ...Field)
	Fatalf(format string, args ...interface{})
	// Sync calls the underlying Core's Sync method, flushing any buffered log
	// entries. Applications should take care to call Sync before exiting.
	Sync() error

	Print(format string, fields ...Field)
	Printf(format string, args ...interface{})
}
