package mysql

import (
	"github.com/zusux/gokit/zlog"
	"gorm.io/gorm"
)

var mysqlClientMapping = map[string]*gorm.DB{} // 用于存储 mysql 连接 初始化后只读

func (c *ConfigMysql) InitMysql() {
	for k, v := range c.Mapping {
		db, err := v.InitMysql()
		if err != nil {
			zlog.Fatalf("mysql init failed: %v", err)
		}
		mysqlClientMapping[k] = db
	}
}

func GetMysqlClient(name string) *gorm.DB {
	return mysqlClientMapping[name]
}
