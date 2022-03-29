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
