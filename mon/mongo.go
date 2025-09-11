package mon

import (
	"context"
	"net/url"
	"time"

	"github.com/zusux/gokit/zlog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *ConfigMongo) Init(ctx context.Context) *mongo.Client {
	// 替换成你的 MongoDB 地址
	if m.Enable {
		clientOptions := options.Client().ApplyURI(url.QueryEscape(m.Uri)).
			SetConnectTimeout(5 * time.Second).
			SetServerSelectionTimeout(5 * time.Second)
		// 连接
		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			zlog.Errorf("MongoDB连接失败 err:%v", err.Error())
			return nil
		}

		// Ping 测试
		err = client.Ping(ctx, nil)
		if err != nil {
			zlog.Errorf("MongoDB Ping失败 err:%v", err.Error())
			return nil
		}
		zlog.Info("MongoDB连接成功！")
		return client
	}
	return nil
}
