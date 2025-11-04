// Package _mongo @Author:冯铁城 [17615007230@163.com] 2025-10-11 17:38:54
package _mongo

import (
	"context"
	"log"
	"time"
)

// client mongo客户端
var client *mongoClient

// CreateMongoClient 创建客户端
func CreateMongoClient() {

	//1.创建MongoDB配置
	_mongoConfig := &mongoConfig{
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
	_mongoClient, err := newMongoClient(_mongoConfig, ctx)
	if err != nil {
		log.Fatalf("mongo create error: %v", err)
	}

	//4.测试链接
	err = _mongoClient.ping()
	if err != nil {
		log.Fatalf("mongo test connect error: %v", err)
	}

	//5.客户端赋值
	client = _mongoClient

	//6.打印日志
	log.Println("mongo client create success")
}

// CloseMongoClient 关闭客户端
func CloseMongoClient() {
	if err := client.close(); err != nil {
		log.Fatalf("mongo close error: %v", err)
	} else {
		log.Println("mongo client close success")
	}
}
