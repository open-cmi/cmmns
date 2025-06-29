package dnsmasq

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/open-cmi/gobase/essential/logger"
)

type DNSConf struct {
	PreferredDNS string
	AlternateDNS string
}

func ApplyResolveConf(preffer string, alter string) error {
	os.Remove("/etc/resolv.conf")
	file, err := os.OpenFile("/etc/resolv.conf", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return errors.New("open resolv.conf failed")
	}
	defer file.Close()

	content := ""
	if preffer != "" {
		content += fmt.Sprintf("nameserver %s\n", preffer)
	}

	if alter != "" {
		content += fmt.Sprintf("nameserver %s\n", alter)
	}

	_, err = file.WriteString(content)
	if err != nil {
		return errors.New("write to resolv.conf failed")
	}

	return nil
}

func SetDNS(preffered string, alternate string, listens []string) error {
	if gConf.Enable {
		logger.Debugf("dnsmasq is not enable, but someone set dns: %s %s\n", preffered, alternate)
		return nil
	}

	var content string = ""
	if preffered != "" {
		content += fmt.Sprintf("server=%s\n", preffered)
	}

	if alternate != "" {
		content += fmt.Sprintf("server=%s\n", alternate)
	}

	wf, err := os.OpenFile(gConf.ConfigFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		logger.Errorf("open dnsmasq conf failed: %s\n", err.Error())
		return err
	}
	defer wf.Close()

	wr := bufio.NewWriter(wf)
	wr.WriteString(content)
	//Flush将缓存的文件真正写入到文件中
	wr.Flush()

	ApplyResolveConf(preffered, alternate)
	err = RestartService()
	return err
}
