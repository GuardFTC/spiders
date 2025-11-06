// Package spider @Author:冯铁城 [17615007230@163.com] 2025-11-03 14:19:45
package spider

import (
	"log"
	"spiders/spider/douban"
	"sync"
)

type Spider interface {
	GetName() string // 获取爬虫名称
	CanRun() bool    // 判断是否可以运行
	Run()            // 运行
}

// Init 初始化所有爬虫
func Init() []Spider {

	//1.创建爬虫切片
	var spiders []Spider

	//2.创建豆瓣top250电影爬虫,写入切片
	top250MovieSpider := douban.NewTop250MovieSpider()
	spiders = append(spiders, top250MovieSpider)

	//3.创建豆瓣top250图书爬虫,写入切片
	top250BookSpider := douban.NewTop250BookSpider()
	spiders = append(spiders, top250BookSpider)

	//4.返回切片集合
	return spiders
}

// RunSpiders 运行所有爬虫
func RunSpiders(spiders []Spider) {

	//1.定义wait group
	var wg sync.WaitGroup

	//2.遍历爬虫集合，异步运行爬取数据
	for _, _spider := range spiders {

		//3.主协程wg加1
		wg.Add(1)

		//4.创建协程异步运行爬虫
		go func(spider Spider) {

			//5.确保最终释放锁
			defer wg.Done()

			//6.打印起止日志，运行爬虫
			if spider.CanRun() {
				log.Printf("[ %s ] 开始运行", spider.GetName())
				spider.Run()
				log.Printf("[ %s ] 运行完毕", spider.GetName())
			} else {
				log.Printf("[ %s ] 运行条件不满足", spider.GetName())
			}
		}(_spider)
	}

	//7.阻塞等待所有爬虫爬取完毕
	wg.Wait()
}
