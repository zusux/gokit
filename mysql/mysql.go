package mysql

import (
	"log"

	"gorm.io/gorm"
)

func (c *ConfigMysql) MustInitMysql() *gorm.DB {
	db, err := c.MysqlConfig.InitMysql()
	if err != nil {
		log.Fatalf("mysql init failed: %v", err)
	}
	return db
}

func (c *ConfigMysql) InitMysql() (*gorm.DB, error) {
	return c.MysqlConfig.InitMysql()
}
