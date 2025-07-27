package encoder

import (
	"go.uber.org/zap/zapcore"
	"os"
)

const (
	GreenColor  = "\033[32m"   //info
	BlueColor   = "\033[34;1m" //warn
	RedColor    = "\033[31m"   //error
	YellowColor = "\033[33m"
	ResetColor  = "\033[0m"
)

// LogEncoder
type LogEncoder struct {
	Server string
	Method string
	Writer string
	zapcore.Encoder
	errFile     *os.File
	file        *os.File
	currentDate string
}

func ColourEncodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch level {
	case zapcore.InfoLevel:
		enc.AppendString(GreenColor + "INFO" + ResetColor)
	case zapcore.WarnLevel:
		enc.AppendString(BlueColor + "WARN" + ResetColor)
	case zapcore.ErrorLevel, zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		enc.AppendString(RedColor + "ERROR" + ResetColor)
	default:
		enc.AppendString(level.String())
	}
}
