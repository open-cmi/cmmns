package auditlog

const (
	LoginType = iota
	OperationType
)

// InsertLog insert audit log
func InsertLog(ip string, username string, logtype int, action string, success bool) error {
	lr := NewLogRecord(ip, logtype, username, action, success)
	return lr.Save()
}
