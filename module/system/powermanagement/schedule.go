package powermanagement

import (
	"encoding/json"
	"os/exec"
	"time"

	"github.com/google/uuid"
	"github.com/open-cmi/cmmns/module/dbkv"
	"github.com/open-cmi/gobase/essential/logger"
	"github.com/open-cmi/gobase/essential/ticker"
)

const PowerScheduleKey = "system-power-scheduled-tasks"

type Schedule struct {
	ID          string `json:"id"`
	Kind        string `json:"kind"`         // reboot, shutdown
	Frequency   string `json:"frequency"`    // once, daily, weekly, monthly
	TimeHM      string `json:"time_hm"`      // HH:mm
	OnceDate    string `json:"once_date"`    // YYYY-MM-DD
	WeeklyDays  []int  `json:"weekly_days"`  // 0-6 (Sun-Sat)
	MonthlyDays []int  `json:"monthly_days"` // 1-31
	Active      bool   `json:"active"`
}

func GetSchedules() ([]Schedule, error) {
	var schedules []Schedule
	val, ok := dbkv.Get(PowerScheduleKey)
	if !ok {
		return []Schedule{}, nil
	}
	err := json.Unmarshal([]byte(val), &schedules)
	return schedules, err
}

func SaveSchedules(schedules []Schedule) error {
	data, err := json.Marshal(schedules)
	if err != nil {
		return err
	}
	return dbkv.Set(PowerScheduleKey, string(data))
}

func AddSchedule(s *Schedule) error {
	schedules, _ := GetSchedules()
	s.ID = uuid.New().String()
	s.Active = true
	schedules = append(schedules, *s)
	return SaveSchedules(schedules)
}

func DeleteSchedule(id string) error {
	schedules, _ := GetSchedules()
	var newList []Schedule
	for _, s := range schedules {
		if s.ID != id {
			newList = append(newList, s)
		}
	}
	return SaveSchedules(newList)
}

func ExecuteAction(kind string) {
	logger.Infof("power management: executing %s", kind)
	if kind == "reboot" {
		exec.Command("reboot").Run()
	} else if kind == "shutdown" {
		exec.Command("shutdown", "-h", "now").Run()
	}
}

func CheckAndExecute() {
	now := time.Now()
	currentHM := now.Format("15:04")
	currentDate := now.Format("2006-01-02")
	currentWeekday := int(now.Weekday())
	currentDay := now.Day()

	schedules, err := GetSchedules()
	if err != nil {
		return
	}

	needsSave := false
	var activeSchedules []Schedule

	for _, s := range schedules {
		if !s.Active {
			activeSchedules = append(activeSchedules, s)
			continue
		}

		shouldRun := false
		if s.TimeHM == currentHM {
			switch s.Frequency {
			case "once":
				if s.OnceDate == currentDate {
					shouldRun = true
				}
			case "daily":
				shouldRun = true
			case "weekly":
				for _, d := range s.WeeklyDays {
					if d == currentWeekday {
						shouldRun = true
						break
					}
				}
			case "monthly":
				for _, d := range s.MonthlyDays {
					if d == currentDay {
						shouldRun = true
						break
					}
				}
			}
		}

		if shouldRun {
			// We run this in a goroutine to not block the checker loop
			go ExecuteAction(s.Kind)

			// If it's a "once" task, remove it after execution
			if s.Frequency == "once" {
				needsSave = true
				continue
			}
		}
		activeSchedules = append(activeSchedules, s)
	}

	if needsSave {
		SaveSchedules(activeSchedules)
	}
}

func init() {
	// Register a ticker to check every minute (0th second of every minute)
	ticker.Register("power_schedule_checker", "0 * * * * *", func(name string, data interface{}) {
		CheckAndExecute()
	}, nil)
}
