package checkip

import (
	"fmt"
	"github.com/ironbang/httpclient"
	"github.com/ironbang/proxypool/common/config"
	"github.com/ironbang/proxypool/common/function"
	"github.com/ironbang/proxypool/database/struct_"
	"net/http"
	"time"
)

var maxProxyChan = config.GetChanelMax("TransferProxyIP")
var proxyChan chan *struct_.ProxyIPInfo

func init() {
	proxyChan = make(chan *struct_.ProxyIPInfo, maxProxyChan)
}

func checkProxyIp(proxy *struct_.ProxyIPInfo) (*struct_.ProxyIPInfo, error) {
	client := &httpclient.HttpClient{ProxyScheme: "http", ProxyIp: proxy.IpPort, DialTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second}
	client, err := client.NewClient()
	t := time.Now() // get current time
	resp, err := client.Get("http://httpbin.org/ip", make(map[string]string))
	elapsed := time.Since(t)
	proxy.LastCheckTime = function.FormatTime(t)
	proxy.Speed = elapsed.Seconds()
	if err != nil {
		proxy.Results = proxy.Results + "0"
		fmt.Printf("[%s] [时长: %f] IP[%s]失效[%s]\n", function.FormatTime(time.Now()), elapsed.Seconds(), proxy.IpPort, err.Error())
	} else {
		if resp != nil {
			if resp.StatusCode == http.StatusOK {
				fmt.Printf("[%s] [时长: %f] 检测IP[%s]成功\n", function.FormatTime(time.Now()), elapsed.Seconds(), proxy.IpPort)
				proxy.Results = proxy.Results + "1"
			} else {
				proxy.Results = proxy.Results + "0"
				fmt.Printf("[%s] [时长: %f] IP[%s]失效[%d]\n", function.FormatTime(time.Now()), elapsed.Seconds(), proxy.IpPort, resp.StatusCode)
			}
		} else {
			proxy.Results = proxy.Results + "0"
			fmt.Printf("[%s] [时长: %f] IP[%s]失效[%d]\n", function.FormatTime(time.Now()), elapsed.Seconds(), proxy.IpPort, resp.StatusCode)
		}
	}

	proxy.Checked = true
	return proxy, nil
}

func PutProxy(proxy *struct_.ProxyIPInfo) {
	for len(proxyChan) >= maxProxyChan {
		time.Sleep(500 * time.Millisecond)
	}
	proxyChan <- proxy
}

func CheckIp() {
	maxPool := config.GetChanelMax("CheckIpCoroutinePool")
	pool := make(chan bool, maxPool)
	fmt.Println("校验模块...")
	for {
		for len(pool) >= maxPool {
			time.Sleep(10 * time.Second)
		}
		proxy := <-proxyChan
		go func(proxy *struct_.ProxyIPInfo) {
			pool <- true
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
				}
				<-pool
			}()
			var err error
			proxy, err = checkProxyIp(proxy)
			if err == nil {
				proxy.Update()
			}
		}(proxy)
	}
}

func CheckStore() {
	fmt.Println("数据库中IP检测模块启动...")
	for {
		// start := time.Now()
		_ = func() int {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
				}
			}()
			ips, err := struct_.GetAll()
			if err != nil {
				fmt.Println(err.Error())
				return 0
			}
			fmt.Printf("数据库中共有%d个IP需要进行检测\n", len(ips))
			for _, ip := range ips {
				PutProxy(ip)
			}
			return len(ips)
		}()
		// fmt.Printf("共检测%d个代理IP，共耗时%0.2f秒",total,time.Since(start).Seconds())
		time.Sleep(5 * time.Millisecond)
	}
}
