package time

import "time"

type AbsoluteTimeObject struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	RefCnt         int    `json:"refcnt"`
	Active         bool   `json:"active"`
	TimestampStart int64  `json:"timestamp_start"`
	TimestampEnd   int64  `json:"timestamp_end"`
}

func (o *AbsoluteTimeObject) IsActive() bool {
	nowUnix := time.Now().Unix()

	if nowUnix < o.TimestampStart || nowUnix > o.TimestampEnd {
		return false
	}
	return true
}

type PeriodTimeObject struct {
	Name            string       `json:"name"`
	Description     string       `json:"description"`
	RefCnt          int          `json:"refcnt"`
	Active          bool         `json:"active"`
	WeekPeriods     []WeekPeriod `json:"week_periods"`
	TimeRangeEnable bool         `json:"time_range_enable"`
	TimeRangeStart  int64        `json:"time_range_start"`
	TimeRangeEnd    int64        `json:"time_range_end"`
}

func (o *PeriodTimeObject) IsActive() bool {
	nowUnix := time.Now().Unix()

	if o.TimeRangeEnable {
		if nowUnix < o.TimeRangeStart || nowUnix > o.TimeRangeEnd {
			return false
		}
	}
	var isActive bool = false
	for _, period := range o.WeekPeriods {
		if period.IsActive() {
			isActive = true
			break
		}
	}
	return isActive
}
