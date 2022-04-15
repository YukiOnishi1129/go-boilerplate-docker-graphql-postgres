package timeutil

import "time"

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

const TimeLayout = "2006-01-02 15:04:05"

func TimeFormat(t time.Time) string {
	return t.In(jst).Format(TimeLayout)
}
