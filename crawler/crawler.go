package crawler

import (
	"fmt"
	"github.com/ironbang/proxypool/crawler/spider"
	"time"
)

func Crawler() {
	fmt.Println("启动爬虫模块...")
	for {
		func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
				}
			}()
			// www.89ip.cn
			spider.IP89Spider()
		}()
		time.Sleep(time.Duration(20) * time.Minute)
	}
}
