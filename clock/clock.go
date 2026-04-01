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
	Month   string
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
		Month:   translateMonth(month),
		Year:    year,
		Hour:    h,
		Minute:  m,
		Second:  s,
	}
}

func translateMonth(i time.Month) string {
	switch i {
	case 1:
		return "January"
	case 2:
		return "February"
	case 3:
		return "March"
	case 4:
		return "April"
	case 5:
		return "May"
	case 6:
		return "June"
	case 7:
		return "July"
	case 8:
		return "August"
	case 9:
		return "September"
	case 10:
		return "October"
	case 11:
		return "November"
	case 12:
		return "December"
	default:
		return "Unknown"
	}
}
