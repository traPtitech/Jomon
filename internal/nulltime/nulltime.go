package nulltime

import (
	"bytes"
	"encoding/json"
	"time"
)

func ParseDate(str string) (NullTime, error) {
	loc, _ := time.LoadLocation("Local")
	t, err := time.ParseInLocation("2006-01-02", str, loc)
	if err != nil {
		return NullTime{}, err
	}
	return FromTime(&t), nil
}

func ParseTime(str string) (NullTime, error) {
	loc, _ := time.LoadLocation("Local")
	t, err := time.ParseInLocation("2006-01-02T15:04:05", str, loc)
	if err != nil {
		return NullTime{}, err
	}
	return FromTime(&t), nil
}

func FromTime(t *time.Time) NullTime {
	if t == nil {
		return NullTime{}
	}
	return NullTime{Time: *t, Valid: !t.IsZero()}
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
