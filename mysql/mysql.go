package mysql

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type ClientMysql struct {
	config *ConfigMysql
	log    logger.Interface
}

func NewMysqlClient(config *ConfigMysql, log logger.Interface) *ClientMysql {
	return &ClientMysql{
		config: config,
		log:    log,
	}
}

func (m *ClientMysql) GetMySqlClient() (*gorm.DB, error) {
	// 判断必填参数
	if m.config.DB == "" || m.config.Addr == "" || m.config.User == "" || m.config.Passwd == "" {
		return nil, fmt.Errorf("gorm: %s", "db, addr, user, passwd is required")
	}
	if m.config.MaxOpenConns == 0 {
		m.config.MaxOpenConns = 10
	}
	if m.config.MaxIdleConns == 0 {
		m.config.MaxIdleConns = 10
	}

	dsn := fmt.Sprintf("%s:%s@"+"tcp(%s)/%s?charset=%s",
		m.config.User, m.config.Passwd, m.config.Addr, m.config.DB, "utf8mb4")
	if m.config.TimeoutSec > 0 {
		// timeout in seconds has "s"
		dsn += fmt.Sprintf("&timeout=%ds", m.config.TimeoutSec)
	}
	dsn += "&parseTime=true"
	if m.config.Options != "" {
		dsn += "&" + strings.Trim(m.config.Options, "&")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:      m.log,
		QueryFields: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(m.config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(m.config.MaxIdleConns)
	return db, nil
}

func (m *ClientMysql) GetLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)
}
