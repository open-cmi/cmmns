package whitelist

import (
	"errors"
	"time"

	"github.com/open-cmi/cmmns/module/nginxconf"
)

func AddWhitelist(address string) error {
	blk := Get(address)
	if blk != nil {
		return errors.New("address is existed")
	}
	blk = New()
	blk.Address = address
	blk.Timestamp = time.Now().Unix()
	err := blk.Save()
	if err != nil {
		return err
	}
	// 这里重新应用nginx配置
	_, lists, err := List(nil)
	if err != nil {
		return err
	}
	var whitelists []string = []string{}
	for _, b := range lists {
		whitelists = append(whitelists, b.Address)
	}
	err = nginxconf.ApplyNginxWhiteConf(whitelists)
	return err
}

func DelWhitelist(address string) error {
	blk := Get(address)
	if blk == nil {
		return errors.New("address is not existed")
	}
	err := blk.Remove()
	if err != nil {
		return err
	}
	// 这里重新应用nginx配置
	_, lists, err := List(nil)
	if err != nil {
		return err
	}
	var whitelist []string = []string{}
	for _, b := range lists {
		whitelist = append(whitelist, b.Address)
	}
	err = nginxconf.ApplyNginxBlackConf(whitelist)
	return err
}