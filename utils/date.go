package utils

import "time"

func Date(format string, t ...time.Time) string {
	date := time.Now()
	if len(t) > 0 {
		date = t[0]
	}
	return date.Format(format)
}

func Time() int64 {
	return time.Now().Unix()
}
