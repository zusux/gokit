package logz

import (
	"fmt"
	"path"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/zusux/gokit/utils"
)

type Logger struct {
	Path            string
	File            string
	Age             int64
	Rotation        int64
	TimestampFormat string
	LoggerLevel     logrus.Level
	LoggerFormat    logrus.Formatter
}
type OptionFunc func(*Logger)

func WithPath(path string) OptionFunc {
	return func(l *Logger) {
		l.Path = path
	}
}
func WithFile(file string) OptionFunc {
	return func(l *Logger) {
		l.File = file
	}
}
func WithAge(age int64) OptionFunc {
	return func(l *Logger) {
		l.Age = age
	}
}
func WithRotation(rotation int64) OptionFunc {
	return func(l *Logger) {
		l.Rotation = rotation
	}
}
func WithLogLevel(level logrus.Level) OptionFunc {
	return func(l *Logger) {
		l.LoggerLevel = level
	}
}
func WithTimestampFormat(timestampFormat string) OptionFunc {
	return func(l *Logger) {
		l.TimestampFormat = timestampFormat
	}
}
func WithFormatJson(timestampFormat string) OptionFunc {
	return func(l *Logger) {
		l.LoggerFormat = &logrus.JSONFormatter{
			TimestampFormat: timestampFormat,
		}
	}
}
func WithFormatText(timestampFormat string) OptionFunc {
	return func(l *Logger) {
		l.LoggerFormat = &logrus.TextFormatter{
			TimestampFormat: timestampFormat,
		}
	}
}

func MustLogger(options ...OptionFunc) *logrus.Logger {
	logger := &Logger{
		Path:            "logs",
		File:            "log",
		Age:             24,
		Rotation:        24,
		TimestampFormat: "2006-01-02 15:04:05",
	}
	for _, opt := range options {
		opt(logger)
	}
	//默认格式化
	if logger.LoggerFormat == nil {
		logger.LoggerFormat = &logrus.TextFormatter{
			TimestampFormat: logger.TimestampFormat,
		}
	}
	log := logrus.New()
	baseLogPath := logger.getFilePath()
	writer, err := rotatelogs.New(
		baseLogPath+".%Y-%m-%d",
		rotatelogs.WithLinkName(baseLogPath),                                  // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Duration(logger.Age)*time.Hour),            // 文件最大保存时间
		rotatelogs.WithRotationTime(time.Duration(logger.Rotation)*time.Hour), // 日志切割时间间隔
	)
	if err != nil {
		panic(errors.New(fmt.Sprintf("[logz][panic] log path is not available err:%v", err.Error())))
	}
	levelMap := make(lfshook.WriterMap, 0)
	for i := logrus.Level(0); i <= logger.LoggerLevel; i++ {
		levelMap[i] = writer
	}
	lfHook := lfshook.NewHook(levelMap, logger.LoggerFormat)
	log.AddHook(lfHook)
	return log
}

func (l *Logger) getFilePath() string {
	baseLogPath := path.Join(l.Path, l.File)
	if !l.checkOrCreateDIr(baseLogPath, "config") {
		workingDir := utils.GetWdDir()
		if !l.checkOrCreateDIr(workingDir, "working") {
			panic(errors.New("[logz][panic] log path is not available"))
		}
		l.Path = workingDir
		baseLogPath = path.Join(l.Path, l.File)
		fmt.Printf("[logz][warning] switch log path to working dir: %s \r\n", workingDir)
	}
	return baseLogPath
}

func (l *Logger) checkOrCreateDIr(baseLogPaht string, where string) bool {
	ok, err := utils.AvailablePath(baseLogPaht)
	if err == nil && ok {
		return true
	}
	fmt.Printf("[logz][warning] %s path is not available, exsit: %v, err: %v \r\n", where, ok, err)
	return false
}
