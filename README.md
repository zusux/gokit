# gokit  日志, 数据库, protobuf 处理工具

## 日志处理
> 配置文件
```yaml
log:
  app: "admin"
  write_console: true
  write_file: true
  path: "logs"
  file: "out.log"
  err_file: "err.log"
  age: 7
  rotation: 24
  timestamp_format: "2006-01-02T15:04:05.000Z"
  logger_level: 0
  logger_format: "console"
```
> go config
```go
package conf
import (
	"github.com/zusux/gokit/zlog"
)

type Config struct {
	Log      zlog.Logger `yaml:"log" json:"log"`
}

```
> usage
```go
package main

import (
	"flag"
	"conf"
	conf2 "github.com/zusux/gokit/conf"
)
var (
	// flagconf is the config flag.
	flagconf string
)
func init() {
	flag.StringVar(&flagconf, "conf", "configs/config.yaml", "config path, eg: -conf config.yaml")
}
//支持函数直接调用与对象调用 底层使用zap log
var bc conf.Config
conf2.MustLoad(flagconf, &bc)
log := conf.Config.Log
loger := zlog.NewSlog( config.InitLog(),2) //2表示调用的级别 可以自由调节
loger.Infof("conf:%v",bc)

zlog.Infof("conf:%v",bc)

```
## 数据库处理
```yaml
database:
  mysql:
    enable: true
    mysql_config:
      db: "dbname"
      addr: "127.0.0.1"
      user: "root"
      passwd: "123456"
      timeout_sec: 3000
      metric: false
      trace: false
      options: ""
  redis:
    host: 127.0.0.1:6379
    pass: ""
    db: 0
  mongo:
    enable: true
    uri: "mongodb://127.0.0.1:27017/manage?retryWrites=true&w=majority"
```

>使用
```go
package conf

import (
	"github.com/zusux/gokit/mon"
	"github.com/zusux/gokit/mysql"
	"github.com/zusux/gokit/redis"
)
type Database struct {
	Mysql mysql.ConfigMysql `yaml:"mysql" json:"mysql"`
	Redis redis.ConfigRedis `yaml:"redis" json:"redis"`
	Mongo mon.ConfigMongo   `yaml:"mongo" json:"mongo"`
}

```
```go
package data

import (
	"conf"
	goRedis "github.com/redis/go-redis/v9"
	"github.com/zusux/gokit/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

func Init(conf *conf.Database){
redisClient := redis.NewRedisClient(&conf.Redis).GetClient()
mdb := conf.Mongo.Init(context.Background())
}

```


## proto tag
* install tools
```
go install github.com/zusux/gokit/protoc-gen-tag/protoc-gen-tag.go
```
* example proto file

user.proto
```proto
 message User {
  string id = 1 [(tag.json) = "id", (tag.bson) = "_id",(tag.gorm) = "column:id"];
  string name = 2 [(tag.json) = "name", (tag.gorm) = "column:name",(tag.yaml) = "name"];
  string age = 3 ;
}
```
* gen file

> protoc --proto_path=. --proto_path=./third_party  --go_out=.  --tag_out=. protoc-gen-tag\example\user.proto 

example user.tag.go
```go
package example
type UserTag struct {
    Id   string `json:"id" bson:"_id" gorm:"column:id"`
    Name string `json:"name" gorm:"column:name"`
    Age  string `json:"age,omitempty"`
}

func (x *User) ToUserTag() *UserTag {
    if x == nil {
        return nil
    }
return &UserTag{
    Id:   x.Id,
    Name: x.Name,
    Age:  x.Age,
    }
}
func (x *UserTag) ToUser() *User {
    if x == nil {
        return nil
    }
    return &User{
        Id:   x.Id,
        Name: x.Name,
        Age:  x.Age,
    }
}
```


