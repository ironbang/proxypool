package spider

import (
	"strings"
)

func SplitIPPort(ip_port string) (ip, port string) {
	ip = ""
	port = ""
	ip_port = strings.TrimSpace(ip_port)
	r := strings.Split(ip_port, ":")
	if len(r) == 2 {
		ip = r[0]
		port = r[1]
	}
	return
}
