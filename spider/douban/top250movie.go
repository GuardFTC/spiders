// Package douban @Author:冯铁城 [17615007230@163.com] 2025-11-03 11:31:40
package douban

import (
	"fmt"
	"log"
	"spiders/common/util"
	"spiders/db/_mongo"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

//-------------------------------------------爬虫结构体相关--------------------------------------------------

// NewTop250MovieSpider 创建豆瓣Top250电影爬虫
func NewTop250MovieSpider() *Top250MovieSpider {

	//1.创建电影列表
	movies := make([]*Movie, 0)

	//2.获取访问链接
	urls := getTop250MovieUrls()

	//3.创建结构体实例
	top250MovieSpider := &Top250MovieSpider{
		urls:   urls,
		movies: movies,
	}

	//4.获取并设置详情页采集器
	infoCollection := getInfoCollection(top250MovieSpider)
	top250MovieSpider.infoCollection = infoCollection

	//5.获取并设置列表页采集器
	listCollection := getListCollection(infoCollection)
	top250MovieSpider.listCollection = listCollection

	//6.返回
	return top250MovieSpider
}

// Top250MovieSpider 爬取豆瓣电影top250
type Top250MovieSpider struct {
	urls           []string         //访问链接切片
	listCollection *colly.Collector //列表页采集器
	infoCollection *colly.Collector //详情页采集器
	movies         []*Movie         //电影列表
	mu             sync.Mutex       //锁
}

// GetName 获取爬虫名称
func (t *Top250MovieSpider) GetName() string {
	return "豆瓣top250电影爬虫"
}

// Run 运=运行爬虫，爬取数据
func (t *Top250MovieSpider) Run() {

	//1.遍历访问链接，爬取数据
	for _, url := range t.urls {
		if err := t.listCollection.Visit(url); err != nil {
			log.Printf("访问地址：%v 异常：%v", url, err)
		}
	}

	//2.保存数据到Mongo
	if err := _mongo.DeleteAndSaveData(t.movies, "top250_movies", "douban"); err != nil {
		log.Fatalf("保存数据异常：%v", err)
	}
}

//-------------------------------------------爬虫方法相关-----------------------------------------------------

// GetUrls 获取访问链接
func getTop250MovieUrls() []string {

	//1.定义访问链接切片
	var urls []string

	//2.循环生成访问链接，添加到切片中
	for i := 0; i < 250; i += 25 {
		urls = append(urls, fmt.Sprintf(top250MovieURL, strconv.Itoa(i)))
	}

	//3.返回访问链接切片
	return urls
}

// GetListCollection 获取列表页采集器
func getListCollection(infoCollection *colly.Collector) *colly.Collector {

	//1.创建采集器
	listCollection := colly.NewCollector(
		colly.UserAgent(defaultUserAgent),
	)

	//2.设置limitRule
	if err := listCollection.Limit(&colly.LimitRule{
		DomainGlob:  "movie.douban.com",
		Parallelism: 1,
		Delay:       3 * time.Second,
		RandomDelay: 500 * time.Millisecond,
	}); err != nil {
		log.Fatalf("设置limitRule异常：%v", err)
	}

	//3.设置请求之前的回调
	listCollection.OnRequest(func(request *colly.Request) {
		log.Printf("开始访问列表页地址：%v", request.URL)
	})

	//4.定义异常回调
	listCollection.OnError(func(response *colly.Response, err error) {
		log.Printf("访问地址：%v 异常：%v", response.Request.URL, err)
	})

	//5.定义响应回调
	listCollection.OnResponse(func(response *colly.Response) {
		if response.StatusCode < 200 || response.StatusCode >= 400 {
			log.Printf("访问地址：%v 状态码异常：%v", response.Request.URL, response.StatusCode)
		}
	})

	//6.定义HTML回调
	listCollection.OnHTML("ol[class='grid_view']", func(element *colly.HTMLElement) {

		//7.遍历li标签
		element.ForEach("li", func(index int, eli *colly.HTMLElement) {

			//8.获取电影详情链接
			movieURL := eli.ChildAttr("div[class='pic'] > a", "href")

			//9.访问详情页
			if err := infoCollection.Visit(movieURL); err != nil {
				log.Printf("访问地址：%v 异常：%v", movieURL, err)
			}
		})
	})

	//10.返回采集器
	return listCollection
}

// GetInfoCollection 获取详情页采集器
func getInfoCollection(top250MovieSpider *Top250MovieSpider) *colly.Collector {

	//1.创建采集器
	infoCollection := colly.NewCollector(
		colly.UserAgent(defaultUserAgent),
	)

	//2.设置limitRule
	if err := infoCollection.Limit(&colly.LimitRule{
		DomainGlob:  "movie.douban.com",
		Parallelism: 1,
		Delay:       4 * time.Second,
		RandomDelay: 700 * time.Millisecond,
	}); err != nil {
		log.Fatalf("设置limitRule异常：%v", err)
	}

	//3.设置请求之前的回调
	infoCollection.OnRequest(func(request *colly.Request) {
		log.Printf("开始访问详情页地址：%v", request.URL)
	})

	//4.定义异常回调
	infoCollection.OnError(func(response *colly.Response, err error) {
		log.Printf("访问地址：%v 异常：%v", response.Request.URL, err)
	})

	//5.定义响应回调
	infoCollection.OnResponse(func(response *colly.Response) {
		if response.StatusCode < 200 || response.StatusCode >= 400 {
			log.Printf("访问地址：%v 状态码异常：%v", response.Request.URL, response.StatusCode)
		}
	})

	//6.定义HTML回调
	infoCollection.OnHTML("#content", func(element *colly.HTMLElement) {

		//7.解析电影信息
		movie := parseMovie(element)

		//8.加锁
		top250MovieSpider.mu.Lock()
		defer top250MovieSpider.mu.Unlock()

		//9.存入集合
		top250MovieSpider.movies = append(top250MovieSpider.movies, movie)
	})

	//10.返回采集器
	return infoCollection
}

//-------------------------------------------被爬数据结构体相关--------------------------------------------------

// Movie 电影信息实体类
type Movie struct {
	Rank         string            `json:"rank"`         // 排名，如 "No.1"
	Name         string            `json:"name"`         // 电影名称
	Year         string            `json:"year"`         // 年份，如 "(1994)"
	CoverURL     string            `json:"coverUrl"`     // 封面图链接
	Directors    map[string]string `json:"directors"`    // 导演: 姓名 -> URL
	Writers      map[string]string `json:"writers"`      // 编剧: 姓名 -> URL
	Actors       map[string]string `json:"actors"`       // 主演: 姓名 -> URL
	Genres       []string          `json:"genres"`       // 类型列表
	Country      string            `json:"country"`      // 制片国家/地区
	Language     []string          `json:"language"`     // 语言
	ReleaseDates []string          `json:"releaseDates"` // 上映日期（多地）
	Runtime      string            `json:"runtime"`      // 片长
	Alias        string            `json:"alias"`        // 又名
	IMDb         string            `json:"imdb"`         // IMDb 编号
	Score        string            `json:"score"`        // 评分
	ScoreCount   string            `json:"scoreCount"`   // 评分人数（含“人评价”或纯数字）
	Summary      string            `json:"summary"`      // 剧情简介（已清理）
}

// NewMovie 初始化 Movie 实例
func NewMovie() *Movie {
	return &Movie{
		Directors:    make(map[string]string),
		Writers:      make(map[string]string),
		Actors:       make(map[string]string),
		Genres:       make([]string, 0),
		Language:     make([]string, 0),
		ReleaseDates: make([]string, 0),
	}
}

// parseMovie 解析电影信息
func parseMovie(element *colly.HTMLElement) *Movie {

	//1.新建电影结构体
	movie := NewMovie()

	//2.解析电影排名
	movie.Rank = element.ChildText("span[class='top250-no']")

	//3.解析电影名称
	movie.Name = element.ChildText("h1:first-of-type > span:first-of-type")

	//4.解析电影年份
	movie.Year = element.ChildText("h1:first-of-type > span[class='year']")

	//5.解析电影封面图片链接
	movie.CoverURL = element.ChildAttr("#mainpic > a > img", "src")

	//5.解析导演信息
	movieDirectors := element.ChildText("span:contains(导演) + span[class='attrs']")
	movieDirectorURLs := element.ChildAttrs("span:contains(导演) + span[class='attrs'] > a", "href")
	for i, director := range strings.Split(movieDirectors, " / ") {
		if i < len(movieDirectorURLs) {
			movie.Directors[strings.TrimSpace(director)] = movieDirectorURLs[i]
		}
	}

	//6.解析编剧信息
	movieWriters := element.ChildText("span:contains(编剧) + span[class='attrs']")
	movieWriterURLs := element.ChildAttrs("span:contains(编剧) + span[class='attrs'] > a", "href")
	for i, writer := range strings.Split(movieWriters, " / ") {
		if i < len(movieWriterURLs) {
			movie.Writers[strings.TrimSpace(writer)] = movieWriterURLs[i]
		}
	}

	//7.解析主演信息
	movieActors := element.ChildText("span:contains(主演) + span[class='attrs']")
	movieActorURLs := element.ChildAttrs("span:contains(主演) + span[class='attrs'] > a", "href")
	for i, actor := range strings.Split(movieActors, " / ") {
		if i < len(movieActorURLs) {
			movie.Actors[strings.TrimSpace(actor)] = movieActorURLs[i]
		}
	}

	//8.解析类型信息
	element.ForEach("span:contains('类型:')", func(i int, e *colly.HTMLElement) {
		var types []string
		e.DOM.NextUntil("br").Each(func(i int, s *goquery.Selection) {
			if s.Is("span[property='v:genre']") {
				types = append(types, s.Text())
			}
		})
		movie.Genres = types
	})

	//9.解析制片国家/地区
	element.ForEach("span:contains('制片国家/地区:')", func(i int, e *colly.HTMLElement) {
		movie.Country = strings.TrimSpace(e.DOM.Get(0).NextSibling.Data)
	})

	//10.解析语言
	element.ForEach("span:contains('语言:')", func(i int, e *colly.HTMLElement) {
		movie.Language = strings.Split(strings.TrimSpace(e.DOM.Get(0).NextSibling.Data), " / ")
	})

	//11.解析上映日期
	element.ForEach("span:contains('上映日期:')", func(i int, e *colly.HTMLElement) {
		var dates []string
		e.DOM.NextUntil("br").Each(func(i int, s *goquery.Selection) {
			if s.Is("span[property='v:initialReleaseDate']") {
				dates = append(dates, s.Text())
			}
		})
		movie.ReleaseDates = dates
	})

	//12.解析片长
	movie.Runtime = element.ChildText("span[property='v:runtime']")

	//13.解析又名
	element.ForEach("span:contains('又名:')", func(i int, e *colly.HTMLElement) {
		movie.Alias = strings.TrimSpace(e.DOM.Get(0).NextSibling.Data)
	})

	//14.解析IMDb
	element.ForEach("span:contains('IMDb:')", func(i int, e *colly.HTMLElement) {
		movie.IMDb = strings.TrimSpace(e.DOM.Get(0).NextSibling.Data)
	})

	//15.解析评分
	movie.Score = element.ChildText("strong[class='ll rating_num']")

	//16.解析评分人数
	movie.ScoreCount = element.ChildText("div[class='rating_sum']")

	//17.解析剧情简介
	movieSummary := element.ChildText("span[class='all hidden']")
	if movieSummary == "" {
		movieSummary = element.ChildText("span[property='v:summary']")
	}
	movie.Summary = util.RemoveSpacesAndNewlines(movieSummary)

	//18.返回movie
	return movie
}
