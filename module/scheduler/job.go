package scheduler

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/open-cmi/cmmns/common/def"
	"github.com/open-cmi/cmmns/common/goparam"
	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/rdb"
	"github.com/open-cmi/cmmns/essential/sqldb"
)

// state Init Pending Running Done Killed Retry

const (
	Init    = "Init"
	Running = "Running"
	Pending = "Pending"
	Done    = "Done"
	Killed  = "Killed"
	Retry   = "Retry"
)

// Job 用于sheduler存储调度使用
type Job struct {
	ID          string `json:"id" db:"id"`
	CronID      string `json:"cron_id" db:"cron_id"`
	Type        string `json:"type" db:"type"`
	RunType     string `json:"run_type" db:"run_type"`
	RunSpec     string `json:"run_spec" db:"run_spec"`
	SchedGroup  string `json:"sched_group" db:"sched_group"`
	SchedObject string `json:"sched_object" db:"sched_object"`
	Priority    int    `json:"priority" db:"priority"`
	State       string `json:"state" db:"state"`
	Count       int    `json:"count" db:"count"`
	Done        int    `json:"done" db:"done"`
	Content     string `json:"content,omitempty" db:"content"`
	Code        int    `json:"code" db:"code"`
	Msg         string `json:"msg" db:"msg"`
	Result      string `json:"result" db:"result"`
	StartedTime int64  `json:"started_time" db:"started_time"`
	StoppedTime int64  `json:"stopped_time" db:"stopped_time"`
	IsNew       bool   `json:"-"`
}

func (m *Job) Save() error {
	db := sqldb.GetConfDB()
	m.setCache()

	if m.IsNew {
		// 存储到数据库
		columns := goparam.GetColumn(*m, []string{})
		values := goparam.GetColumnInsertNamed(columns)

		insertClause := fmt.Sprintf("insert into job(%s) values(%s)",
			strings.Join(columns, ","), strings.Join(values, ","))

		logger.Debugf("start to exec sql clause: %s\n", insertClause)

		_, err := db.NamedExec(insertClause, m)
		if err != nil {
			logger.Errorf("create model failed: %s", err.Error())
			return errors.New("create model failed")
		}
		m.IsNew = false
	} else {
		columns := goparam.GetColumn(*m, []string{"id", "started_time"})

		var updates []string = []string{}
		for _, column := range columns {
			updates = append(updates, fmt.Sprintf(`%s=:%s`, column, column))
		}
		updateClause := fmt.Sprintf("update job set %s where id=:id", strings.Join(updates, ","))
		logger.Debugf("start to exec sql clause: %s", updateClause)
		_, err := db.NamedExec(updateClause, m)
		if err != nil {
			logger.Errorf("update job model failed: %s", err.Error())
			return errors.New("update model failed")
		}
	}

	return nil
}

func (m *Job) Remove() error {
	db := sqldb.GetConfDB()

	deleteClause := "delete from job where id=:id"
	_, err := db.NamedExec(deleteClause, m)
	if err != nil {
		return errors.New("delete model failed")
	}
	return nil
}

func (j *Job) setCache() error {
	cache := rdb.GetClient("job")

	key := fmt.Sprintf("job.%s", j.ID)
	jobMap := make(map[string]interface{})
	jobMap["id"] = j.ID
	jobMap["type"] = j.Type
	jobMap["cron_id"] = j.CronID
	jobMap["count"] = j.Count
	jobMap["run_type"] = j.RunType
	jobMap["run_spce"] = j.RunSpec
	jobMap["sched_group"] = j.SchedGroup
	jobMap["sched_object"] = j.SchedObject
	jobMap["priority"] = j.Priority
	jobMap["state"] = j.State
	jobMap["done"] = j.Done
	jobMap["content"] = j.Content
	jobMap["code"] = j.Code
	jobMap["msg"] = j.Msg
	jobMap["result"] = j.Result
	jobMap["started_time"] = j.StartedTime
	jobMap["stopped_time"] = j.StoppedTime

	_, err := cache.HSet(context.TODO(), key, jobMap).Result()
	if err != nil {
		return err
	}
	cache.Expire(context.TODO(), key, 5*time.Minute)
	return nil
}

func getCache(id string) *Job {
	cache := rdb.GetClient("job")

	key := fmt.Sprintf("job.%s", id)
	jobMap, err := cache.HGetAll(context.TODO(), key).Result()
	if err != nil {
		return nil
	}
	var job Job
	job.ID = jobMap["id"]
	job.Type = jobMap["type"]
	job.CronID = jobMap["cron_id"]
	job.Count, _ = strconv.Atoi(jobMap["count"])
	job.RunType = jobMap["run_type"]
	job.RunSpec = jobMap["run_spce"]
	job.SchedGroup = jobMap["sched_group"]
	job.SchedObject = jobMap["sched_object"]
	job.Priority, _ = strconv.Atoi(jobMap["priority"])
	job.State = jobMap["state"]
	job.Done, _ = strconv.Atoi(jobMap["done"])
	job.Content = jobMap["content"]
	job.Code, _ = strconv.Atoi(jobMap["code"])
	job.Msg = jobMap["msg"]
	job.Result = jobMap["result"]
	job.StartedTime, _ = strconv.ParseInt(jobMap["started_time"], 10, 64)
	job.StoppedTime, _ = strconv.ParseInt(jobMap["stopped_time"], 10, 64)
	return &job
}

func init() {
	rdb.Register("job", def.RDBJob)
}
