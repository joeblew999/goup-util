---
title: "Quick Start"
date: 2025-12-21
draft: false
weight: 1
---

# Quick Start

Get up and running with goup-util in 5 minutes.

## Prerequisites

- [Go 1.24+](https://golang.org/)
- macOS, Linux, or Windows
- [Task](https://taskfile.dev/) (optional but recommended for running predefined workflows)

## Install goup-util

```bash
# Clone the repository
git clone https://github.com/joeblew999/goup-util
cd goup-util

# Build goup-util
go build .

# Verify it works
./goup-util --help
```

Or install system-wide:

```bash
go run . self setup
```

## Build Your First App

The `hybrid-dashboard` example is the best starting point -- it's a Gio UI app with an embedded webview.

### macOS

```bash
# Build the hybrid dashboard
goup-util build macos examples/hybrid-dashboard

# Open it
open examples/hybrid-dashboard/.bin/macos/hybrid-dashboard.app
```

Or use Task:
```bash
task run:hybrid
```

### Android

```bash
# Install Android SDK and NDK (one-time)
goup-util install android-sdk
goup-util install android-ndk

# Build APK
goup-util build android examples/hybrid-dashboard
```

### iOS (requires macOS + Xcode)

```bash
# Install Xcode from App Store first, then:
xcode-select --install

# Build for iOS Simulator
goup-util build ios-simulator examples/hybrid-dashboard

# Build for device
goup-util build ios examples/hybrid-dashboard
```

### Windows

```bash
# On a Windows machine:
goup-util build windows examples/hybrid-dashboard
```

## Gio Version Compatibility

**Important:** If you're creating a new project that uses gio-plugins (webviewer, hyperlink), you must pin specific versions. Mismatched versions cause runtime panics.

```bash
go get gioui.org@7bcb315ee174
go get github.com/gioui-plugins/gio-plugins@v0.9.1
go mod tidy
```

See [Platform Support](/users/platforms/#gio-version-compatibility) for details.

## Generate Icons

Generate platform-specific icons from a source image:

```bash
goup-util icons examples/hybrid-dashboard
```

This reads `icon-source.png` from the project directory and generates icons for all platforms (icns, ico, Android drawables).

Requirements:
- `icon-source.png` must exist in the project root
- Square PNG, 512x512 or larger recommended

## Common Commands

```bash
# Build and run in one step
goup-util run macos examples/hybrid-dashboard

# Force rebuild (ignore cache)
goup-util build --force macos examples/hybrid-dashboard

# Check if rebuild needed
goup-util build --check macos examples/hybrid-dashboard

# Build with deep linking schemes
goup-util build macos examples/hybrid-dashboard --schemes "myapp://"

# List available SDKs
goup-util list

# Show configuration
goup-util config

# Check goup-util installation health
goup-util self doctor
```

## Using Task

goup-util comes with a comprehensive [Taskfile](https://taskfile.dev/) for common workflows:

```bash
# See all available tasks
task --list

# Quick demo (build and run hybrid dashboard)
task dev:demo

# Build all examples for macOS
task build:examples:macos

# Start Hugo docs server
task hugo:start
```

## Zero-Compile Option

Want to ship a website as a desktop app without writing any Go code?

The [Webviewer Shell](/users/webviewer-shell/) is a pre-built binary that loads any URL from an `app.json` config file. Download, edit the URL, run. No Go, no SDKs, no compilation.

## Next Steps

- **[Platform Support](/users/platforms/)** -- Platform-specific details, requirements, and known limitations
- **[Packaging](/users/packaging/)** -- Create signed bundles and distribution archives
- **[Webviewer Shell](/users/webviewer-shell/)** -- Zero-compile option for shipping web apps
- **[Architecture](/architecture/)** -- How goup-util and the webview system work

## Troubleshooting

**Build fails with "SDK not found"**
- Run `goup-util install <sdk-name>` to install the required SDK
- Run `goup-util list` to see available SDKs

**Icons not generating**
- Ensure `icon-source.png` exists in your project directory
- Use a square PNG, 512x512 or larger

**macOS "can't be opened because Apple cannot check it"**
- Right-click the app, click Open, then click Open in the dialog
- Or run: `xattr -cr path/to/app.app`

**Gio version panic**
- Pin versions: `go get gioui.org@7bcb315ee174 && go get github.com/gioui-plugins/gio-plugins@v0.9.1 && go mod tidy`

**Getting help**
- `goup-util --help` for command reference
- File issues at [github.com/joeblew999/goup-util/issues](https://github.com/joeblew999/goup-util/issues)
