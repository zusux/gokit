package zlog

import (
	"fmt"
	"go.uber.org/zap"
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

// logEncoder
type logEncoder struct {
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
func (e *logEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
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

type Logger struct {
	Path            string
	File            string
	Age             int64
	Rotation        int64
	TimestampFormat string // 2006-01-02 15:04:05
	LoggerLevel     uint8  // -1:debug 0:info 1:warn 2:error 3:DPanic 4:Panic 5:fatal
	LoggerFormat    string //json
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) InitLog() *zap.Logger {
	// 使用 zap 的 NewDevelopmentConfig 快速配置
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.Level(l.LoggerLevel)),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         l.LoggerFormat,
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(l.TimestampFormat) // 替换时间格式化方式
	cfg.EncoderConfig.EncodeLevel = ColourEncodeLevel
	// 创建自定义的 Encoder
	encoder := &logEncoder{
		Encoder: zapcore.NewConsoleEncoder(cfg.EncoderConfig), // 使用 Console 编码器
	}
	// 创建 Core
	core := zapcore.NewCore(
		encoder,                      // 使用自定义的 Encoder
		zapcore.AddSync(os.Stdout),   // 输出到控制台
		zapcore.Level(l.LoggerLevel), // 设置日志级别
	)
	// 创建 Logger
	logger := zap.New(core, zap.AddCaller())

	zap.ReplaceGlobals(logger)
	return logger
}
