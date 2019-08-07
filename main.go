package main

import (
	"fmt"
	"github.com/ironbang/proxypool/checkip"
	"github.com/ironbang/proxypool/crawler"
	"github.com/ironbang/proxypool/database"
	"github.com/ironbang/proxypool/restful"
	"sync"
)

func main() {
	func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		sysChan := make(chan string, 2000)

		wait := sync.WaitGroup{}
		wait.Add(4)

		go database.CheckStore(sysChan)

		// RESTFul模块
		go restful.RESTFul()

		// 爬虫模块
		go crawler.Crawler(sysChan)

		// 校验模块
		go checkip.CheckIp(sysChan)

		wait.Wait()
	}()
}
