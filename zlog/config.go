package zlog

import (
	"dario.cat/mergo"
	"github.com/zusux/gokit/zlog/encoder"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Logger struct {
	Path            string `yaml:"path" json:"path"`
	File            string `yaml:"file" json:"file"`
	Age             int64  `yaml:"age" json:"age"`
	Rotation        int64  `yaml:"rotation" json:"rotation"`
	TimestampFormat string `yaml:"timestamp_format" json:"timestamp_format"` // 2006-01-02 15:04:05
	LoggerLevel     uint8  `yaml:"logger_level" json:"logger_level"`         // -1:debug 0:info 1:warn 2:error 3:DPanic 4:Panic 5:fatal
	LoggerFormat    string `yaml:"logger_format" json:"logger_format"`       //json
}

func (l *Logger) InitLog() *zap.Logger {
	defaultConf := Logger{
		Path:            "logs",
		File:            "out.log",
		Age:             7,
		Rotation:        24,
		TimestampFormat: "2006-01-02 15:04:05",
		LoggerLevel:     0,
		LoggerFormat:    "console",
	}
	err := mergo.Merge(l, defaultConf)
	if err != nil {
		panic(err)
	}
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
	cfg.EncoderConfig.EncodeLevel = encoder.ColourEncodeLevel
	// 创建自定义的 Encoder
	encoderObj := &encoder.LogEncoder{
		Encoder: zapcore.NewConsoleEncoder(cfg.EncoderConfig), // 使用 Console 编码器
	}
	// 创建 Core
	core := zapcore.NewCore(
		encoderObj,                   // 使用自定义的 Encoder
		zapcore.AddSync(os.Stdout),   // 输出到控制台
		zapcore.Level(l.LoggerLevel), // 设置日志级别
	)
	// 创建 Logger
	logger := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(logger)
	return logger
}
