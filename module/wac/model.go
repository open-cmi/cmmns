package wac

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/netip"
	"strings"

	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/sqldb"
)

// web access control

var gWebAccessControl WebAccessControl

type WebAccessControl struct {
	Mode          string
	PermitAddrs   []netip.Addr
	PermitPrefixs []netip.Prefix
	DenyAddrs     []netip.Addr
	DenyPrefixs   []netip.Prefix
}

type Model struct {
	Enable       bool   `json:"enable" db:"enable"`
	Mode         string `json:"mode" db:"mode"` // blacklist or whitelist
	RawWhitelist string `json:"raw_whitelist" db:"raw_whitelist"`
	RawBlacklist string `json:"raw_blacklist" db:"raw_blacklist"`
	isNew        bool
}

func (m *Model) Key() string {
	return "web-access-control"
}

func (m *Model) Value() string {
	v, _ := json.Marshal(m)
	return string(v)
}

func (m *Model) Save() error {
	db := sqldb.GetConfDB()

	if m.isNew {
		// 存储到数据库
		columns := []string{"key", "value"}
		values := []string{"$1", "$2"}

		insertClause := fmt.Sprintf("insert into k_v_table(%s) values(%s)",
			strings.Join(columns, ","), strings.Join(values, ","))

		logger.Debugf("start to exec sql clause: %s\n", insertClause)
		_, err := db.Exec(insertClause, m.Key(), m.Value())
		if err != nil {
			logger.Errorf("create model failed: %s", err.Error())
			return errors.New("create model failed")
		}
		m.isNew = false
	} else {
		updateClause := "update k_v_table set value=$1 where key=$2"
		logger.Debugf("start to exec sql clause: %s", updateClause)
		_, err := db.Exec(updateClause, m.Value(), m.Key())
		if err != nil {
			logger.Errorf("update model failed: %s", err.Error())
			return errors.New("update model failed")
		}
	}

	return nil
}

func New() *Model {
	return &Model{
		isNew: true,
	}
}

func Get() *Model {
	db := sqldb.GetConfDB()
	var m Model
	queryClause := `select value from k_v_table where key=$1`
	row := db.QueryRowx(queryClause, m.Key())
	if row == nil {
		return nil
	}
	var value string
	err := row.Scan(&value)
	if err != nil {
		logger.Infof("wac scan model failed: %s\n", err.Error())
		return nil
	}

	err = json.Unmarshal([]byte(value), &m)
	if err != nil {
		logger.Errorf("wac json unmarshal failed: %s\n", err.Error())
		return nil
	}

	return &m
}

func (m *Model) ConvertoWAC() WebAccessControl {
	var wac WebAccessControl
	wac.Mode = m.Mode

	var arrs []string
	lines := strings.Split(m.RawWhitelist, "\n")
	for _, line := range lines {
		tmp := strings.Split(line, ",")
		arrs = append(arrs, tmp...)
	}
	for _, ar := range arrs {
		ar = strings.Trim(ar, "\t ")
		if strings.Contains(ar, "/") {
			// prefix
			p, err := netip.ParsePrefix(ar)
			if err != nil {
				continue
			}
			wac.PermitPrefixs = append(wac.PermitPrefixs, p)
		} else {
			// addr
			a, err := netip.ParseAddr(ar)
			if err != nil {
				continue
			}
			wac.PermitAddrs = append(wac.PermitAddrs, a)
		}
	}

	arrs = []string{}
	lines = strings.Split(m.RawBlacklist, "\n")
	for _, line := range lines {
		tmp := strings.Split(line, ",")
		arrs = append(arrs, tmp...)
	}
	for _, ar := range arrs {
		ar = strings.Trim(ar, "\t ")
		if strings.Contains(ar, "/") {
			// prefix
			p, err := netip.ParsePrefix(ar)
			if err != nil {
				continue
			}
			wac.DenyPrefixs = append(wac.DenyPrefixs, p)
		} else {
			// addr
			a, err := netip.ParseAddr(ar)
			if err != nil {
				continue
			}
			wac.DenyAddrs = append(wac.DenyAddrs, a)
		}
	}
	return wac
}
