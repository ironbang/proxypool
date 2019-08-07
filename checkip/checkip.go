package checkip

import (
	"fmt"
	"github.com/ironbang/httpclient"
	"github.com/ironbang/proxypool/common/function"
	"github.com/ironbang/proxypool/database"
	"net/http"
	"time"
)

func checkProxyIp(proxy string) (*database.ProxyIPInfo, error) {
	// 写入数据库
	ipinfo := &database.ProxyIPInfo{IpPort: proxy}

	client := &httpclient.HttpClient{ProxyScheme: "http", ProxyIp: proxy, DialTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second}
	client, err := client.NewClient()
	t := time.Now() // get current time
	resp, err := client.Get("http://httpbin.org/ip")
	elapsed := time.Since(t)
	ipinfo.LastCheckTime = function.FormatTime(t)
	ipinfo.Speed = elapsed.Seconds()
	if err != nil {
		ipinfo.Results = "0"
		fmt.Printf("[时长: %f] IP[%s]失效[%s]\n", elapsed.Seconds(), proxy, err.Error())
	} else {
		if resp != nil {
			if resp.StatusCode == http.StatusOK {
				fmt.Printf("[时长: %f] 检测IP[%s]成功\n", elapsed.Seconds(), proxy)
				ipinfo.Results = "1"
			} else {
				ipinfo.Results = "0"
				fmt.Printf("[时长: %f] IP[%s]失效[%d]\n", elapsed.Seconds(), proxy, resp.StatusCode)
			}
		} else {
			ipinfo.Results = "0"
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
			time.Sleep(time.Duration(100) * time.Millisecond)
			go func(proxy string) {
				ipinfo, _ := checkProxyIp(proxy)
				store.Put(ipinfo)
			}(ip)
		}()
	}
}
