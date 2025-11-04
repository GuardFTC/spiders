// Package client @Author:冯铁城 [17615007230@163.com] 2025-10-11 10:53:06
package client

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoClient mongo客户端
type MongoClient struct {
	client *mongo.Client
	ctx    context.Context
}

// NewMongoClient 创建mongo客户端
func NewMongoClient(config *MongoConfig, ctx context.Context) (*MongoClient, error) {

	//1.创建连接配置项
	clientOptions := options.Client().ApplyURI(config.Uri).
		SetMaxPoolSize(config.MaxPoolSize).
		SetMinPoolSize(config.MinPoolSize).
		SetConnectTimeout(config.ConnectTimeout).
		SetSocketTimeout(config.SocketTimeout)

	//2.建立连接,获取客户端
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	//3.创建客户端对象，返回
	return &MongoClient{
		client: client,
		ctx:    ctx,
	}, nil
}

// Ping 测试连接
func (m *MongoClient) Ping() error {
	return m.client.Ping(m.ctx, nil)
}

// Close 关闭连接
func (m *MongoClient) Close() error {
	return m.client.Disconnect(m.ctx)
}

// GetClient 获取客户端
func (m *MongoClient) GetClient() *mongo.Client {
	return m.client
}

// GetCtx 获取上下文
func (m *MongoClient) GetCtx() context.Context {
	return m.ctx
}
