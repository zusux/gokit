package mongo

import (
	"context"
	"github.com/zusux/gokit/zlog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClientMapping = map[string]*mongo.Client{} // 用于存储 mongo 连接 初始化后只读

func (m *ConfigMongo) Init(ctx context.Context) {
	// 替换成你的 MongoDB 地址
	if m.Enable {
		for k, v := range m.Mapping {
			clientOptions := options.Client().ApplyURI(v)
			// 连接
			client, err := mongo.Connect(ctx, clientOptions)
			if err != nil {
				zlog.Fatalf("连接失败: err:%v", err)
			}

			// Ping 测试
			err = client.Ping(ctx, nil)
			if err != nil {
				zlog.Fatalf("Ping失败: err:%v", err)
			}
			zlog.Info("MongoDB连接成功！")
			mongoClientMapping[k] = client
		}
	}
}

func GetMongoClient(name string) *mongo.Client {
	return mongoClientMapping[name]
}
