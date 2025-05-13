package time

import (
	"strconv"
	"strings"
	"time"
)

type WeekPeriod struct {
	Days            string `json:"days"`
	PeriodTimeStart int64  `json:"period_time_start"`
	PeriodTimeEnd   int64  `json:"period_time_end"`
}

// days取值: 0-6的数字
// days取值: Sun、Mon、Tue、Wed、Thu、Fri和Sat
// days取值：working-day
// days取值: off-day
// days取值: daily
func (o *WeekPeriod) isDaysActive() bool {
	dow := time.Now().Weekday()
	if o.Days == "daily" {
		return true
	}
	if o.Days == "off-day" {
		if dow == time.Saturday || dow == time.Sunday {
			return true
		}
		return false
	}
	if o.Days == "working-day" {
		if dow >= time.Monday && dow <= time.Friday {
			return true
		}
		return false
	}
	days := strings.Split(o.Days, ",")

	for _, day := range days {
		daysInt, err := strconv.Atoi(day)
		if err != nil {
			if strings.EqualFold(dow.String()[0:3], day) {
				return true
			}
		} else if int(dow) == daysInt {
			return true
		}
	}

	return false
}

func (o *WeekPeriod) IsActive() bool {
	if !o.isDaysActive() {
		return false
	}
	if o.PeriodTimeStart == 0 && o.PeriodTimeEnd == 0 {
		return true
	}

	now := time.Now()
	nowUnix := now.Unix()

	ztm := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, nil).Unix()
	if ztm+o.PeriodTimeStart < nowUnix && nowUnix < ztm+o.PeriodTimeEnd {
		return true
	}
	return false
}
