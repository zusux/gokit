package zlog

import (
	"dario.cat/mergo"
	"github.com/zusux/gokit/zlog/encoder"
	"github.com/zusux/gokit/zlog/writer"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type Logger struct {
	App             string `yaml:"app" json:"app"`
	WriteConsole    bool   `yaml:"write_console" json:"write_console"`
	WriteFile       bool   `yaml:"write_file" json:"write_file"`
	Path            string `yaml:"path" json:"path"`
	File            string `yaml:"file" json:"file"`
	ErrFile         string `yaml:"err_file" json:"err_file"`
	Age             int64  `yaml:"age" json:"age"`
	Rotation        int64  `yaml:"rotation" json:"rotation"`
	TimestampFormat string `yaml:"timestamp_format" json:"timestamp_format"` // 2006-01-02 15:04:05
	LoggerLevel     uint8  `yaml:"logger_level" json:"logger_level"`         // -1:debug 0:info 1:warn 2:error 3:DPanic 4:Panic 5:fatal
	LoggerFormat    string `yaml:"logger_format" json:"logger_format"`       //json
}

func (l *Logger) InitLog() *zap.Logger {
	defaultConf := Logger{
		App:             "",
		WriteConsole:    true,
		WriteFile:       true,
		Path:            "logs",
		File:            "out.log",
		ErrFile:         "err.log",
		Age:             7,
		Rotation:        24,
		TimestampFormat: "2006-01-02T15:04:05.000Z",
		LoggerLevel:     0,
		LoggerFormat:    "console",
	}
	err := mergo.Merge(l, defaultConf)
	if err != nil {
		panic(err)
	}

	level := zapcore.Level(l.LoggerLevel)
	var cores []zapcore.Core

	// ========== 控制台输出 ==========
	consoleEncoder := &encoder.AppPrefixEncoder{
		App:     l.App,
		Format:  l.LoggerFormat,
		Encoder: buildEncoder(l.LoggerFormat, l.TimestampFormat, true),
	}
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level)
	if l.WriteConsole {
		cores = append(cores, consoleCore)
	}

	// ========== 文件输出 ==========
	if l.WriteFile {
		_ = os.MkdirAll(l.Path, os.ModePerm)
		fileEncoder := &encoder.AppPrefixEncoder{
			App:     l.App,
			Format:  l.LoggerFormat,
			Encoder: buildEncoder(l.LoggerFormat, l.TimestampFormat, false),
		}
		splitWrite := writer.NewDateSplitWriter(l.Path, l.File, l.Rotation)
		splitWrite.StartCleaner(l.Age, time.Hour) // ✅ 启动后台定时清理任务
		fileWriter := zapcore.AddSync(splitWrite)

		fileCore := zapcore.NewCore(fileEncoder, fileWriter, level)
		cores = append(cores, fileCore)

		// 单独的 error 日志输出
		errWriter := zapcore.AddSync(writer.NewDateSplitWriter(l.Path, l.ErrFile, l.Rotation))
		errLevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l >= zapcore.ErrorLevel
		})
		errCore := zapcore.NewCore(fileEncoder, errWriter, errLevel)
		cores = append(cores, errCore)
	}

	if len(cores) == 0 {
		cores = append(cores, consoleCore)
	}
	core := zapcore.NewTee(cores...)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	zap.ReplaceGlobals(logger)
	return logger
}

func buildEncoder(format, tsFormat string, withColor bool) zapcore.Encoder {
	cfg := zapcore.EncoderConfig{
		TimeKey:    "time",
		LevelKey:   "level",
		MessageKey: "msg",
		EncodeTime: zapcore.TimeEncoderOfLayout(tsFormat),
	}

	switch format {
	case "json":
		// JSON 格式不带颜色
		cfg.EncodeLevel = zapcore.CapitalLevelEncoder
		return zapcore.NewJSONEncoder(cfg)
	default:
		// Console 格式带颜色或不带颜色
		if withColor {
			cfg.EncodeLevel = encoder.ColourEncodeLevel
		} else {
			cfg.EncodeLevel = zapcore.CapitalLevelEncoder
		}
		return zapcore.NewConsoleEncoder(cfg)
	}
}
