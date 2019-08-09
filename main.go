package main

import (
	"fmt"
	"github.com/ironbang/proxypool/checkip"
	"github.com/ironbang/proxypool/crawler"
	"github.com/ironbang/proxypool/restful"
)

func main() {
	func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()

		// 校验模块
		go checkip.CheckIp()

		go checkip.CheckStore()

		// 爬虫模块
		go crawler.Crawler()

		// RESTFul模块
		restful.RESTFul()
	}()
}
