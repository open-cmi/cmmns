package syslog

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/open-cmi/gobase/essential/logger"
	"github.com/open-cmi/gobase/essential/sqldb"
	"github.com/open-cmi/gobase/pkg/ginparam"
)

// Model syslog record
type Model struct {
	ID        string `json:"id" db:"id"`
	LogTime   int64  `json:"log_time" db:"log_time"`
	Process   string `json:"process" db:"process"`
	PID       int    `json:"pid" db:"pid"`
	Hostname  string `json:"hostname" db:"hostname"`
	Level     string `json:"level" db:"level"`
	Message   string `json:"message" db:"message"`
	CreatedAt int64  `json:"created_at" db:"created_at"`
}

// Save inserts a syslog record into database
func (m *Model) Save() error {
	db := sqldb.GetDB()

	columns := sqldb.ExtractDBColumns(*m)
	values := sqldb.GetNamedInsertValuesClauses(columns)

	insertClause := fmt.Sprintf("insert into syslog(%s) values(%s)",
		strings.Join(columns, ","), values)

	_, err := db.NamedExec(insertClause, m)
	if err != nil {
		logger.Errorf("insert syslog failed: %s\n", err.Error())
	}
	return err
}

// QueryFilter filter options for querying syslog
type QueryFilter struct {
	Process   string
	Processes []string
	Keyword   string
	Level     string
	TimeStart int64
	TimeEnd   int64
}

// QueryList queries syslog records with pagination and filters
func QueryList(p *ginparam.ListOption, filter *QueryFilter) (int, []Model, error) {
	db := sqldb.GetDB()

	var logs []Model = []Model{}
	var paramnum int = 1
	var whereClause string
	var whereArgs []interface{}

	appendWhere := func(cond string, args ...interface{}) {
		if whereClause != "" {
			whereClause += " and "
		} else {
			whereClause += " where "
		}
		whereClause += cond
		whereArgs = append(whereArgs, args...)
	}

	if filter.Process != "" {
		appendWhere(fmt.Sprintf("process like %s", sqldb.LikePlaceHolder(paramnum)), filter.Process)
		paramnum++
	}

	if len(filter.Processes) > 0 {
		placeholders := make([]string, 0, len(filter.Processes))
		args := make([]interface{}, 0, len(filter.Processes))
		for _, proc := range filter.Processes {
			placeholders = append(placeholders, sqldb.PlaceHolder(paramnum))
			args = append(args, proc)
			paramnum++
		}
		appendWhere(fmt.Sprintf("process in (%s)", strings.Join(placeholders, ",")), args...)
	}

	if filter.Keyword != "" {
		appendWhere(fmt.Sprintf("message like %s", sqldb.LikePlaceHolder(paramnum)), filter.Keyword)
		paramnum++
	}

	if filter.Level != "" {
		appendWhere(fmt.Sprintf("level = %s", sqldb.PlaceHolder(paramnum)), filter.Level)
		paramnum++
	}

	if filter.TimeStart > 0 && filter.TimeEnd > 0 {
		appendWhere(
			fmt.Sprintf("log_time >= %s and log_time <= %s", sqldb.PlaceHolder(paramnum), sqldb.PlaceHolder(paramnum+1)),
			filter.TimeStart, filter.TimeEnd,
		)
		paramnum += 2
	}

	countClause := "select count(*) from syslog" + whereClause
	row := db.QueryRow(countClause, whereArgs...)

	var count int
	if err := row.Scan(&count); err != nil {
		logger.Errorf("get syslog count failed: %s\n", err.Error())
		return 0, logs, nil
	}

	queryClause := "select * from syslog" + whereClause
	queryClause += sqldb.BuildFinalClause(p, []string{"log_time", "created_at"})

	rows, err := db.Queryx(queryClause, whereArgs...)
	if err != nil {
		logger.Errorf("syslog queryx failed: %s\n", err.Error())
		return count, logs, nil
	}
	defer rows.Close()

	for rows.Next() {
		var item Model
		if err := rows.StructScan(&item); err != nil {
			logger.Errorf("syslog struct scan failed: %s\n", err.Error())
			break
		}
		logs = append(logs, item)
	}
	return count, logs, nil
}

// CleanExpired deletes syslog records older than keepDays
func CleanExpired(keepDays int) error {
	if keepDays <= 0 {
		return nil
	}
	db := sqldb.GetDB()
	cutoff := time.Now().AddDate(0, 0, -keepDays).Unix()
	_, err := db.Exec(fmt.Sprintf("delete from syslog where created_at < %s", sqldb.PlaceHolder(1)), cutoff)
	if err != nil {
		logger.Errorf("clean expired syslog failed: %s\n", err.Error())
	}
	return err
}

// RFC3164 syslog line pattern: "Nov 22 13:15:34 hostname process[1234]: message"
var rfc3164Regex = regexp.MustCompile(`^([A-Z][a-z]{2}\s+\d+\s+\d{2}:\d{2}:\d{2})\s+(\S+)\s+([^\s\[:]+)(?:\[(\d+)\])?:\s*(.*)$`)

// ParseLine parses a standard RFC3164 syslog line into Model
func ParseLine(line string) (*Model, bool) {
	matches := rfc3164Regex.FindStringSubmatch(strings.TrimSpace(line))
	if len(matches) != 6 {
		return nil, false
	}

	now := time.Now()
	// RFC3164 timestamp lacks year; attach current year and parse in local timezone
	tsStr := fmt.Sprintf("%d %s", now.Year(), matches[1])
	parsedTime, err := time.ParseInLocation("2006 Jan _2 15:04:05", tsStr, time.Local)
	if err != nil {
		parsedTime = now
	}

	pid := 0
	if matches[4] != "" {
		pid, _ = strconv.Atoi(matches[4])
	}

	return &Model{
		ID:        uuid.New().String(),
		LogTime:   parsedTime.Unix(),
		Hostname:  matches[2],
		Process:   matches[3],
		PID:       pid,
		Message:   matches[5],
		CreatedAt: now.Unix(),
	}, true
}

// CollectFromFile reads syslog from filePath starting at offset, filters by allowedProcesses,
// saves matched records to DB, and returns the new offset.
func CollectFromFile(filePath string, offset int64, allowedProcesses []string) (int64, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return offset, err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return offset, err
	}
	// File was truncated or rotated
	if info.Size() < offset {
		offset = 0
	}

	if _, err := f.Seek(offset, 0); err != nil {
		return offset, err
	}

	allowed := make(map[string]struct{}, len(allowedProcesses))
	for _, p := range allowedProcesses {
		if p != "" {
			allowed[p] = struct{}{}
		}
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		m, ok := ParseLine(line)
		if !ok {
			continue
		}
		if len(allowed) > 0 {
			if _, exists := allowed[m.Process]; !exists {
				continue
			}
		}
		if err := m.Save(); err != nil {
			logger.Errorf("save syslog record failed: %s\n", err.Error())
		}
	}

	newOffset, err := f.Seek(0, 1)
	if err != nil {
		return offset, err
	}
	return newOffset, scanner.Err()
}

var priorityLevelMap = map[string]string{
	"0": "emerg",
	"1": "alert",
	"2": "crit",
	"3": "err",
	"4": "warning",
	"5": "notice",
	"6": "info",
	"7": "debug",
}

// CollectFromJournalctl reads syslog entries from journalctl starting after cursor,
// filters by allowedProcesses, saves matched records to DB, and returns the new cursor.
func CollectFromJournalctl(cursor string, allowedProcesses []string) (string, error) {
	args := []string{"-o", "json", "--no-pager"}
	if cursor != "" {
		args = append(args, "--after-cursor="+cursor)
	} else {
		args = append(args, "-n", "1000")
	}
	for _, proc := range allowedProcesses {
		if proc != "" {
			args = append(args, "-t", proc)
		}
	}

	cmd := exec.Command("journalctl", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return cursor, err
	}
	if err := cmd.Start(); err != nil {
		return cursor, err
	}

	newCursor := cursor
	now := time.Now().Unix()
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		var entry map[string]interface{}
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			continue
		}

		if c, ok := entry["__CURSOR"].(string); ok && c != "" {
			newCursor = c
		}

		proc, _ := entry["SYSLOG_IDENTIFIER"].(string)
		if proc == "" {
			proc, _ = entry["_COMM"].(string)
		}

		msg, _ := entry["MESSAGE"].(string)
		if msg == "" {
			continue
		}

		hostname, _ := entry["_HOSTNAME"].(string)
		pidStr, _ := entry["_PID"].(string)
		pid, _ := strconv.Atoi(pidStr)

		prioStr, _ := entry["PRIORITY"].(string)
		level := priorityLevelMap[prioStr]

		logTime := now
		if tsStr, ok := entry["__REALTIME_TIMESTAMP"].(string); ok {
			if usec, err := strconv.ParseInt(tsStr, 10, 64); err == nil {
				logTime = usec / 1_000_000
			}
		}

		m := &Model{
			ID:        uuid.New().String(),
			LogTime:   logTime,
			Process:   proc,
			PID:       pid,
			Hostname:  hostname,
			Level:     level,
			Message:   msg,
			CreatedAt: now,
		}
		if err := m.Save(); err != nil {
			logger.Errorf("save journalctl syslog record failed: %s\n", err.Error())
		}
	}

	if err := scanner.Err(); err != nil {
		_ = cmd.Wait()
		return newCursor, err
	}
	if err := cmd.Wait(); err != nil {
		return newCursor, err
	}
	return newCursor, nil
}
