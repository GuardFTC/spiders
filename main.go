// Package main @Author:冯铁城 [17615007230@163.com] 2025-11-03 11:31:00
package main

import (
	"spiders/db"
	"spiders/spider"
)

func main() {

	//1.初始化Mongo客户端
	db.InitDatabase()
	defer db.CloseDatabase()

	//2.初始化所有爬虫
	spiders := spider.Init()

	//3.运行所有爬虫
	spider.RunSpiders(spiders)
}
