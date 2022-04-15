package timeutil

import (
	"time"
)

const TimeLayout = "2006-01-02 15:04:05"

func TimeFormat(t time.Time) string {
	return t.Format(TimeLayout)
}
