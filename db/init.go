// Package db @Author:冯铁城 [17615007230@163.com] 2025-11-04 16:51:22
package db

import "spiders/db/_mongo"

// InitDatabase 初始化数据库
func InitDatabase() {

	//1.初始化Mongo客户端
	_mongo.CreateMongoClient()
}

// CloseDatabase 关闭数据库
func CloseDatabase() {

	//1.关闭Mongo客户端
	_mongo.CloseMongoClient()
}
