// Package main @Author:冯铁城 [17615007230@163.com] 2025-11-03 11:31:00
package main

import (
	"log"
	"spiders/spider"
	"sync"
)

func main() {

	//1.初始化所有爬虫
	spiders := spider.Init()

	//2.定义wait group
	var wg sync.WaitGroup

	//3.遍历爬虫集合，异步运行爬取数据
	for _, _spider := range spiders {

		//4.主协程wg加1
		wg.Add(1)

		//5.创建协程异步运行爬虫
		go func() {

			//6.确保最终释放锁
			defer wg.Done()

			//7.打印起止日志，运行爬虫
			log.Printf("爬虫:%s 开始运行", _spider.GetName())
			_spider.Run()
			log.Printf("爬虫:%s 运行完毕", _spider.GetName())
		}()
	}

	//8.阻塞等待所有爬虫爬取完毕
	wg.Wait()
}
