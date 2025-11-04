// Package main @Author:冯铁城 [17615007230@163.com] 2025-11-03 11:31:00
package main

import "spiders/spider"

func main() {

	//1.初始化所有爬虫
	spiders := spider.Init()

	//2.遍历爬虫集合，运行爬取数据
	for _, _spider := range spiders {
		_spider.Run()
	}
}
