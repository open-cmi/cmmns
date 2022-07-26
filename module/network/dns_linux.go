package network

import (
	"bufio"
	"bytes"
	"os/exec"
	"strings"
)

func GetDNS() []string {
	var dns []string
	cmd := exec.Command("systemd-resolve", "--status")
	output, err := cmd.Output()
	if err != nil {
		// 如果systemd-resolve执行失败，则从resolve文件中取
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
	rd := bufio.NewReader(bytes.NewBuffer(output))
	for linebyte, _, err := rd.ReadLine(); err == nil; linebyte, _, err = rd.ReadLine() {
		line := strings.Trim(string(linebyte), "\r\n\t ")
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "DNS Servers:") {
			tmp := strings.TrimPrefix(line, "DNS Servers:")
			tmp = strings.Trim(tmp, "\r\n\t ")
			dns = append(dns, tmp)
		}
	}
	return dns
}
