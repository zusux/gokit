package mon

import (
	"context"
	"github.com/zusux/gokit/zlog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *ConfigMongo) Init(ctx context.Context) *mongo.Client {
	// 替换成你的 MongoDB 地址
	if m.Enable {
		clientOptions := options.Client().ApplyURI(m.Uri)
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
		return client
	}
	return nil
}
