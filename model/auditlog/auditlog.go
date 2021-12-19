package auditlog

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/open-cmi/cmmns/storage/db"
	"github.com/open-cmi/cmmns/utils"

	commonmsg "github.com/open-cmi/cmmns/msg/common"
)

type Model struct {
	ID        string `json:"id"`
	IP        string `json:"ip"`
	Type      int    `json:"type"`
	Username  string `json:"username"`
	Action    string `json:"action"`
	Timestamp int    `json:"timestamp"`
}

// List list
func List(p *commonmsg.RequestQuery) (int, []Model, error) {
	dbsql := db.GetDB()

	var logs []Model = []Model{}

	countClause := fmt.Sprintf("select count(*) from audit_log")
	whereClause := utils.BuildSQLClause(p)
	countClause += whereClause
	row := dbsql.QueryRow(countClause)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, logs, errors.New("get count failed")
	}

	queryClause := fmt.Sprintf(`select id,ip,type,username,action,timestamp from audit_log`)
	queryClause += whereClause
	rows, err := dbsql.Query(queryClause)
	if err != nil {
		// 没有的话，也不需要报错
		return count, logs, nil
	}

	for rows.Next() {
		var item Model
		err := rows.Scan(&item.ID, &item.IP, &item.Type, &item.Username, &item.Action, &item.Timestamp)
		if err != nil {
			break
		}

		logs = append(logs, item)
	}
	return count, logs, err
}

func InsertLog(ip string, username string, logtype int, action string) error {
	timestamp := time.Now().Unix()
	id := uuid.New().String()
	insertClause := fmt.Sprintf(`insert into audit_log(id, type, ip, username, action, timestamp) 
		values('%s', %d, '%s', '%s', '%s', %d)`, id, logtype, ip, username, action, timestamp)

	sqldb := db.GetDB()
	_, err := sqldb.Exec(insertClause)
	return err
}
