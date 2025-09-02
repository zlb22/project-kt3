package constant

import "time"

const (
	TimeFormatDef    = "2006-01-02 15:04:05"
	TimeFormatDate   = "2006-01-02"
	TimeFormatWithMS = "2006-01-02 15:04:05.000" //其他类型
	TimeFormatTime   = "15:04:05"                //其他类型

	SecondDuration = time.Second        // SecondDuration 1秒
	HourDuration   = time.Hour          // HourDuration 1小时
	DayDuration    = time.Hour * 24     // DayDuration 1天
	WeekDuration   = time.Hour * 24 * 7 // WeekDuration 1周
)
