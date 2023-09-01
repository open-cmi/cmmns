package auditlog

import (
	"time"

	"github.com/google/uuid"
	"github.com/open-cmi/cmmns/essential/sqldb"
)

const (
	LoginType = iota
	OperationType
	SystemType
	WarningType
)

// InsertLog insert audit log
func InsertLog(ip string, username string, logtype int, action string) error {
	timestamp := time.Now().Unix()
	id := uuid.New().String()

	// 这里应该通过tunnel来传递，不能直接写数据库，后续再调整
	insertClause := `insert into audit_log(id, type, ip, username, action, timestamp) 
		values($1, $2, $3, $4, $5, $6)`

	db := sqldb.GetConfDB()
	_, err := db.Exec(insertClause, id, logtype, ip, username, action, timestamp)
	return err
}
