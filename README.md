# rend

A desktop widget system for Linux and macOS. Each widget is a native window rendering HTML powered by Go.

No Electron. No Node. No framework. Just a Go binary, a webview, and whatever HTML you want on your desktop.

---

## Installation

```sh
curl -sSL https://raw.githubusercontent.com/tomiok/rend/main/install.sh | sh
```

**Linux** requires the following system library, included by default on Ubuntu 22.04+ and Pop!_OS 22.04+:

```
libwebkit2gtk-4.0
```

**macOS** has no additional dependencies.

---

## Building from source

```sh
git clone https://github.com/tomiok/rend
cd rend
go build -o rend .
```

Linux also requires the development headers to compile:

```sh
sudo apt install libgtk-3-dev libwebkit2gtk-4.0-dev
```

---

## How it works

`rend` launches widgets as native windows using the system's built-in web renderer — WebKitGTK on Linux, WKWebView on macOS. Each widget runs as an independent process with its own window. Data flows from Go to the widget via JavaScript evaluation on a ticker interval.

---

## Usage

```
rend add <widget>            launch a built-in widget
rend add --html file.html    launch a custom widget from an HTML file
```

---

## Built-in widgets

| Widget | Description | Refresh |
|--------|-------------|---------|
| `clock` | Current time and date | 1 second |
| `host` | CPU, memory, disk and system information | 1 minute |

### clock

```sh
rend add clock
```

Displays the current time and date in a small floating window.

### host

```sh
rend add host
```

Displays a system overview: CPU model and usage, memory, disk and general system information including hostname, OS, uptime and kernel version.

---

## Custom widgets

Any HTML file can be a widget. Data is injected from Go via `window.eval()` on a configurable interval.

```sh
rend add --html ~/.config/rend/widgets/mywidget.html
```

---

## Built with

- [webview/webview_go](https://github.com/webview/webview_go) — Go bindings for the webview library
- [shirou/gopsutil](https://github.com/shirou/gopsutil) — system metrics

---

## Philosophy

One binary. No config files required. No background services. Each widget is a process — if it misbehaves, kill it. The HTML is yours to edit.

Inspired by the suckless philosophy: software that does one thing, whose source you can read in an afternoon.

---

## License

MIT