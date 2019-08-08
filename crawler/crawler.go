package crawler

import (
	"fmt"
	"github.com/ironbang/proxypool/crawler/spider"
	"sync"
	"time"
)

func Crawler(sysChan chan<- string, group *sync.WaitGroup) {
	fmt.Println("启动爬虫模块...")
	for {
		func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
				}
			}()
			// www.89ip.cn
			spider.IP89Spider(sysChan)
		}()
		time.Sleep(time.Duration(20) * time.Minute)
	}
}
