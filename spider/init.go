// Package spider @Author:冯铁城 [17615007230@163.com] 2025-11-03 14:19:45
package spider

import (
	"spiders/spider/douban"
)

type Spider interface {
	GetName() string
	Run()
}

// Init 初始化所有爬虫
func Init() []Spider {

	//1.创建爬虫切片
	var spiders []Spider

	//2.创建豆瓣top250电影爬虫,写入切片
	top250MovieSpider := douban.NewTop250MovieSpider()
	spiders = append(spiders, top250MovieSpider)

	//3.返回切片集合
	return spiders
}
