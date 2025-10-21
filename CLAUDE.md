# Claude Instructions for goup-util

## Project Overview

**goup-util** is a specialized build tool for creating **cross-platform hybrid applications** using Go and Gio UI.

### The Real Mission

Enable developers to build **one codebase** that runs everywhere:
- 🖥️ Desktop: macOS, Windows, Linux
- 📱 Mobile: iOS, Android  
- 🌐 Web: Browser (via WASM)

**Key capability**: Hybrid apps mixing **native Gio UI** (for shell/controls) with **native webviews** (for rich content).

### Why This Matters

Traditional cross-platform tools require multiple languages (Swift, Kotlin, JavaScript). **goup-util enables pure Go development** for hybrid apps by:
1. Managing platform SDKs (Android SDK, Xcode tools)
2. Building platform-specific binaries from Go source
3. Handling native integrations (webviews, icons, packaging)
4. Supporting the full app lifecycle (build → package → release)

### Key Principles
- **Pure Go development**: One language for all platforms
- **Hybrid architecture**: Native UI + webview content
- **Idempotent operations**: All operations are safe to run multiple times
- **DRY (Don't Repeat Yourself)**: Centralized path management and shared utilities
- **Developer-focused**: Clean CLI interface with clear commands
- **True cross-platform**: Web, desktop, and mobile from single codebase

## Project Structure

```
goup-util/
├── cmd/                    # CLI commands (Cobra-based)
│   ├── root.go            # Root command
│   ├── build.go           # Build Gio apps for platforms
│   ├── install.go         # Install SDKs
│   ├── self.go            # Build/update goup-util itself
│   ├── icons.go           # Generate platform icons
│   ├── package.go         # Package apps for distribution
│   ├── workspace.go       # Manage Go workspaces
│   └── ...
├── pkg/                   # Shared packages
│   ├── config/           # Configuration management
│   ├── installer/        # SDK installation logic
│   ├── icons/            # Icon generation
│   ├── workspace/        # Go workspace utilities
│   └── ...
├── examples/             # Example Gio applications
│   ├── gio-basic/                # Simple Gio UI demo
│   ├── gio-plugin-hyperlink/     # Hyperlink plugin demo
│   └── gio-plugin-webviewer/     # Multi-tab webview browser (THE KEY EXAMPLE)
├── docs/                 # End-user documentation
│   ├── agents/           # AI assistant collaboration guides
│   └── WEBVIEW-ANALYSIS.md  # Cross-platform webview deep dive
├── .src/                 # Dependency source code (gitignored)
│   └── gio-plugins/      # gio-plugins source for reference
└── main.go              # Entry point

```

## The Hybrid App Vision

**goup-util exists to make this possible**:

```
┌─────────────────────────────────────┐
│     Your App (Pure Go)              │
│                                     │
│  ┌─────────────────────────────┐   │
│  │  Gio UI (Native Controls)   │   │
│  │  - Tabs, buttons, layout    │   │
│  │  - Native performance       │   │
│  └─────────────────────────────┘   │
│                                     │
│  ┌─────────────────────────────┐   │
│  │  Native WebView             │   │
│  │  - Rich web content         │   │
│  │  - HTML/CSS/JavaScript      │   │
│  │  - Platform webview engine  │   │
│  └─────────────────────────────┘   │
│                                     │
│  ↕ Go ↔ JavaScript Bridge          │
└─────────────────────────────────────┘

Built once → Runs on all platforms
```

## Development Workflow

### Building the Tool

```bash
# Build goup-util itself
go run . self build

# Build for all platforms
go run . self build --all

# Run tests
go test ./...

# Run integration tests
go test -v integration_test.go
```

### Building Hybrid Apps

```bash
# The webviewer example is THE reference implementation
go run . build macos examples/gio-plugin-webviewer
go run . build windows examples/gio-plugin-webviewer
go run . build android examples/gio-plugin-webviewer
go run . build ios examples/gio-plugin-webviewer

# Install required SDKs
go run . install android-sdk
go run . install android-ndk

# Generate platform icons
go run . icons examples/gio-plugin-webviewer
```

## Key Commands to Understand

- `build <platform> <app>` - Build Gio apps for different platforms (macos, windows, android, ios)
- `install <sdk>` - Install platform SDKs (Android SDK, NDK, etc.)
- `self build` - Build goup-util binaries for distribution
- `icons <app>` - Generate platform-specific icons from source images
- `package <app>` - Package built apps for distribution
- `workspace` - Manage Go workspace files for multi-module projects
- `gitignore` - Manage .gitignore templates for Gio projects

## Important Files

- `cmd/*.go` - All CLI command implementations
- `pkg/config/` - Config file handling and directory management
- `pkg/installer/` - SDK installation logic
- `examples/gio-plugin-webviewer/main.go` - **THE KEY EXAMPLE** - Multi-tab browser showing full webview capabilities
- `go.mod` - Dependencies (cobra, progressbar, icns, gio-plugins, etc.)
- `.gitignore` - Build binaries are excluded (goup-util*)

## Common Tasks

### Understanding Webview Integration

**This is the core use case**. Study these files:

1. **Local example**: `examples/gio-plugin-webviewer/main.go`
2. **Plugin source**: `.src/gio-plugins/webviewer/`
3. **Demo app**: `.src/gio-plugins/webviewer/demo/demo.go`
4. **Analysis**: `docs/WEBVIEW-ANALYSIS.md`
5. **Agent guide**: `docs/agents/gio-plugins.md`

### Adding a New Command

1. Create `cmd/newcommand.go`
2. Use Cobra patterns from existing commands
3. Add to root command in `cmd/root.go`
4. Test with `go run . newcommand`

### Modifying Build Logic

- See `cmd/build.go` for platform-specific build commands
- Build logic uses idempotent patterns
- Platform detection and SDK path resolution in `pkg/`

### Working with Icons

- Icon generation in `pkg/icons/`
- Supports PNG source → platform formats (icns, ico, Android drawables)
- Test with example projects

## Dependencies

### Core Dependencies
- `github.com/spf13/cobra` - CLI framework
- `github.com/schollz/progressbar/v3` - Progress display
- `github.com/JackMordaunt/icns` - macOS icon generation

### Gio Ecosystem (THE IMPORTANT ONES)
- `gioui.org` - Core Gio UI framework
- `github.com/gioui-plugins/gio-plugins` - Native plugins
  - **webviewer** - Native webview integration (WKWebView, WebView2, etc.)
  - **hyperlink** - Open URLs in system browser
  - **auth** - OAuth flows
  - **explorer** - File system access

### Platform Tools
- Android SDK, NDK - For Android builds
- Xcode - For iOS/macOS builds  
- WebView2 - For Windows (Edge-based webview)

## Testing Guidelines

- Test commands using `go run .` before building
- **Use the webviewer example for integration testing**
- Verify idempotency (running twice should produce same result)
- Test on target platforms when modifying build logic
- Check webview behavior on each platform (they differ!)

## CI/CD

- GitHub Actions in `.github/workflows/`
- `build.yml` - Main CI pipeline using `go run . self build`
- Builds binaries for all platforms
- Artifacts uploaded as releases

## Future Plans (See TODO.md)

- **UTM integration** - Automated Windows VM testing from macOS
- **Winget** - Windows package management for dependencies
- **Automated testing infrastructure** - Test builds on all platforms
- **JavaScript ↔ Go bridge patterns** - Better hybrid app communication
- **Production templates** - Ready-to-use hybrid app starters

## Source Code References (.src/)

The `.src/` folder contains cloned source code of key dependencies for easy reference. This folder is gitignored and local-only.

### Available Sources

- **gio-plugins** (`.src/gio-plugins/`) - Gio UI native plugins
  - WebViewer implementation and examples
  - Platform-specific code for macOS, iOS, Android, Windows, Linux
  - See [docs/agents/gio-plugins.md](docs/agents/gio-plugins.md) for detailed guide

### Usage

When working with dependencies:

```bash
# Clone a new dependency for reference
git clone --depth 1 https://github.com/org/repo .src/repo

# Search for implementations
grep -r "pattern" .src/gio-plugins/

# View platform-specific code
ls .src/gio-plugins/webviewer/webview/webview_*.go

# Read the webview demo (our example is based on this)
cat .src/gio-plugins/webviewer/demo/demo.go
```

### Agent Collaboration

For AI assistants working on this project:

1. **Read source before asking** - Check `.src/` for dependency behavior
2. **Update agent docs** - Add guides in `docs/agents/` when learning new patterns
3. **See agent guides** - Read `docs/agents/README.md` for collaboration patterns

The agent documentation helps multiple AI assistants work effectively on the codebase by providing context about dependencies, patterns, and architecture.

## Tips for Claude

1. **The webviewer example is CRITICAL** - This shows the real use case
2. **Always test with `go run .`** - Don't build binaries during development
3. **Maintain idempotency** - Operations should be safe to run multiple times
4. **Follow existing patterns** - Look at similar commands for consistency
5. **Update docs/** - Keep end-user docs in sync with code changes
6. **Check .gitignore** - Don't commit build binaries (goup-util*)
7. **Use examples/** - Test changes with the example Gio projects
8. **Cross-platform awareness** - Code runs on macOS, Linux, Windows, Android, iOS
9. **Hybrid apps are the goal** - Native UI + webview content in pure Go
10. **Read .src/ dependencies** - Source code is available locally

## Common Debugging

```bash
# Check configuration
go run . config

# List available SDKs
go run . list

# Verbose output (add -v flag if available)
go run . build macos examples/gio-basic -v

# Check Go workspace
go run . workspace list

# Test webviewer on desktop (fastest iteration)
go run . build macos examples/gio-plugin-webviewer
open examples/gio-plugin-webviewer/.bin/macos/gio-plugin-webviewer.app
```

## Code Style

- Follow standard Go conventions
- Use `cobra` for CLI structure
- Error handling with clear messages
- Progress bars for long operations
- Idempotent file operations
- Platform-specific code in separate files (`*_darwin.go`, `*_android.go`, etc.)

## The Big Picture

**goup-util is a developer tool for building a specific class of apps:**

**Cross-platform hybrid applications where:**
- Shell/controls are native Gio UI (Go)
- Content can be web (via native webviews)
- Everything is written in Go
- Deploys to web, desktop, and mobile from one codebase

This is about **enabling pure Go development** for the kind of apps that traditionally require Swift + Kotlin + JavaScript. The webview integration is what makes hybrid apps possible while keeping native performance.

## Documentation Best Practices

### Screenshots for README

**IMPORTANT**: The README should have visual proof that this works!

Use the Playwright MCP to capture screenshots of running apps:

1. **Desktop Apps** (macOS):
```bash
# Launch the app
open examples/gio-plugin-webviewer/.bin/gio-plugin-webviewer.app

# Use Playwright MCP to capture
mcp__playwright__browser_take_screenshot
# Save to: docs/screenshots/webviewer-macos.png
```

2. **Mobile Simulators** (iOS/Android):
```bash
# iOS Simulator
open -a Simulator
xcrun simctl install booted examples/gio-plugin-webviewer/.bin/gio-plugin-webviewer.app
xcrun simctl launch booted com.example.gio-plugin-webviewer

# Capture with screenshot tool
# Save to: docs/screenshots/webviewer-ios.png

# Android Emulator
adb install examples/gio-plugin-webviewer/.bin/gio-plugin-webviewer.apk
adb shell am start -n com.example.webviewer/.MainActivity

# Capture
adb exec-out screencap -p > docs/screenshots/webviewer-android.png
```

3. **README Structure**:
```markdown
# goup-util

![macOS Screenshot](docs/screenshots/webviewer-macos.png)
![iOS Screenshot](docs/screenshots/webviewer-ios.png)
![Android Screenshot](docs/screenshots/webviewer-android.png)

Build cross-platform hybrid apps in pure Go...
```

### Example Apps Should Be Complete

**Current Issue**: webviewer example only loads external URLs (https://google.com)

**Better**: Include embedded web server with example content

```go
// examples/hybrid-app-complete/
// 
// main.go - Gio UI + embedded HTTP server
package main

import (
    "embed"
    "net/http"
    
    "gioui.org/app"
    "github.com/gioui-plugins/gio-plugins/webviewer"
)

//go:embed web/*
var webContent embed.FS

func main() {
    // Start embedded web server
    go func() {
        http.Handle("/", http.FileServer(http.FS(webContent)))
        http.ListenAndServe("localhost:8080", nil)
    }()
    
    // Launch Gio app with webview pointing to localhost
    // ...
}

// web/index.html - HTML/CSS/JS content
// web/app.js - JavaScript that calls Go functions
// web/styles.css - Styling
```

This shows:
- ✅ Embedded web content (no external dependencies)
- ✅ Go HTTP server (real backend)
- ✅ Go ↔ JavaScript bridge (function calls)
- ✅ Complete, working example
- ✅ Can be used as template for real apps

### When Creating Documentation

1. **Always include screenshots** - Visual proof is powerful
2. **Show it running** - Not just code, but actual output
3. **Complete examples** - Should work out of the box
4. **Link screenshots in README** - First thing people see
5. **Update CLAUDE.md** when adding new patterns

This helps AI assistants maintain documentation quality.
