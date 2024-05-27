# gokit

```yaml
ConfigKafka:
	Brokers: [""]
	Async: true
	Timeout:  500
	Name:    string
	Group:   string
	Topic:   string

ConfigRedis:
  Host: ""
  Pass: ""
  Db: 0

ConfigMysql:
  DB: ""
  Addr: "127.0.0.1:3306"
  User:  "root"
  Passwd: "123456"
  TimeoutSec:   500
  MaxOpenConns: 20
  MaxIdleConns: 20
  Metric:       false
  Trace:        false
  Options: ""

```