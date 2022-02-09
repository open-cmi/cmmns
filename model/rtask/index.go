package rtask

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/open-cmi/cmmns/logger"
	"github.com/open-cmi/cmmns/storage/db"
	"github.com/open-cmi/goutils/logutil"
)

// realtime task

var TaskPoolMutex sync.Mutex
var TaskPool map[string]*Task

type Task struct {
	ID          string   `json:"id"`
	Type        string   `json:"type"`
	Total       int      `json:"total"`
	Success     int      `json:"success"`
	Failed      int      `jons:"failed"`
	StartTime   int64    `json:"start_time"`
	EndTime     int64    `json:"end_time"`
	LogFileName string   `json:"log_file_name"`
	LogFile     *os.File `json:"-"`
	LogContent  string   `json:"log_content"`
}

func (t *Task) AppendLog(log string) {
	t.LogContent = log
	t.LogFile.WriteString(log + "\n")
}

// Persist save to database
func (t *Task) Persist() {
	t.EndTime = time.Now().Unix()

	// 存储到实时任务中
	clause := fmt.Sprintf(`insert into 
		realtime_task(id, type, total, success, failed, start_time, end_time) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		ON CONFLICT (id) DO
		update set total=$3, success=$4, failed=$5, end_time=$7`,
	)

	_, err := db.DB.Exec(clause, t.ID, t.Type, t.Total,
		t.Success, t.Failed, t.StartTime, t.EndTime)
	if err != nil {
		logger.Logger.Printf(logutil.Error, err.Error())
	}

	// 存储到任务日志中
	logClause := fmt.Sprintf(`insert into 
		task_log(id, type, content, file, ctime) 
		VALUES ($1, $2, $3, $4, $5) 
		ON CONFLICT (id) DO
		update set content=$3`,
	)

	_, err = db.DB.Exec(logClause, t.ID, t.Type, t.LogContent, t.LogFileName, t.EndTime)
	if err != nil {
		logger.Logger.Printf(logutil.Error, err.Error())
	}
	return
}

// Release release task
func (t *Task) Release() {
	t.LogFile.Close()

	TaskPoolMutex.Lock()
	TaskPool[t.ID] = nil
	TaskPoolMutex.Unlock()
	return
}

func Get(id string) *Task {
	TaskPoolMutex.Lock()
	t, ok := TaskPool[id]
	TaskPoolMutex.Unlock()
	if ok {
		return t
	}
	return nil
}

// New new task
func New(id string) *Task {
	task := new(Task)
	task.ID = id
	task.Type = "deploy"
	task.StartTime = time.Now().Unix()
	filename := filepath.Join(os.TempDir(), id+".log")
	task.LogFileName = filename

	filePtr, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil
	}
	task.LogFile = filePtr

	// assign
	TaskPool[id] = task
	return task
}

func init() {
	TaskPool = make(map[string]*Task, 0)
}
