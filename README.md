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


