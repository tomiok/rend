package clock

import (
	_ "embed"
	"time"
)

//go:embed clock.html
var TodayHtml string

type Today struct {
	Weekday time.Weekday
	Day     int
	Month   int
	Year    int
	Hour    int
	Minute  int
	Second  int
}

func Clock() Today {
	now := time.Now()
	weekday := now.Weekday()
	day := now.Day()
	month := now.Month()
	year := now.Year()
	h, m, s := now.Clock()

	return Today{
		Weekday: weekday,
		Day:     day,
		Month:   int(month),
		Year:    year,
		Hour:    h,
		Minute:  m,
		Second:  s,
	}
}
