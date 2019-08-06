package main

import (
	"ProxyPool/checkip"
	"ProxyPool/crawler"
	"ProxyPool/database"
	"ProxyPool/restful"
	"fmt"
	"net/http"
	"runtime/pprof"
	"sync"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	p := pprof.Lookup("goroutine")
	p.WriteTo(w, 1)
}

func main() {

	go func() {
		http.HandleFunc("/", handler)
		http.ListenAndServe(":11181", nil)
	}()

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
