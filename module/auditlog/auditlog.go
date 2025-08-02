package auditlog

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/open-cmi/gobase/essential/i18n"
	"github.com/open-cmi/gobase/essential/logger"
	"github.com/open-cmi/gobase/essential/sqldb"
	"github.com/open-cmi/gobase/pkg/goparam"
)

type Model struct {
	ID        string `json:"id" db:"id"`
	IP        string `json:"ip" db:"ip"`
	Type      int    `json:"type" db:"type"`
	Username  string `json:"username" db:"username"`
	Action    string `json:"action" db:"action"`
	Result    string `json:"result" db:"result"`
	Timestamp int64  `json:"timestamp" db:"timestamp"`
}

func (m *Model) Save() error {
	db := sqldb.GetDB()

	columns := goparam.GetColumn(*m, []string{})
	values := goparam.GetColumnInsertNamed(columns)

	insertClause := fmt.Sprintf("insert into audit_log(%s) values(%s)",
		strings.Join(columns, ","), strings.Join(values, ","))

	logger.Debugf("start to exec sql clause: %s\n", insertClause)

	_, err := db.NamedExec(insertClause, m)
	if err != nil {
		logger.Errorf("insert log failed: %s\n", err.Error())
	}
	return err
}

func NewLogRecord(ip string, logtype int, username string, action string, success bool) *Model {
	var result string
	if success {
		result = i18n.Sprintf("success")
	} else {
		result = i18n.Sprintf("fail")
	}

	return &Model{
		ID:        uuid.New().String(),
		Timestamp: time.Now().Unix(),
		IP:        ip,
		Type:      logtype,
		Username:  username,
		Action:    action,
		Result:    result,
	}
}

// List list
type QueryFilter struct {
	IP        string
	Username  string
	TimeStart int64
	TimeEnd   int64
}

func QueryList(p *goparam.Param, filter *QueryFilter) (int, []Model, error) {
	db := sqldb.GetDB()

	var logs []Model = []Model{}

	var paramnum int = 1
	var whereClause string
	var whereArgs []interface{}

	if filter.IP != "" {
		if whereClause != "" {
			whereClause += " and "
		} else {
			whereClause += " where "
		}
		whereClause += fmt.Sprintf(`ip like %s`, sqldb.LikePlaceHolder(paramnum))
		whereArgs = append(whereArgs, filter.IP)
		paramnum += 1
	}

	if filter.Username != "" {
		if whereClause != "" {
			whereClause += " and "
		} else {
			whereClause += " where "
		}
		whereClause += fmt.Sprintf(`username like %s`, sqldb.LikePlaceHolder(paramnum))
		whereArgs = append(whereArgs, filter.Username)
		paramnum += 1
	}

	if filter.TimeStart != 0 && filter.TimeEnd != 0 {
		if whereClause != "" {
			whereClause += " and "
		} else {
			whereClause += " where "
		}
		whereClause += fmt.Sprintf(`timestamp > %d and timestamp < %d`, paramnum, paramnum+1)
		whereArgs = append(whereArgs, filter.TimeStart, filter.TimeEnd)
		paramnum += 2
	}

	countClause := "select count(*) from audit_log"
	countClause += whereClause
	row := db.QueryRow(countClause, whereArgs...)

	var count int
	err := row.Scan(&count)
	if err != nil {
		// 这里只打印错误，不暴露错误，防止存在sql注入时给用于错误提示
		logger.Errorf("get count failed: %s\n", err.Error())
		return 0, logs, nil
	}

	queryClause := `select * from audit_log`
	queryClause += whereClause
	finalClause := goparam.BuildFinalClause(p, []string{"timestamp"})
	queryClause += finalClause
	rows, err := db.Queryx(queryClause, whereArgs...)
	if err != nil {
		logger.Errorf("audit log queryx failed: %s\n", err.Error())
		// 没有的话，也不需要报错
		return count, logs, nil
	}
	defer rows.Close()
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
