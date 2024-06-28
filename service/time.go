package service

import "time"

func StrToDate(str string) (time.Time, error) {
	loc, _ := time.LoadLocation("Local")
	return time.ParseInLocation("2006-01-02", str, loc)
}

func StrToTime(str string) (time.Time, error) {
	loc, _ := time.LoadLocation("Local")
	return time.ParseInLocation("2006-01-02T15:04:05", str, loc)
}
