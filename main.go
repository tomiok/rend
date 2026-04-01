package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
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
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	subcmd := strings.ToLower(os.Args[1])

	switch subcmd {
	case "add":
		runAdd(os.Args[2:])
	case "list":
		runList()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", subcmd)
		printUsage()
		os.Exit(1)
	}
}

func runAdd(args []string) {
	fs := flag.NewFlagSet("add", flag.ExitOnError)
	htmlPath := fs.String("html", "", "path to custom HTML widget")
	interval := fs.Duration("interval", 0, "refresh interval (e.g. 5s, 1m)")
	width := fs.Int("width", 0, "window width")
	height := fs.Int("height", 0, "window height")
	_ = fs.Parse(args)

	wv := buildWebView()

	if *htmlPath != "" {
		runCustomHTML(wv, *htmlPath, *interval, *width, *height)
	} else if fs.NArg() > 0 {
		runBuiltin(fs.Arg(0), wv)
	} else {
		fmt.Fprintln(os.Stderr, "rend add: specify a widget name or --html <file>")
		os.Exit(1)
	}

	wv.Run()
	wv.Destroy()
}

func printUsage() {
	fmt.Println(`Usage:
  rend add <widget>               launch a built-in widget
  rend add --html <file.html>     launch a custom widget
  rend list                       list available built-in widgets
 
Built-in widgets: clock, host`)
}

// runBuiltin lanza uno de los widgets integrados por nombre.
// Reemplaza al switch original en run(), pero ahora recibe el nombre
// limpio (sin el prefijo "add") gracias al parsing previo en runAdd.
func runBuiltin(name string, wv webview.WebView) {
	switch strings.ToLower(name) {
	case widgetClock:
		runClock(wv)
	case widgetHostInfo:
		runHostInfo(wv)
	default:
		_, _ = fmt.Fprintf(os.Stderr, "rend: unknown widget %q\n", name)
		_, _ = fmt.Fprintln(os.Stderr, "Run 'rend list' to see available widgets.")
		os.Exit(1)
	}
}

// runList imprime los widgets disponibles y sus descripciones.
// Se invoca con: rend list
func runList() {
	widgets := []struct {
		Name        string
		Description string
		Refresh     string
	}{
		{widgetClock, "Current time and date", "1s"},
		{widgetHostInfo, "CPU, memory, disk and system information", "1m"},
	}

	fmt.Println("Available built-in widgets:")
	fmt.Printf("  %-12s  %-44s  %s\n", "WIDGET", "DESCRIPTION", "REFRESH")
	fmt.Printf("  %-12s  %-44s  %s\n", "------", "-----------", "-------")
	for _, w := range widgets {
		fmt.Printf("  %-12s  %-44s  %s\n", w.Name, w.Description, w.Refresh)
	}
	fmt.Println()
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
					"updateClock(%d,%d,%d,'%s',%d,'%s',%d)",
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

func runCustomHTML(wv webview.WebView, path string, interval time.Duration, w int, h int) {
	data, err := os.ReadFile(path)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "rend: cannot read %s: %v\n", path, err)
		os.Exit(1)
	}

	title := filepath.Base(path)
	wv.SetTitle(title)

	width := 600
	height := 400

	if w > 0 {
		width = w
	}

	if h > 0 {
		height = h
	}

	wv.SetSize(width, height, webview.HintNone)
	wv.SetHtml(string(data))

	if interval > 0 {
		ticker := time.NewTicker(interval)
		go func() {
			defer ticker.Stop()
			for range ticker.C {
				wv.Dispatch(func() {
					wv.Eval("if(typeof onRendTick==='function') onRendTick()")
				})
			}
		}()
	}
}
