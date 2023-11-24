package myutils

import (
	"time"
)

func ParseTime(startTime, endTime string) (start, end int64, err error) {
	if endTime == "" {
		endTime = time.Now().Format("2006-01-02 15:04:05")
	}
	s, err := time.Parse("2006-01-02 15:04:05", startTime)
	if err != nil {
		return 0, 0, err
	}
	e, err := time.Parse("2006-01-02 15:04:05", endTime)
	if err != nil {
		return 0, 0, err
	}
	return s.Unix(), e.Unix(), nil
}
