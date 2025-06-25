# examples

gio is at version 0.8.0

gio-plugins need a lower version of gio: 0.6.0

gio-plugins themselves are at 0.8.0: https://github.com/gioui-plugins/gio-plugins/releases/tag/v0.8.0

so the go.mod for any gio app using gio plugins is:

```go

module ...

go 1.24.4

require (
	gioui.org v0.8.0
	github.com/gioui-plugins/gio-plugins v0.8.0
)

```

