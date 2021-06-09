package common

import (
	"time"
)
const (
	TokenDateFormat = time.RFC3339Nano //"2006-01-02T15:04:05.012345678-08:00"	
	ModelDateFormat string = time.RFC3339 //"2006-01-02T15:04:05-08:00"
	CutOffDateFormat string = "2006-01-02"	
	ExactDateFormat string = "2006-01-02 15:04:05.012345678"
	RoundedDateFormat string = "2006-01-02"
)

func ShowDate(t time.Time) string {
	return t.Format(ExactDateFormat)
}

func IsBefore(t1, t2 time.Time) bool {
	return t2.Sub(t1) > 0
}

func IsFutureDate(t time.Time) bool {
	now := time.Now()
	return IsBefore(t, now)
}

// is it a correct cut off data?
func IsCutOffDate(v interface{}) bool {
	var t time.Time
	var err error
	
	// convert to time stamp
	switch v.(type) {
	case string:
		t, err = time.Parse(CutOffDateFormat, v.(string))
		if err != nil {
			return false
		}
		
	case time.Time:
		t = v.(time.Time)
	}

	// 1. in the future
	if IsFutureDate(t) {
		return false
	}
	
	// 2. 1st day of the month
	d := t.Day()
	h := t.Hour()
	m := t.Minute()
	s := t.Second()
	n := t.Nanosecond()
	if d == 1 && h == 0 && m == 0 && s == 0 && n == 0 {
		return true
	}
	
	return false
}

// closest 1st day of the month in the future of time buing
func NextCutOffDate() (cutoff time.Time, err error) {
	now := time.Now()
	y := now.Year()
    mon := now.Month()
    d := 1
    h := 0
    m := 0
    s := 0
	n := 0

	if mon == 12 { // dec -> next year
		mon = 1
		y++
	} else {
		mon++
	}

	nextCutOffDate := time.Date(y, mon, d, h, m, s, n, time.UTC)	
	return time.Parse(RoundedDateFormat, nextCutOffDate.Format(RoundedDateFormat))
}
