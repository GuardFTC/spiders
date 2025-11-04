// Package douban @Author:冯铁城 [17615007230@163.com] 2025-11-04 19:04:43
package douban

import (
	"fmt"
	"log"
	"regexp"
	"spiders/db/_mongo"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

//-------------------------------------------爬虫结构体相关--------------------------------------------------

// NewTop250BookSpider 创建豆瓣Top250图书爬虫
func NewTop250BookSpider() *Top250MovieSpider {

	//1.创建电影列表
	movies := make([]*Movie, 0)

	//2.获取访问链接
	urls := getTop250MovieUrls()

	//3.创建结构体实例
	top250MovieSpider := &Top250MovieSpider{
		urls:   urls,
		movies: movies,
		canRun: true,
	}

	//4.获取并设置详情页采集器
	infoCollection := getTop250MovieInfoCollection(top250MovieSpider)
	top250MovieSpider.infoCollection = infoCollection

	//5.获取并设置列表页采集器
	listCollection := getTop250MovieListCollection(infoCollection, top250MovieSpider)
	top250MovieSpider.listCollection = listCollection

	//6.返回
	return top250MovieSpider
}

type Top250BookSpider struct {
	urls           []string         //访问链接切片
	listCollection *colly.Collector //列表页采集器
	books          []*Book          //图书列表
	mu             sync.Mutex       //锁
	canRun         bool             //是否可以运行
}

// GetName 获取爬虫名称
func (t *Top250BookSpider) GetName() string {
	return "豆瓣top250图书爬虫"
}

// CanRun 获取是否可以运行
func (t *Top250BookSpider) CanRun() bool {
	return t.canRun
}

// Run 运=运行爬虫，爬取数据
func (t *Top250BookSpider) Run() {

	//1.遍历访问链接，爬取数据
	for _, url := range t.urls {
		if err := t.listCollection.Visit(url); err != nil {
			log.Printf("访问地址：%v 异常：%v", url, err)
		}
	}

	//2.保存数据到Mongo
	if err := _mongo.DeleteAndSaveData(t.books, defaultDbName, top250BookCollectionName); err != nil {
		log.Printf("保存数据异常：%v", err)
		return
	}
}

//-------------------------------------------爬虫方法相关-----------------------------------------------------

// 获取top250图书访问URL
func getTop250BookUrls() []string {

	//1.定义url切片
	var urls []string

	//2.循环总页数，获取所有URL，并加入url切片
	for i := 0; i < 10; i++ {

		//3.格式化url
		url := fmt.Sprintf(top250BookURL, i*25)

		//4.存入切片
		urls = append(urls, url)
	}

	//5.返回切片
	return urls
}

//-------------------------------------------被爬数据结构体相关--------------------------------------------------

// Book 图书信息结构体
type Book struct {
	CoverImg    string `json:"cover_img"`    // 图书封面
	Title       string `json:"title"`        // 图书标题
	Author      string `json:"author"`       // 作者
	Publisher   string `json:"publisher"`    // 出版社
	PublishTime string `json:"publish_time"` // 出版时间
	Price       string `json:"price"`        // 单价
	Rating      string `json:"rating"`       // 评分
	RatingCount string `json:"rating_count"` // 评价数
	Description string `json:"description"`  // 描述
	EbookLink   string `json:"ebook_link"`   // 电子版链接
}

// newBook 创建图书结构体实例
func parseBook(el *colly.HTMLElement) *Book {

	//1.创建图书结构体实例
	book := new(Book)

	//2.获取图书封面
	book.CoverImg = el.ChildAttr("td:nth-of-type(1) > a > img", "src")

	//3.获取图书标题
	book.Title = el.ChildText("td:nth-of-type(2) > div[class='pl2'] > a")

	//4.填充图书信息相关字段
	infoText := el.ChildText("td:nth-of-type(2) > p[class='pl']")
	infos := strings.Split(infoText, "/")
	if len(infos) == 4 {
		book.Author = strings.TrimSpace(infos[0])
		book.Publisher = strings.TrimSpace(infos[1])
		book.PublishTime = strings.TrimSpace(infos[2])
		book.Price = strings.TrimSpace(infos[3])
	} else if len(infos) > 4 {
		book.Author = strings.TrimSpace(infos[0])
		book.Publisher = strings.TrimSpace(infos[1]) + "/" + strings.TrimSpace(infos[2])
		book.PublishTime = strings.TrimSpace(infos[3])
		book.Price = strings.TrimSpace(infos[4])
	}

	//5.获取评价数以及评分
	book.Rating = el.ChildText("td:nth-of-type(2) > div[class='star clearfix'] > span[class='rating_nums']")
	RatingCountStr := el.ChildText("td:nth-of-type(2) > div[class='star clearfix'] > span[class='pl']")
	book.RatingCount = regexp.MustCompile(`\d+人`).FindString(strings.ReplaceAll(RatingCountStr, "\n", ""))

	//6.获取描述
	book.Description = el.ChildText("td:nth-of-type(2) > p[class='quote'] > span[class='inq']")

	//7.返回book
	return book
}
