package service

import "time"

func (s *Services) StrToDate(str string) (time.Time, error) {
	loc, _ := time.LoadLocation("Local")
	return time.ParseInLocation("2006-01-02", str, loc)
}

func (s *Services) StrToTime(str string) (time.Time, error) {
	loc, _ := time.LoadLocation("Local")
	return time.ParseInLocation("2006-01-02T15:04:05", str, loc)
}
