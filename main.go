package main

import (
	"encoding/json"
	"fmt"
	"os"
	"rend/clock"
	"rend/host"
	"strings"
	"time"

	webview "github.com/webview/webview_go"
)

const (
	widgetClock    = "clock"
	widgetHostInfo = "host"
)

func main() {
	args := os.Args[1:]
	var name string
	if args != nil && len(args) > 0 {
		name = args[0]
	}

	wv := buildWebView()
	run(strings.ToLower(name), wv)

	wv.Run()
	wv.Destroy()
}

func run(name string, wv webview.WebView) {
	switch name {
	case widgetClock:
		runClock(wv)
	case widgetHostInfo:
		runHostInfo(wv)
	default:
		runClock(wv)
	}
}

func buildWebView() webview.WebView {
	return webview.New(false)
}

func runClock(wv webview.WebView) {
	wv.SetTitle("clock")
	wv.SetSize(300, 150, webview.HintFixed)
	wv.SetHtml(clock.TodayHtml)
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		defer ticker.Stop()
		for range ticker.C {
			today := clock.Clock()
			wv.Dispatch(func() {
				wv.Eval(fmt.Sprintf(
					"updateClock(%d,%d,%d,'%s',%d,'%d',%d)",
					today.Hour, today.Minute, today.Second, today.Weekday, today.Day, today.Month, today.Year,
				))
			})
		}
	}()
}

func runHostInfo(wv webview.WebView) {
	wv.SetTitle("host information")
	wv.SetSize(900, 500, webview.HintFixed)
	wv.SetHtml(host.InfoHtml)

	ticker := time.NewTicker(time.Minute)
	go func() {
		time.Sleep(500 * time.Millisecond)

		info, _ := host.NewInformation()
		data, _ := json.Marshal(info)
		wv.Dispatch(func() {
			wv.Eval(fmt.Sprintf("updateHost(%s)", string(data)))
		})

		defer ticker.Stop()
		for range ticker.C {
			info, _ = host.NewInformation()
			data, _ = json.Marshal(info)
			wv.Dispatch(func() {
				wv.Eval(fmt.Sprintf("updateHost(%s)", data))
			})
		}
	}()
}
