package licmng

import (
	"fmt"
	"time"

	"github.com/open-cmi/gobase/essential/sqldb"
)

type Statistics struct {
	Total        int            `json:"total"`
	VersionStats map[string]int `json:"version_stats"`
	ModelStats   map[string]int `json:"model_stats"`
}

type MonthlyStat struct {
	Month int64 `json:"month" db:"month"`
	Count int   `json:"count" db:"count"`
}

func GetMonthlyStatistics(months int) ([]MonthlyStat, error) {
	db := sqldb.GetDB()

	// 计算起始时间（当前月份往前推 N 个月）
	// 为了包含本月完整数据，通常从该月 1 号开始算
	now := time.Now()
	startTime := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).AddDate(0, -months+1, 0).Unix()

	var query string
	var args []interface{}
	args = append(args, startTime)

	if db.DriverName() == "postgres" {
		query = `SELECT CAST(EXTRACT(EPOCH FROM date_trunc(
            	'month',
            	to_timestamp(created_time) AT TIME ZONE 'Asia/Shanghai'
        		)
			) AS BIGINT
		) AS month,
		COUNT(*) AS count
		FROM license
		WHERE created_time >= $1
		GROUP BY month
		ORDER BY month ASC;`
	} else {
		// 默认适配 sqlite
		query = `SELECT CAST(strftime('%s', datetime(created_time, 'unixepoch', 'localtime', 'start of month')) AS BIGINT) as month, count(*) as count 
				 FROM license WHERE created_time >= ? 
				 GROUP BY month ORDER BY month ASC`
	}

	var stats []MonthlyStat
	err := db.Select(&stats, query, args...)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func GetStatistics(year int, product string) (*Statistics, error) {
	db := sqldb.GetDB()

	// 基础查询，按版本和模型进行分组统计
	query := `SELECT version, model, count(*) as count FROM license WHERE 1=1`
	var args []interface{}
	argIdx := 1

	if product != "" {
		query += fmt.Sprintf(" AND prod = $%d", argIdx)
		args = append(args, product)
		argIdx++
	}

	if year > 0 {
		// 假设使用 Unix 时间戳存储 created_time
		startTime := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
		endTime := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
		query += fmt.Sprintf(" AND created_time >= $%d AND created_time < $%d", argIdx, argIdx+1)
		args = append(args, startTime, endTime)
	}

	query += " GROUP BY version, model"

	rows, err := db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := Statistics{
		VersionStats: make(map[string]int),
		ModelStats:   make(map[string]int),
	}

	for rows.Next() {
		var item struct {
			Version string `db:"version"`
			Model   string `db:"model"`
			Count   int    `db:"count"`
		}
		if err := rows.StructScan(&item); err != nil {
			return nil, err
		}

		stats.Total += item.Count
		stats.VersionStats[item.Version] += item.Count
		stats.ModelStats[item.Model] += item.Count
	}

	return &stats, nil
}
