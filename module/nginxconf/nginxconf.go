package nginxconf

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/open-cmi/cmmns/essential/logger"
)

func ApplyNginxBlackConf(blacklists []string) error {
	blkConf := path.Join(gConf.Path, "ngx_blacklist.conf")
	wf, err := os.OpenFile(blkConf, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		logger.Errorf("open %s for writing failed\n", blkConf)
		return err
	}
	defer wf.Close()
	for _, line := range blacklists {
		line = strings.Trim(line, "\t ")
		if line == "" {
			continue
		}
		content := fmt.Sprintf("deny %s;\n", line)
		_, err = wf.WriteString(content)
	}
	if gConf.Reload != "" {
		cmd := exec.Command("bash", "-c", gConf.Reload)
		err = cmd.Run()
	}
	return err
}

func ApplyNginxWhiteConf(whitelists []string) error {
	whiteConf := path.Join(gConf.Path, "ngx_whitelist.conf")
	wf, err := os.OpenFile(whiteConf, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		logger.Errorf("open %s for writing failed\n", whiteConf)
		return err
	}
	defer wf.Close()
	for _, line := range whitelists {
		line = strings.Trim(line, "\t ")
		if line == "" {
			continue
		}
		content := fmt.Sprintf("allow %s;\n", line)
		_, err = wf.WriteString(content)
	}
	if gConf.Reload != "" {
		cmd := exec.Command("bash", "-c", gConf.Reload)
		err = cmd.Run()
	}
	return err
}

func applyAccessControl(mode string) error {
	var err error

	if gConf.Path == "" {
		return errors.New("nginx path is not set, please set it in config file first")
	}

	nginxConf := filepath.Join(gConf.Path, "nginx.conf")
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
		if mode == "blacklist" {
			newContent += "    include ngx_blacklist.conf;\n    allow all;\n"
		} else {
			newContent += "    include ngx_whitelist.conf;\n    deny all;\n"
		}
		newContent += string(content)[index+1:]
		// 写文件
		wf, err := os.OpenFile(nginxConf, os.O_RDWR|os.O_TRUNC, 0644)
		if err != nil {
			logger.Errorf("open %s for writing failed\n", nginxConf)
			return err
		}
		_, err = wf.WriteString(newContent)
		wf.Close()
		return err
	}

	// 配置已经存在
	if mode == "blacklist" {
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
		if err != nil {
			logger.Errorf("write %s failed\n", nginxConf)
			return err
		}
		wf.Close()
	} else {
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
		wf.Close()
	}
	return err
}

func ApplyNginxAccessControl(mode string) error {
	err := applyAccessControl(mode)
	if err != nil {
		return err
	}
	//
	if gConf.Reload != "" {
		cmd := exec.Command("bash", "-c", gConf.Reload)
		err = cmd.Run()
	}
	return err
}

func CancelNginxAccessControl() error {
	var err error

	if gConf.Path == "" {
		logger.Errorf("nginx path is not set, please set it in config file first\n")
		return nil
	}

	nginxConf := filepath.Join(gConf.Path, "nginx.conf")
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
	newContent := strings.Replace(string(content), "    include ngx_blacklist.conf;\n    allow all;\n", "", -1)
	newContent = strings.Replace(newContent, "    include ngx_whitelist.conf;\n    deny all;\n", "", -1)
	// 写文件
	wf, err := os.OpenFile(nginxConf, os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		logger.Errorf("open %s for writing failed\n", nginxConf)
		return err
	}

	_, err = wf.WriteString(newContent)
	if err != nil {
		wf.Close()
		logger.Errorf("write %s failed\n", nginxConf)
		return err
	}
	wf.Close()

	//
	if gConf.Reload != "" {
		cmd := exec.Command("bash", "-c", gConf.Reload)
		err = cmd.Run()
	}
	return err
}

func ApplyServicePort(httpPort int, httpsPort int) error {
	var err error

	m := GetServiceModel()
	if m == nil {
		m = NewServiceModel()
		m.HttpPort = 80
		m.HttpsPort = 443
	}

	if gConf.Path == "" {
		logger.Errorf("nginx path is not set, please set it in config file first\n")
		return nil
	}
	nginxConf := filepath.Join(gConf.Path, "nginx.conf")
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
	var newContent string
	if strings.Contains(string(content), "listen 80;") {
		newHttpString := fmt.Sprintf("listen %d;", httpPort)
		newContent = strings.Replace(string(content), "listen 80;", newHttpString, -1)
	} else {
		oldHttpString := fmt.Sprintf("listen %d;", m.HttpPort)
		newHttpString := fmt.Sprintf("listen %d;", httpPort)
		newContent = strings.Replace(string(content), oldHttpString, newHttpString, -1)
	}
	if strings.Contains(newContent, "listen 443 ssl;") {
		newHttpsString := fmt.Sprintf("listen %d ssl;", httpsPort)
		newContent = strings.Replace(newContent, "listen 443 ssl;", newHttpsString, -1)

		if httpsPort != 443 {
			newRedirectString := fmt.Sprintf("return 301 https://$host:%d$request_uri;", httpsPort)
			newContent = strings.Replace(newContent, "return 301 https://$host$request_uri;", newRedirectString, -1)
		}
	} else {
		oldHttpsString := fmt.Sprintf("listen %d ssl;", m.HttpsPort)
		newHttpsString := fmt.Sprintf("listen %d ssl;", httpsPort)
		newContent = strings.Replace(newContent, oldHttpsString, newHttpsString, -1)

		var newRedirectString string
		oldRedirectString := fmt.Sprintf("return 301 https://$host:%d$request_uri;", m.HttpsPort)

		if httpsPort == 443 {
			newRedirectString = "return 301 https://$host$request_uri;"
		} else {
			newRedirectString = fmt.Sprintf("return 301 https://$host:%d$request_uri;", httpsPort)
		}
		newContent = strings.Replace(newContent, oldRedirectString, newRedirectString, -1)
	}

	// 写文件
	wf, err := os.OpenFile(nginxConf, os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		logger.Errorf("open %s for writing failed\n", nginxConf)
		return err
	}
	defer wf.Close()
	_, err = wf.WriteString(newContent)
	if err != nil {
		logger.Errorf("write %s failed\n", nginxConf)
		return err
	}
	m.HttpPort = httpPort
	m.HttpsPort = httpsPort
	err = m.Save()
	if err != nil {
		logger.Errorf("nginx conf service model save failed: %s\n", err.Error())
	}
	//
	if gConf.Reload != "" {
		cmd := exec.Command("bash", "-c", gConf.Reload)
		err = cmd.Run()
	}
	return err
}
