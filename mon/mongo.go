package mon

import (
	"context"
	"time"

	"github.com/zusux/gokit/zlog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *ConfigMongo) Init(ctx context.Context) (*mongo.Client, error) {
	// 替换成你的 MongoDB 地址
	if m.Enable {
		clientOptions := options.Client().ApplyURI(m.Uri).
			SetConnectTimeout(5 * time.Second).
			SetServerSelectionTimeout(5 * time.Second)
		// 连接
		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			zlog.Errorf("MongoDB连接失败 err:%v", err.Error())
			return nil, err
		}
		if client == nil {
			zlog.Warn("MongoDB不可用，将跳过数据库操作")
			return nil, err
		}
		// Ping 测试
		err = client.Ping(ctx, nil)
		if err != nil {
			zlog.Errorf("MongoDB Ping失败 err:%v", err.Error())
			return nil, err
		}
		zlog.Info("MongoDB连接成功！")
		return client, err
	}
	zlog.Warn("MongoDB 连接开关 关闭!")
	return nil, nil
}

func (m *ConfigMongo) MustInit(ctx context.Context) *mongo.Client {
	clientOptions := options.Client().ApplyURI(m.Uri).
		SetConnectTimeout(5 * time.Second).
		SetServerSelectionTimeout(5 * time.Second)
	// 连接
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		zlog.Errorf("MongoDB连接失败 err:%v", err.Error())
		return nil
	}
	if client == nil {
		zlog.Warn("MongoDB不可用，将跳过数据库操作")
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
