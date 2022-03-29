package network

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func GetDNS() []string {
	var dns []string
	rf, err := os.Open("/etc/resolv.conf")
	if err != nil {
		return dns
	}
	rd := bufio.NewReader(rf)
	for linebyte, _, err := rd.ReadLine(); err == nil; linebyte, _, err = rd.ReadLine() {
		line := string(linebyte)
		fmt.Println(line)
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "nameserver") {
			tmp := strings.TrimPrefix(line, "nameserver")
			tmp = strings.Trim(tmp, "\r\n\t ")
			dns = append(dns, tmp)
		}
	}
	return dns
}
