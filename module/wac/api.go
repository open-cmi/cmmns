package wac

import (
	"errors"
	"fmt"
	"io"
	"net/netip"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/open-cmi/cmmns/essential/logger"
)

var globalSeq int

func ApplyNginxBlackConf(m *Model) error {
	blkConf := path.Join(gConf.NginxPath, "ngx_blacklist.conf")
	wf, err := os.OpenFile(blkConf, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		logger.Errorf("open %s for writing failed\n", blkConf)
		return err
	}
	defer wf.Close()
	lines := strings.Split(m.RawBlacklist, "\n")
	for _, line := range lines {
		line = strings.Trim(line, "\t ")
		if line == "" {
			continue
		}
		content := fmt.Sprintf("deny %s;\n", line)
		_, err = wf.WriteString(content)
	}
	return err
}

func ApplyNginxWhiteConf(m *Model) error {
	whiteConf := path.Join(gConf.NginxPath, "ngx_whitelist.conf")
	wf, err := os.OpenFile(whiteConf, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		logger.Errorf("open %s for writing failed\n", whiteConf)
		return err
	}
	defer wf.Close()
	lines := strings.Split(m.RawWhitelist, "\n")
	for _, line := range lines {
		line = strings.Trim(line, "\t ")
		if line == "" {
			continue
		}
		content := fmt.Sprintf("allow %s;\n", line)
		_, err = wf.WriteString(content)
	}
	return err
}

func ApplyNginxConf(m *Model) error {
	var err error
	if m.Mode == "blacklist" {
		err = ApplyNginxBlackConf(m)
	} else {
		err = ApplyNginxWhiteConf(m)
	}
	if err != nil {
		logger.Errorf("write %s conf failed\n", m.Mode)
		return err
	}

	nginxConf := filepath.Join(gConf.NginxPath, "nginx.conf")
	rf, err := os.Open(nginxConf)
	if err != nil {
		logger.Errorf("open %s failed\n", nginxConf)
		return err
	}
	defer rf.Close()
	content, err := io.ReadAll(rf)
	if err != nil {
		logger.Errorf("read %s failed\n", nginxConf)
		return err
	}
	// 判断配置是否存在
	if !strings.Contains(string(content), "deny all;") && !strings.Contains(string(content), "allow all;") {
		// 两个配置都不存在，相当于新添加
		index := strings.Index(string(content), "http {")
		if index == -1 {
			return errors.New("nginx conf invalid: not contains location /")
		}
		for content[index] != '\n' {
			index++
		}
		newContent := string(content)[0 : index+1]
		newContent += "    include ngx_blacklist.conf;\n    allow all;\n"
		newContent += string(content)[index+1:]
		// 写文件
		wf, err := os.OpenFile(nginxConf, os.O_RDWR|os.O_TRUNC, 0644)
		if err != nil {
			logger.Errorf("open %s for writing failed\n", nginxConf)
			return err
		}
		_, err = wf.WriteString(newContent)
		return err
	}

	// 配置已经存在
	if m.Mode == "blacklist" {
		// 如果之前是白名单deny all，则替换成黑名单allow all
		newContent := strings.Replace(string(content), "deny all;", "allow all;", -1)
		newContent = strings.Replace(newContent, "include ngx_whitelist.conf", "include ngx_blacklist.conf", -1)
		// 写文件
		wf, err := os.OpenFile(nginxConf, os.O_RDWR|os.O_TRUNC, 0644)
		if err != nil {
			logger.Errorf("open %s for writing failed\n", nginxConf)
			return err
		}
		_, err = wf.WriteString(newContent)
		return err
	}
	// 白名单
	// 如果之前是黑名单deny all，则替换成黑名单deny all
	newContent := strings.Replace(string(content), "allow all;", "deny all;", -1)
	newContent = strings.Replace(newContent, "include ngx_blacklist.conf", "include ngx_whitelist.conf", -1)
	// 写文件
	wf, err := os.OpenFile(nginxConf, os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		logger.Errorf("open %s for writing failed\n", nginxConf)
		return err
	}
	_, err = wf.WriteString(newContent)
	if err != nil {
		logger.Errorf("write %s failed\n", nginxConf)
		return err
	}

	//
	if gConf.Reload != "" {
		cmd := exec.Command("bash", "-c", gConf.Reload)
		err = cmd.Run()
	}
	return err
}

func GetWAC() Model {
	m := Get()
	if m == nil {
		m = New()
		m.Mode = "blacklist"
	}
	return *m
}

func SetWAC(req *SetRequest) error {
	m := Get()
	if m == nil {
		m = New()
	}

	globalSeq++
	m.Enable = req.Enable
	m.Mode = req.Mode
	m.RawWhitelist = req.RawWhitelist
	m.RawBlacklist = req.RawBlacklist
	m.Seq = globalSeq

	if gConf.NginxPath != "" {
		// 保存之前先应用nginx配置，nginx应用成功再保存
		err := ApplyNginxConf(m)
		if err != nil {
			logger.Errorf("wac apply nginx conf failed: %s\n", err.Error())
			return err
		}
	}

	err := m.Save()
	if err != nil {
		logger.Errorf("wac save failed: %s\n", err.Error())
		return err
	}

	return err
}

func IPIsInList(rawList string, ip string) bool {
	ipaddr, err := netip.ParseAddr(ip)
	if err != nil {
		return false
	}

	lines := strings.Split(rawList, "\n")
	for _, line := range lines {
		line = strings.Trim(line, "\t ")
		if line == "" {
			continue
		}
		if strings.Contains(line, "/") {
			// prefix
			p, err := netip.ParsePrefix(line)
			if err != nil {
				continue
			}
			if p.Contains(ipaddr) {
				return true
			}
		} else {
			// addr
			a, err := netip.ParseAddr(line)
			if err != nil {
				continue
			}
			if a.Compare(ipaddr) == 0 {
				return true
			}
		}
	}
	return false
}

func CheckPermit(ip string) bool {

	if globalModel == nil || globalModel.Seq != globalSeq {
		globalModel = Get()
	}
	// 如果无配置，默认是允许的
	if globalModel == nil {
		return true
	}

	// 未开启，默认放行
	if !globalModel.Enable {
		return true
	}

	if globalModel.Mode == "blacklist" {
		// 默认通过
		return !IPIsInList(globalModel.RawBlacklist, ip)
	}
	// 如果是白名单，则在列表中放行，不在则deny
	return IPIsInList(globalModel.RawWhitelist, ip)
}
