package service

import (
	"bytes"
	"encoding/json"
	"time"
)

func StrToDate(str string) (time.Time, error) {
	loc, _ := time.LoadLocation("Local")
	return time.ParseInLocation("2006-01-02", str, loc)
}

func StrToTime(str string) (time.Time, error) {
	loc, _ := time.LoadLocation("Local")
	return time.ParseInLocation("2006-01-02T15:04:05", str, loc)
}

func NullTimeToTime(t *time.Time) NullTime {
	if t == nil {
		return NullTime{}
	}
	return NullTime{Time:*t,Valid: *t!=time.Time{}}
}

type NullTime struct {
	Time  time.Time
	Valid bool
}

var jsonNull = []byte("null")

func (t NullTime) MarshalJSON() ([]byte, error) {
	if t.Valid {
		return t.Time.MarshalJSON()
	}
	return jsonNull, nil
}

func (t *NullTime) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, jsonNull) {
		*t = NullTime{}
		return nil
	}
	err := json.Unmarshal(b, &t.Time)
	t.Valid = err == nil
	return err
}
