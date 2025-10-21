package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
	"github.com/gioui-plugins/gio-plugins/plugin/gioplugins"
	"github.com/gioui-plugins/gio-plugins/webviewer/giowebview"
)

//go:embed web/*
var webContent embed.FS

// SystemStats represents system information exposed to JavaScript
type SystemStats struct {
	Platform    string  `json:"platform"`
	GoVersion   string  `json:"goVersion"`
	CPUUsage    float64 `json:"cpuUsage"`
	MemoryUsage float64 `json:"memoryUsage"`
	Uptime      int64   `json:"uptime"`
}

var (
	startTime = time.Now()
	th        = material.NewTheme()
)

func main() {
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	
	// Start embedded HTTP server
	serverURL := startWebServer()
	fmt.Printf("Web server started at %s\n", serverURL)

	// Launch Gio UI app
	go runApp(serverURL)
	app.Main()
}

func startWebServer() string {
	// Find available port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	// Serve embedded web content
	webFS, err := fs.Sub(webContent, "web")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(webFS)))
	
	// API endpoint: Get system stats (called from JavaScript)
	mux.HandleFunc("/api/stats", handleStats)
	
	// API endpoint: Say hello from Go
	mux.HandleFunc("/api/hello", handleHello)

	serverAddr := fmt.Sprintf("127.0.0.1:%d", port)
	go func() {
		log.Printf("HTTP server listening on http://%s\n", serverAddr)
		if err := http.ListenAndServe(serverAddr, mux); err != nil {
			log.Fatal(err)
		}
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	return fmt.Sprintf("http://%s", serverAddr)
}

func handleStats(w http.ResponseWriter, r *http.Request) {
	stats := SystemStats{
		Platform:    fmt.Sprintf("%s", os.Getenv("GOOS")),
		GoVersion:   "1.25.0",
		CPUUsage:    rand.Float64() * 100,      // Simulated
		MemoryUsage: 50 + rand.Float64()*40,    // Simulated
		Uptime:      time.Since(startTime).Milliseconds(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func handleHello(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"message": "Hello from Go! 🚀",
		"time":    time.Now().Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func runApp(serverURL string) {
	window := &app.Window{}
	window.Option(app.Title("Hybrid Dashboard - Gio + WebView"))
	window.Option(app.Size(1200, 800))

	var ops op.Ops
	webviewTag := new(int)

	for {
		evt := gioplugins.Hijack(window)

		switch evt := evt.(type) {
		case app.DestroyEvent:
			os.Exit(0)
			return

		case app.FrameEvent:
			gtx := app.NewContext(&ops, evt)

			// Layout: WebView fills entire window
			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					// Render WebView
					defer giowebview.WebViewOp{Tag: webviewTag}.Push(gtx.Ops).Pop(gtx.Ops)
					
					// Position and size
					giowebview.OffsetOp{Point: f32.Point{X: 0, Y: 0}}.Add(gtx.Ops)
					giowebview.RectOp{
						Size: f32.Point{
							X: float32(gtx.Constraints.Max.X),
							Y: float32(gtx.Constraints.Max.Y),
						},
					}.Add(gtx.Ops)

					// Navigate to embedded web server on first load
					for {
						ev, ok := gioplugins.Event(gtx, giowebview.Filter{Target: webviewTag})
						if !ok {
							break
						}
						
						switch ev.(type) {
						case giowebview.NavigationEvent:
							// Handle navigation events if needed
						}
					}

					// Initial navigation
					gioplugins.Execute(gtx, giowebview.NavigateCmd{
						View: webviewTag,
						URL:  serverURL,
					})

					return layout.Dimensions{Size: gtx.Constraints.Max}
				}),
			)

			evt.Frame(gtx.Ops)
		}
	}
}
