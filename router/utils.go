package router

import "time"

func StrToDate(str string) (time.Time, error) {
	loc, _ := time.LoadLocation("Local")
	return time.ParseInLocation("2006-01-02", str, loc)
}
