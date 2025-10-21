# goup-util

**Build cross-platform hybrid applications entirely in Go.**

Native UI + Web Content = Pure Go Development 🚀

![Status](https://img.shields.io/badge/status-alpha-orange)
![Go Version](https://img.shields.io/badge/go-1.25%2B-blue)
![Platforms](https://img.shields.io/badge/platforms-macOS%20%7C%20iOS%20%7C%20Android%20%7C%20Windows%20%7C%20Linux-green)

---

## What is goup-util?

A specialized build tool that enables you to create **hybrid applications** using:
- **Gio UI** for native controls and navigation
- **Native webviews** for rich web content  
- **Pure Go** for everything - no Swift, Kotlin, or JavaScript required

**Write once in Go → Deploy everywhere**

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
│  │  - HTML/CSS/JavaScript      │   │
│  │  - Platform webview engine  │   │
│  │  - WKWebView, WebView2      │   │
│  └─────────────────────────────┘   │
│                                     │
│  ↕ Go ↔ JavaScript Bridge          │
└─────────────────────────────────────┘
```

---

## Platform Support

| Platform | Build | Status | Notes |
|----------|-------|--------|-------|
| **macOS** | ✅ | Tested | Native .app bundles with WKWebView |
| **iOS** | ✅ | Tested | App Store ready, WKWebView integration |
| **Android** | ✅ | Tested | APK generation, Chromium WebView |
| **Windows** | ⚠️ | Untested | WebView2 support (cross-compile issue) |
| **Linux** | ⚠️ | Untested | WebKitGTK support (cross-compile issue) |
| **Web** | 🔜 | Planned | WASM deployment |

**All tested platforms work successfully!** ✨

---

## Quick Start

### Prerequisites

- **Go 1.25+**
- **macOS**: Xcode Command Line Tools
- **Android**: Auto-installed via `goup-util install ndk-bundle`

### Install

```bash

## Using Taskfile (Recommended)

We provide a [Taskfile](https://taskfile.dev) for common operations:

```bash
# Install Task first (if you don't have it)
brew install go-task/tap/go-task

# See all available tasks
task --list

# Quick demo - builds and runs hybrid-dashboard
task demo

# Build and run specific examples
task run:hybrid        # Hybrid dashboard with embedded server
task run:webviewer     # Multi-tab browser
task run:basic         # Simple Gio app

# Build for different platforms
task build:hybrid:macos
task build:hybrid:ios
task build:hybrid:android
task build:hybrid:all   # All platforms

# SDK management
task install:ndk        # Install Android NDK
task list:sdks          # Show available SDKs

# Development
task workspace:init     # Initialize Go workspace
task clean             # Clean build artifacts
task test              # Run tests

# Full setup from scratch
task setup             # Initialize workspace + install NDKs
```

**Quick start with Task:**
```bash
git clone https://github.com/joeblew999/goup-util.git
cd goup-util
task setup    # One command to set everything up
task demo     # See hybrid-dashboard in action!
```

git clone https://github.com/joeblew999/goup-util.git
cd goup-util
```

### Build Your First Hybrid App

```bash
# Set up workspace
go work init
go work use . examples/gio-plugin-webviewer

# Build for macOS
go run . build macos examples/gio-plugin-webviewer

# Build for iOS
go run . build ios examples/gio-plugin-webviewer

# Build for Android (installs NDK if needed)
go run . install ndk-bundle
go run . build android examples/gio-plugin-webviewer

# Launch the app!
open examples/gio-plugin-webviewer/.bin/gio-plugin-webviewer.app
```

**That's it!** You just built a multi-tab browser app in pure Go that runs on macOS, iOS, and Android.

---

## What Can You Build?

### Hybrid Apps with Embedded Web Content
- **Dashboards** - Native shell + web charts/graphs
- **Content Apps** - Native navigation + HTML articles
- **Dev Tools** - Native IDE + web inspector
- **Documentation** - Native app + rendered markdown

### Pure Native Apps
- **Productivity Tools** - All-native Gio UI
- **Utilities** - System integration apps
- **Games** - 2D/3D with Gio rendering

### Progressive Web Apps (PWA)
- **Web-first apps** packaged as native
- **Offline-capable** with service workers
- **App store distribution** of web apps

---

## Examples

### Basic Gio App
```bash
go run . build macos examples/gio-basic
```
Simple pure-Gio application showing native UI capabilities.

### Webviewer (Multi-tab Browser)
```bash
go run . build macos examples/gio-plugin-webviewer
```
**The key example** - demonstrates:
- ✅ Multiple webviews (tabs)
- ✅ URL navigation
- ✅ JavaScript execution
- ✅ Storage access (cookies, localStorage)
- ✅ Native UI + web content integration

### Hyperlink Integration
```bash
go run . build macos examples/gio-plugin-hyperlink
```
Shows how to open system browser from your app.

---

## Features

### 🎨 Automatic Icon Generation
```bash
# Generate platform-specific icons from one source
go run . icons macos examples/my-app
go run . icons android examples/my-app
go run . icons ios examples/my-app
```

### 📦 SDK Management
```bash
# Auto-installs and caches SDKs
go run . install ndk-bundle        # Android NDK
go run . install android-sdk        # Android SDK
go run . list                       # Show available SDKs
```

### 🔧 Workspace Integration
```bash
# Manage multi-module projects
go run . workspace list
go run . ensure-workspace examples/my-app
```

### 🚀 Self-Building
```bash
# Build goup-util itself
go run . self build
```

---

## Architecture

**Idempotent**: All operations are safe to run multiple times  
**DRY**: Centralized path management, no duplication  
**Clean**: Service layer ready for future API use  
**Caching**: SDKs downloaded once, reused forever

---

## Documentation

- **[IMPROVEMENTS.md](docs/IMPROVEMENTS.md)** - Roadmap and future enhancements
- **[WEBVIEW-ANALYSIS.md](docs/WEBVIEW-ANALYSIS.md)** - Deep dive into cross-platform webview support
- **[TODO.md](TODO.md)** - Current tasks and priorities
- **[CLAUDE.md](CLAUDE.md)** - AI assistant guide (for development)
- **[docs/agents/](docs/agents/)** - Dependency guides for AI collaboration
- **[docs/platforms.md](docs/platforms.md)** - Platform-specific build information

---

## Project Status

**Current Phase**: Proof of Concept → Production Ready

**What Works**:
- ✅ Builds succeed on macOS, iOS, Android
- ✅ Webviewer hybrid apps work on all tested platforms
- ✅ Icon generation for all platforms
- ✅ SDK caching and management
- ✅ Multi-module workspace support

**What's Next** (see [IMPROVEMENTS.md](docs/IMPROVEMENTS.md)):
1. **Better UX** - Progress bars, error messages, feedback
2. **Performance** - Incremental builds, parallel operations
3. **Webview Excellence** - Go ↔ JS bridge, TypeScript defs, DevTools
4. **Testing** - Automated testing, deployment helpers
5. **Windows/Linux** - Fix cross-compilation issues

---

## Why goup-util?

### vs Electron/Tauri
- ✅ **Much smaller binaries** (native webviews, not embedded browser)
- ✅ **Better performance** (no Node.js/Chromium overhead)
- ✅ **Pure Go** (one language, one ecosystem)
- ✅ **Mobile support** (iOS + Android, not just desktop)

### vs Flutter
- ✅ **Pure Go** (no Dart required)
- ✅ **Native webviews** (leverage existing web content)
- ✅ **Simpler** (no custom rendering engine)

### vs Native (SwiftUI/Jetpack Compose)
- ✅ **Cross-platform** (write once, deploy everywhere)
- ✅ **One language** (Go for all platforms)
- ✅ **Hybrid capable** (mix native + web seamlessly)

---

## Contributing

We're in active development! See [TODO.md](TODO.md) for current priorities.

**Quick wins needed**:
- Better progress feedback during builds
- Error messages with suggestions
- Complete hybrid app example with embedded server
- Screenshots and visual documentation
- Windows/Linux cross-compilation fixes

---

## License

[Check LICENSE file]

---

## Credits

Built on top of:
- **[Gio UI](https://gioui.org)** - Pure Go immediate-mode UI
- **[gio-plugins](https://github.com/gioui-plugins/gio-plugins)** - Native webview integration
- **[Cobra](https://github.com/spf13/cobra)** - CLI framework

---

## Vision

**Make Go the best choice for cross-platform hybrid application development.**

No Swift. No Kotlin. No JavaScript required.*

Just Go. Everywhere. 🚀

<sub>* JavaScript optional for web content in webviews</sub>
