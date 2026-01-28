# goup-util Documentation

Build cross-platform hybrid apps in pure Go.

## Quick Navigation

### 📚 Getting Started
- [Quick Start Guide](quickstart.md) - Get up and running in 5 minutes
- [Platform Support](platforms.md) - Platform-specific features and requirements

### 🔧 Core Features
- [Packaging Guide](PACKAGING.md) - Create distribution-ready packages
- [CI/CD Integration](cicd.md) - Automated build pipelines

### 🏗️ Architecture
- [Webview Analysis](WEBVIEW-ANALYSIS.md) - Deep dive into hybrid app architecture
- [Improvements Roadmap](IMPROVEMENTS.md) - Architectural overview and future plans

### 🤖 AI Collaboration
- [AI Assistant Guide](agents/README.md) - Collaboration patterns for AI assistants
- [Gio Plugins Reference](agents/gio-plugins.md) - gio-plugins deep dive
- [robotgo Reference](agents/robotgo.md) - Screenshot system reference

## What is goup-util?

A build tool for creating **cross-platform hybrid applications** using Go and Gio UI.

**One codebase → Runs everywhere:**
- 🖥️ Desktop: macOS, Windows, Linux
- 📱 Mobile: iOS, Android
- 🌐 Web: Browser (WASM)
- 🔀 Hybrid: Native Gio UI + native webviews

## Key Capabilities

✅ **Pure Go Development** - One language for all platforms
✅ **Hybrid Architecture** - Native UI + webview content
✅ **SDK Management** - Automated install and caching
✅ **Asset Generation** - Icons for all platforms
✅ **Idempotent Builds** - Safe to run multiple times
✅ **Screenshot Capture** - App Store screenshot generation

## Quick Commands

```bash
# Build for macOS
go run . build macos examples/hybrid-dashboard

# Install Android SDK
go run . install android-sdk

# Generate icons
go run . icons examples/hybrid-dashboard

# Capture screenshots
task screenshot-hybrid
```

## Examples

Working demonstrations in `examples/`:

- **hybrid-dashboard** - Recommended starting point (Gio + webview)
- **gio-basic** - Simple Gio UI demo
- **gio-plugin-hyperlink** - URL handling
- **gio-plugin-webviewer** - Multi-tab browser

## Getting Help

- [TODO.md](../TODO.md) - Known issues and roadmap
- [CLAUDE.md](../CLAUDE.md) - Development guidelines
- [GitHub Issues](https://github.com/joeblew99/goup-util/issues) - Report bugs

---

**Start here:** [Quick Start Guide](quickstart.md)
