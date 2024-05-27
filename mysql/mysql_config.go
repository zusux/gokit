package mysql

type ConfigMysql struct {
	DB           string `yaml:"DB"`
	Addr         string `yaml:"Addr"`
	User         string `yaml:"User"`
	Passwd       string `yaml:"Passwd"`
	TimeoutSec   int    `yaml:"TimeoutSec"`
	MaxOpenConns int    `yaml:"MaxOpenConns"`
	MaxIdleConns int    `yaml:"MaxIdleConns"`
	Metric       bool   `yaml:"Metric"`
	Trace        bool   `yaml:"Trace"`
	Options      string `yaml:"Options"`
}
