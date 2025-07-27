package mysql

import (
	"dario.cat/mergo"
	"fmt"
	"github.com/zusux/gokit/mysql/mlog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strings"
)

type ConfigMysql struct {
	Enable  bool             `yaml:"enable" json:"enable"`
	Mapping map[string]Mysql `yaml:"mapping" json:"mapping"`
}

type Mysql struct {
	DB          string `yaml:"db" json:"db"`
	Addr        string `yaml:"addr" json:"addr"`
	User        string `yaml:"user" json:"user"`
	Passwd      string `yaml:"passwd" json:"passwd"`
	TimeoutSec  int    `yaml:"timeout_sec" json:"timeout_sec"`
	MaxOpenConn int    `yaml:"max_open_conn" json:"max_open_conn"`
	MaxIdleConn int    `yaml:"max_idle_conn" json:"max_idle_conn"`
	Metric      bool   `yaml:"metric" json:"metric"`
	Trace       bool   `yaml:"trace" json:"trace"`
	Options     string `yaml:"options" json:"options"`
}

func (m *Mysql) InitMysql() (*gorm.DB, error) {
	defaultConfig := Mysql{
		MaxOpenConn: 10,
		MaxIdleConn: 10,
	}
	err := mergo.Merge(&m, defaultConfig)
	if err != nil {
		return nil, err
	}
	// 判断必填参数
	if m.DB == "" || m.Addr == "" || m.User == "" || m.Passwd == "" {
		return nil, fmt.Errorf("gorm: %s", "db, addr, user, passwd is required")
	}

	dsn := fmt.Sprintf("%s:%s@"+"tcp(%s)/%s?charset=%s",
		m.User, m.Passwd, m.Addr, m.DB, "utf8mb4")
	if m.TimeoutSec > 0 {
		// timeout in seconds has "s"
		dsn += fmt.Sprintf("&timeout=%ds", m.TimeoutSec)
	}
	dsn += "&parseTime=true"
	if m.Options != "" {
		dsn += "&" + strings.Trim(m.Options, "&")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:      mlog.GetLogger(),
		QueryFields: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, err
	}
	_, err = db.DB()
	if err != nil {
		return nil, err
	}
	return db, nil
}
