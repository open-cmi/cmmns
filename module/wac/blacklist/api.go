package blacklist

import (
	"errors"
	"time"

	"github.com/open-cmi/cmmns/module/nginxconf"
	"github.com/open-cmi/gobase/essential/i18n"
)

func AddBlacklist(address string) error {
	blk := Get(address)
	if blk != nil {
		return errors.New((i18n.Sprintf("address %s is existing", address)))
	}

	blk = New()
	blk.Address = address
	blk.Timestamp = time.Now().Unix()
	err := blk.Save()
	if err != nil {
		return err
	}
	// 这里重新应用nginx配置
	lists, err := ListAll()
	if err != nil {
		return err
	}
	var blacklist []string = []string{}
	for _, b := range lists {
		blacklist = append(blacklist, b.Address)
	}
	err = nginxconf.ApplyNginxBlackConf(blacklist)
	return err
}

func DelBlacklist(address string) error {
	blk := Get(address)
	if blk == nil {
		return errors.New("address is not existing")
	}
	err := blk.Remove()
	if err != nil {
		return err
	}
	// 这里重新应用nginx配置
	lists, err := ListAll()
	if err != nil {
		return err
	}
	var blacklist []string = []string{}
	for _, b := range lists {
		blacklist = append(blacklist, b.Address)
	}
	err = nginxconf.ApplyNginxBlackConf(blacklist)
	return err
}
