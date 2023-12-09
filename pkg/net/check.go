package net

import (
	"os/exec"
	"time"

	"github.com/go-ping/ping"
)

// PingCheck run as root
func PingCheck(address string, count int) *ping.Statistics {
	ping, err := ping.NewPinger(address)
	if err != nil {
		return nil
	}
	ping.Count = count                                    // ping次数
	ping.Timeout = time.Duration(3000 * time.Millisecond) // 超时时间为 3s
	ping.SetPrivileged(true)
	ping.Run() // blocks until finished
	return ping.Statistics()
}

func CurlCheck(domain string) (code string, err error) {
	args := []string{
		"-s",
		"-o",
		"/dev/null",
		"-w",
		"%{http_code}",
		domain,
	}
	cmd := exec.Command("curl", args...)
	o, err := cmd.Output()
	return string(o), err
}
