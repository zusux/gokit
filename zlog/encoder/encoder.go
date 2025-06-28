package encoder

import (
	"fmt"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
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
func (e *LogEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	// 先调用原始的 EncodeEntry 方法生成日志行
	buff, err := e.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		return nil, err
	}
	data := buff.String()
	buff.Reset()
	buff.AppendString("[myApp] " + data)
	data = buff.String()
	// 时间分片
	now := time.Now().Format("2006-01-02")
	if e.currentDate != now {
		os.MkdirAll(fmt.Sprintf("logs/%s", now), 0666)
		// 时间不同，先创建目录
		name := fmt.Sprintf("logs/%s/out.log", now)
		file, _ := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		e.file = file
		e.currentDate = now
	}

	switch entry.Level {
	case zapcore.ErrorLevel:
		if e.errFile == nil {
			name := fmt.Sprintf("logs/%s/err.log", now)
			file, _ := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
			e.errFile = file
		}
		e.errFile.WriteString(buff.String())
	}

	if e.currentDate == now {
		e.file.WriteString(data)
	}
	return buff, nil
}
