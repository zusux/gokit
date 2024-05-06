package mysql

type ConfigMysql struct {
	DB           string
	Addr         string
	User         string
	Passwd       string
	TimeoutSec   int
	MaxOpenConns int
	MaxIdleConns int
	Metric       bool
	Trace        bool
	Options      string
}
