// Package client @Author:冯铁城 [17615007230@163.com] 2025-10-11 17:38:54
package client

import (
	"context"
	"log"
	"time"
)

// CreateMongoClient 创建客户端
func CreateMongoClient() *MongoClient {

	//1.创建MongoDB配置
	mongoConfig := &MongoConfig{
		//Uri: "mongodb://myuser:mypassword@localhost:27017/testdb?replicaSet=rs0", //带用户名和密码，以及副本集
		Uri:            "mongodb://127.0.0.1:27017",
		MaxPoolSize:    50,
		MinPoolSize:    5,
		ConnectTimeout: 10 * time.Second,
		SocketTimeout:  10 * time.Second,
	}

	//2.创建上下文
	ctx := context.Background()

	//3.创建客户端
	mongoClient, err := NewMongoClient(mongoConfig, ctx)
	if err != nil {
		log.Fatalf("mongo create error: %v", err)
	}

	//4.测试链接
	err = mongoClient.Ping()
	if err != nil {
		log.Fatalf("mongo test connect error: %v", err)
	}

	//5.打印日志
	log.Println("mongo client create success")

	//6.返回客户端
	return mongoClient
}

// CloseMongoClient 关闭客户端
func CloseMongoClient(mongoClient *MongoClient) {
	if err := mongoClient.Close(); err != nil {
		log.Fatalf("mongo close error: %v", err)
	} else {
		log.Println("mongo client close success")
	}
}
