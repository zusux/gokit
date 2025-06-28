package mlog

import (
	"github.com/zusux/gokit/zlog"
	"gorm.io/gorm/logger"
	"time"
)

func GetLogger() logger.Interface {
	return logger.New(
		zlog.NewSLog(), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)
}
