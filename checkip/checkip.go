package checkip

import (
	"ProxyPool/database"
	"fmt"
	"github.com/ironbang/httpclient"
	"net/http"
	"time"
)

func checkProxyIp(proxy string) (*database.IPInfo, error) {
	// 写入数据库
	ipinfo := database.NewIPInfo(proxy)

	client := &httpclient.HttpClient{ProxyScheme: "http", ProxyIp: proxy, DialTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second}
	client, err := client.NewClient()
	t := time.Now() // get current time
	resp, err := client.Get("http://httpbin.org/ip")
	ipinfo.LastCheckTime = t.Format("2006-01-02 15:04:05")
	elapsed := time.Since(t)
	ipinfo.Speed = elapsed.Seconds()
	if err != nil {
		ipinfo.Result = append(ipinfo.Result, false)
		fmt.Printf("[时长: %f] IP[%s]失效[%s]\n", elapsed.Seconds(), proxy, err.Error())
	} else {
		if resp != nil {
			if resp.StatusCode == http.StatusOK {
				fmt.Printf("[时长: %f] 检测IP[%s]成功\n", elapsed.Seconds(), proxy)
				ipinfo.Result = append(ipinfo.Result, true)
			} else {
				ipinfo.Result = append(ipinfo.Result, false)
				fmt.Printf("[时长: %f] IP[%s]失效[%d]\n", elapsed.Seconds(), proxy, resp.StatusCode)
			}
		} else {
			ipinfo.Result = append(ipinfo.Result, false)
			fmt.Printf("[时长: %f] IP[%s]失效[%d]\n", elapsed.Seconds(), proxy, resp.StatusCode)
		}
	}
	return ipinfo, nil
}

func CheckIp(sysChan <-chan string) {
	fmt.Println("校验模块...")
	store := database.NewStore()
	for {
		func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
				}
			}()
			ip := <-sysChan

			go func(proxy string) {
				ipinfo, _ := checkProxyIp(proxy)
				store.Put(ipinfo)
			}(ip)
		}()
	}
}
