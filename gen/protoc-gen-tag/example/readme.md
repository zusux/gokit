go install github.com/zusux/gokit/gen/protoc-gen-tag/protoc-gen-tag.go

protoc --proto_path=. --proto_path=./third_party  --go_out=.  --tag_out=. gen\protoc-gen-tag\example\user.proto 