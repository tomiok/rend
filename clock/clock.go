package clock

import (
	"time"
)

var TodayHtml = `<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<style>
  * { margin: 0; padding: 0; box-sizing: border-box; }
  
  body {
    background: #1a1a2e;
    color: #e0e0e0;
    font-family: monospace;
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100vh;
    user-select: none;
  }

  .clock {
    text-align: center;
  }

  .time {
    font-size: 2.5rem;
    font-weight: bold;
    color: #a78bfa;
    letter-spacing: 0.1em;
  }

  .date {
    font-size: 1rem;
    color: #888;
    margin-top: 4px;
  }
</style>
</head>
<body>
  <div class="clock">
    <div class="time" id="time">00:00:00</div>
    <div class="date" id="date">Loading...</div>
  </div>

<script>
  function updateClock(h, m, s, weekday, day, month, year) {
    const pad = n => String(n).padStart(2, '0');
    document.getElementById('time').textContent = 
      pad(h) + ':' + pad(m) + ':' + pad(s);
    document.getElementById('date').textContent = 
      weekday + ' ' + day + ' de ' + month + ' ' + year;
  }
</script>
</body>
</html>`

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
