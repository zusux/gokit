package mysql

import (
	"github.com/zusux/gokit/zlog"
	"gorm.io/gorm"
)

func (c *ConfigMysql) InitMysql() *gorm.DB {
	db, err := c.MysqlConfig.InitMysql()
	if err != nil {
		zlog.Fatalf("mysql init failed: %v", err)
	}
	return db
}
