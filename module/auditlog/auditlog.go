package auditlog

import (
	"fmt"
	"strings"

	"github.com/open-cmi/cmmns/common/api"
	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/sqldb"
)

type Model struct {
	ID        string `json:"id" db:"id"`
	IP        string `json:"ip" db:"ip"`
	Type      int    `json:"type" db:"type"`
	Username  string `json:"username" db:"username"`
	Action    string `json:"action" db:"action"`
	Timestamp int    `json:"timestamp" db:"timestamp"`
}

// List list
func List(p *api.Option) (int, []Model, error) {
	db := sqldb.GetDB()

	var logs []Model = []Model{}

	countClause := "select count(*) from audit_log"
	whereClause, args := api.BuildWhereClause(p)
	countClause += whereClause
	row := db.QueryRow(countClause, args...)

	var count int
	err := row.Scan(&count)
	if err != nil {
		// 这里只打印错误，不暴露错误，防止存在sql注入时给用于错误提示
		logger.Errorf("get count failed: %s\n", err.Error())
		return 0, logs, nil
	}

	columns := api.GetColumn(Model{}, []string{})

	queryClause := fmt.Sprintf(`select %s from audit_log`, strings.Join(columns, ","))
	queryClause += whereClause
	finalClause := api.BuildFinalClause(p)
	queryClause += finalClause
	rows, err := db.Queryx(queryClause, args...)
	if err != nil {
		logger.Errorf("audit log queryx failed: %s\n", err.Error())
		// 没有的话，也不需要报错
		return count, logs, nil
	}

	for rows.Next() {
		var item Model
		err := rows.StructScan(&item)
		if err != nil {
			logger.Errorf("audit log struct scan failed: %s\n", err.Error())
			break
		}

		logs = append(logs, item)
	}
	return count, logs, err
}
