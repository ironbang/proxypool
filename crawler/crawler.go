package crawler

import (
	"ProxyPool/crawler/spider"
	"fmt"
	"time"
)

func Crawler(sysChan chan<- string) {
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
		time.Sleep(time.Duration(10) * time.Minute)
	}
}
