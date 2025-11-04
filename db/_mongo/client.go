// Package _mongo @Author:冯铁城 [17615007230@163.com] 2025-10-11 10:53:06
package _mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongoClient mongo客户端
type mongoClient struct {
	client *mongo.Client
	ctx    context.Context
}

// newMongoClient 创建mongo客户端
func newMongoClient(config *mongoConfig, ctx context.Context) (*mongoClient, error) {

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
	return &mongoClient{
		client: client,
		ctx:    ctx,
	}, nil
}

// ping 测试连接
func (m *mongoClient) ping() error {
	return m.client.Ping(m.ctx, nil)
}

// close 关闭连接
func (m *mongoClient) close() error {
	return m.client.Disconnect(m.ctx)
}

// getClient 获取客户端
func (m *mongoClient) getClient() *mongo.Client {
	return m.client
}

// getCtx 获取上下文
func (m *mongoClient) getCtx() context.Context {
	return m.ctx
}
