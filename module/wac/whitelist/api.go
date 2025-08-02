package whitelist

import (
	"errors"
	"time"

	"github.com/open-cmi/cmmns/module/nginxconf"
	"github.com/open-cmi/gobase/essential/i18n"
	"github.com/open-cmi/gobase/essential/logger"
	"github.com/open-cmi/gobase/essential/sqldb"
	"github.com/open-cmi/gobase/pkg/goparam"
)

func AddWhitelist(address string) error {
	blk := Get(address)
	if blk != nil {
		return errors.New(i18n.Sprintf("address %s is existing", address))
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
	var whitelist []string = []string{}
	for _, b := range lists {
		whitelist = append(whitelist, b.Address)
	}
	err = nginxconf.ApplyNginxWhiteConf(whitelist)
	return err
}

func QueryList(query *goparam.Param) (int, []Whitelist, error) {
	db := sqldb.GetDB()

	var users []Whitelist = []Whitelist{}
	countClause := "select count(*) from wac_whitelist"
	row := db.QueryRow(countClause)

	var count int
	err := row.Scan(&count)
	if err != nil {
		logger.Errorf("wac_whitelist list count failed, %s\n", err.Error())
		return 0, users, errors.New("list count failed")
	}

	queryClause := `select * from wac_whitelist`
	finalClause := goparam.BuildFinalClause(query, []string{"timestamp"})
	queryClause += finalClause
	rows, err := db.Queryx(queryClause)
	if err != nil {
		// 没有的话，也不需要报错
		return count, users, nil
	}
	defer rows.Close()
	for rows.Next() {
		var item Whitelist
		err := rows.StructScan(&item)
		if err != nil {
			logger.Errorf("whitelist struct scan failed %s\n", err.Error())
			break
		}

		users = append(users, item)
	}
	return count, users, err
}

func ListAll() ([]Whitelist, error) {
	db := sqldb.GetDB()

	var users []Whitelist = []Whitelist{}

	queryClause := `select * from wac_whitelist`
	rows, err := db.Queryx(queryClause)
	if err != nil {
		// 没有的话，也不需要报错
		return users, nil
	}
	defer rows.Close()
	for rows.Next() {
		var item Whitelist
		err := rows.StructScan(&item)
		if err != nil {
			logger.Errorf("whitelist struct scan failed %s\n", err.Error())
			break
		}

		users = append(users, item)
	}
	return users, err
}
