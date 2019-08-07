package spider

import (
	"fmt"
	"github.com/ironbang/requests"
	"regexp"
)

func IP89Spider(sysChan chan<- string) {
	fmt.Println("开始爬取89ip")
	url := `http://www.89ip.cn/tqdl.html?num=9999&address=&kill_address=&port=&kill_port=&isp=`
	// url = `http://www.baidu.com`
	resp, err := requests.Get(url)
	if err != nil {
		return
	}
	ip_reg := regexp.MustCompile(`((25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\:([0-9]+)`)
	ips := ip_reg.FindAllString(resp.Text(), -1)
	for _, ip := range ips {
		//ip = "http://" + ip
		sysChan <- ip
	}
}
