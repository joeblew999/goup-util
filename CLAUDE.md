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

## CRITICAL: Gio Version Compatibility

**VERSION MANAGEMENT IS CRITICAL** - Gio and gio-plugins version mismatches cause runtime panics!

### The Problem

- Using `gioui.org v0.8.0` with `gio-plugins v0.8.0` causes: `panic: Gio version not supported`
- The version tags don't guarantee compatibility - specific commit hashes are required
- See issue: https://github.com/gioui-plugins/gio-plugins/issues/104

### The Solution

**Always use these specific versions** (as of 2025-10-21):

```bash
# For projects using gio-plugins (webviewer, hyperlink, etc.)
go get gioui.org@1a17e9ea3725cf5bcb8bdd363e8c6b310669e936
go get github.com/gioui-plugins/gio-plugins@main
go mod tidy

# For projects using only Gio UI (no plugins)
go get gioui.org@1a17e9ea3725cf5bcb8bdd363e8c6b310669e936
go mod tidy
```

This gives you:
- `gioui.org v0.8.1-0.20250526181049-1a17e9ea3725` (commit after v0.8.0 tag)
- `github.com/gioui-plugins/gio-plugins v0.8.1-0.20250616220248-653221ccd770` (main branch)

### When Adding New Examples

1. **ALWAYS** use the version commands above after `go mod init`
2. **NEVER** use `@latest` tags - they are incompatible
3. **TEST** the app actually launches before committing
4. **UPDATE** this section if recommended versions change

### Version Management TODO

goup-util should eventually automate this:
- `goup-util init <project>` - Initialize with correct versions
- `goup-util doctor` - Check and fix version compatibility
- Auto-update go.mod when building examples

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

## Idempotency Guarantees

**ALL build operations are idempotent** - Safe to run multiple times, skips unnecessary work.

### Build Cache System

Located in `pkg/buildcache/`, tracks:
- **Source hashes** (SHA256 of .go, .mod, .sum files)
- **Output paths** and timestamps
- **Build success** status
- **Platform-specific** caching

### Smart Rebuild Detection

```bash
# First build - compiles everything
go run . build macos examples/hybrid-dashboard
# Building hybrid-dashboard for macos...
# ✓ Built successfully

# Second build - skips (no changes)
go run . build macos examples/hybrid-dashboard
# ✓ hybrid-dashboard for macos is up-to-date (use --force to rebuild)

# Force rebuild
go run . build --force macos examples/hybrid-dashboard
# Rebuilding: forced rebuild requested

# Check if rebuild needed (for CI/CD)
go run . build --check macos examples/hybrid-dashboard
echo $?  # 0=up-to-date, 1=needs rebuild
```

### What Triggers Rebuilds

✅ **Triggers rebuild:**
- Source code changes (.go files)
- Dependencies change (go.mod, go.sum)
- Assets change (.png, .jpg for icons)
- Output missing or corrupted
- `--force` flag

❌ **Skips rebuild:**
- No source changes
- Output exists and valid
- Previous build successful
- Same platform

### Build Flags

All build commands support:
- `--force` - Force rebuild even if up-to-date
- `--check` - Check if rebuild needed (exit code 0=no, 1=yes)

## Three-Tier Packaging System

goup-util provides **three distinct operations** for the app lifecycle:

### 1. Build - Compile & Create Basic Structures

```bash
go run . build <platform> <app>
```

**Purpose:** Fast iteration during development
- Compiles Go source to binaries
- Creates basic app structures (.app bundles, APKs)
- **Idempotent**: Uses build cache
- **Fast**: Skips unnecessary rebuilds
- Output: `<app>/.bin/`

### 2. Bundle - Create Signed App Bundles

```bash
go run . bundle <platform> <app> [--bundle-id ID] [--sign IDENTITY]
```

**Purpose:** Prepare for distribution
- Creates proper app bundles with metadata
- Generates Info.plist from templates
- **Code signing** (auto-detect or specified)
- Hardened runtime entitlements (macOS)
- **Pure Go**: No bash scripts
- Output: `<app>/.dist/<name>.app`

```bash
# Examples
go run . bundle macos examples/hybrid-dashboard
go run . bundle macos examples/hybrid-dashboard --bundle-id com.company.app
go run . bundle macos examples/hybrid-dashboard --sign "Developer ID Application: Name"
```

**Code Signing:**
- Auto-detects "Developer ID Application" certificates
- Falls back to "Apple Development" if needed
- Uses ad-hoc signature (`-`) if no certificates found
- Ad-hoc suitable for local testing, not public distribution

### 3. Package - Create Distribution Archives

```bash
go run . package <platform> <app>
```

**Purpose:** Final distribution packages
- Creates tar.gz (macOS/iOS) or zip (Windows) archives
- Copies APKs (Android)
- Ready for upload/distribution
- **Pure Go**: Uses pkg/packaging/archive.go
- Output: `<app>/.dist/<name>-<platform>.tar.gz`

### Complete Release Workflow

```bash
# 1. Build (idempotent)
go run . build macos examples/hybrid-dashboard

# 2. Create signed bundle
go run . bundle macos examples/hybrid-dashboard \
  --bundle-id com.company.myapp \
  --version 1.0.0

# 3. Test the bundle
open examples/hybrid-dashboard/.dist/hybrid-dashboard.app

# 4. Package for distribution
go run . package macos examples/hybrid-dashboard

# 5. Upload the archive
ls examples/hybrid-dashboard/.dist/*.tar.gz
```

See [docs/PACKAGING.md](docs/PACKAGING.md) for complete details.

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

- **robotgo** (`.src/robotgo/`) - Desktop automation and screenshots
  - Screenshot capabilities (CaptureScreen, CaptureImg, SaveCapture)
  - Multi-display support, keyboard/mouse automation
  - Platform-specific C code for macOS, Windows, Linux
  - See [docs/agents/robotgo.md](docs/agents/robotgo.md) for detailed guide
  - **Note**: Used optionally via build tags to avoid CGO in main build

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

## Screenshot Integration

**IMPLEMENTED**: goup-util now has built-in screenshot capabilities using robotgo!

### Quick Usage

```bash
# Show display information
task screenshot:info

# Take a screenshot
task screenshot

# Capture with delay (for menus/tooltips)
task screenshot:delay

# Capture all displays
task screenshot:all

# Or use directly
CGO_ENABLED=1 go run . screenshot output.png
```

### macOS Permission Setup

**CRITICAL**: On macOS 10.15+, grant Screen Recording permission:

1. System Settings → Privacy & Security → Screen Recording
2. Add Terminal.app (or your IDE)
3. Restart the terminal

The command will show a helpful error if permission is missing.

### Documentation Best Practices

### Screenshots for README

**IMPORTANT**: The README should have visual proof that this works!

Use the built-in screenshot command to capture running apps:

1. **Desktop Apps** (macOS):
```bash
# Build and launch the app
task run:hybrid

# Capture screenshot with delay (wait for app to load)
CGO_ENABLED=1 go run . screenshot --delay 2000 docs/screenshots/hybrid-dashboard-macos.png

# Or use the task
task docs:screenshots
```

2. **Mobile Simulators** (iOS/Android):
```bash
# iOS Simulator - use native screenshot
# Build and launch
task build:hybrid:ios
open -a Simulator
xcrun simctl install booted examples/hybrid-dashboard/.bin/hybrid-dashboard.app
xcrun simctl launch booted com.example.hybrid-dashboard

# Capture using simctl
xcrun simctl io booted screenshot docs/screenshots/hybrid-dashboard-ios.png

# Android Emulator - use adb
task build:hybrid:android
adb install examples/hybrid-dashboard/.bin/hybrid-dashboard.apk
adb shell am start -n com.example.hybrid/.MainActivity

# Capture using adb
adb exec-out screencap -p > docs/screenshots/hybrid-dashboard-android.png
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

## Taskfile Maintenance

**CRITICAL**: When adding new commands, examples, or features, **ALWAYS update Taskfile.yml**!

### Why This Matters

The Taskfile is the **primary developer interface**. Users run `task --list` to discover what they can do. If you add a feature but don't add a task for it, **nobody will know it exists**.

### When to Update Taskfile

**Add a new task whenever you:**
- ✅ Create a new example app (`examples/new-app/`)
- ✅ Add a new command to `cmd/`
- ✅ Add a new platform target
- ✅ Create a new workflow or common operation
- ✅ Add testing or deployment capabilities

### Task Naming Convention

Follow this pattern:
```yaml
# Format: <action>:<target>:<platform>
task build:hybrid:macos        # Build hybrid-dashboard for macOS
task build:hybrid:ios          # Build hybrid-dashboard for iOS
task run:hybrid                # Build and run hybrid-dashboard
task build:examples:android    # Build all examples for Android
```

**Categories:**
- `run:*` - Build and launch (for quick testing)
- `build:*` - Build only
- `install:*` - Install SDKs/dependencies
- `test:*` - Run tests
- `clean:*` - Clean up artifacts
- `workspace:*` - Workspace management
- `setup` - One-time setup tasks
- `demo` - Quick demonstrations

### Example: Adding a New Example App

When you create `examples/new-app/`:

```yaml
# Add these tasks to Taskfile.yml

vars:
  NEW_APP_EXAMPLE: examples/new-app

tasks:
  run:new-app:
    desc: Build and run new-app example (macOS)
    cmds:
      - "{{.GOUP}} build macos {{.NEW_APP_EXAMPLE}}"
      - open {{.NEW_APP_EXAMPLE}}/.bin/new-app.app

  build:new-app:macos:
    desc: Build new-app for macOS
    cmds:
      - "{{.GOUP}} build macos {{.NEW_APP_EXAMPLE}}"

  build:new-app:ios:
    desc: Build new-app for iOS
    cmds:
      - "{{.GOUP}} build ios {{.NEW_APP_EXAMPLE}}"

  build:new-app:android:
    desc: Build new-app for Android
    cmds:
      - "{{.GOUP}} build android {{.NEW_APP_EXAMPLE}}"
```

**Then update the composite tasks:**
```yaml
  build:examples:macos:
    desc: Build all examples for macOS
    cmds:
      - "{{.GOUP}} build macos {{.BASIC_EXAMPLE}}"
      - "{{.GOUP}} build macos {{.WEBVIEWER_EXAMPLE}}"
      - "{{.GOUP}} build macos {{.HYBRID_EXAMPLE}}"
      - "{{.GOUP}} build macos {{.NEW_APP_EXAMPLE}}"  # ADD THIS
```

### Testing Your Tasks

Before committing, **always test**:
```bash
# Verify task syntax
task --list

# Test the new task
task run:new-app

# Test composite tasks still work
task build:examples:macos
```

### Taskfile Anti-Patterns

**DON'T:**
- ❌ Add features without corresponding tasks
- ❌ Use inconsistent naming
- ❌ Forget to update composite tasks (build:examples:all, etc.)
- ❌ Hardcode paths (use vars instead)
- ❌ Create duplicate tasks

**DO:**
- ✅ Keep tasks simple and composable
- ✅ Use descriptive names
- ✅ Add helpful descriptions
- ✅ Test before committing
- ✅ Update README if adding major workflows

### Quick Reference

```bash
# See all tasks
task --list

# Run a task
task demo

# Run with verbose output
task -v demo

# See what a task will do (dry run)
task --dry demo
```

### Remember

**The Taskfile is the front door.** Keep it updated, or features will be invisible to users.

**Golden Rule**: If you can do it with `go run .`, there should be a task for it.

## Testing Taskfile Targets

**IMPORTANT**: The Taskfile is our PRIMARY testing mechanism right now. Always verify tasks work!

### Why Test Tasks?

Currently, goup-util has **limited unit tests**. The Taskfile tasks serve as:
- ✅ Integration tests (build → run workflows)
- ✅ Smoke tests (does it build at all?)
- ✅ Platform compatibility tests
- ✅ User workflow validation

**If a task is broken, users can't use the tool!**

### Test Before Committing

**Always run these before pushing**:

```bash
# 1. Verify all tasks are listed
task --list

# 2. Test info/config tasks (fast)
task config
task list:sdks
task workspace:list

# 3. Test icon generation (fast)
task icons:hybrid

# 4. Test at least one build (moderate)
task run:hybrid      # Builds + launches

# 5. Run Go tests (if they exist)
task test
```

### Full Task Test Suite

Create a test script to verify all tasks:

```bash
#!/bin/bash
# test-tasks.sh

echo "Testing Taskfile targets..."

# Info tasks
task config || echo "❌ config failed"
task list:sdks || echo "❌ list:sdks failed"
task workspace:list || echo "❌ workspace:list failed"

# Icon tasks
task icons:hybrid || echo "❌ icons:hybrid failed"

# Build tasks (one per platform to save time)
task build:hybrid:macos || echo "❌ build:hybrid:macos failed"

# Run task (check if app launches)
task run:hybrid || echo "❌ run:hybrid failed"

echo "✓ Task testing complete"
```

### Common Task Issues

**Problem**: Task fails with "invalid keys in command"
- **Cause**: YAML syntax error (usually unescaped colons in strings)
- **Fix**: Use single quotes or escape colons

**Problem**: Task fails with "command not found"
- **Cause**: Wrong path to binary or missing dependency
- **Fix**: Check {{.GOUP}} variable, verify binary exists

**Problem**: Task hangs
- **Cause**: Waiting for user input or long operation
- **Fix**: Add timeout or make operation non-interactive

### Task Testing Checklist

When modifying Taskfile.yml:

- [ ] Run `task --list` (verify syntax)
- [ ] Test the modified task
- [ ] Test any dependent tasks
- [ ] Verify task description is clear
- [ ] Check task works on clean environment
- [ ] Update this checklist if adding new categories

### Integration with CI/CD

**Future**: Add GitHub Actions workflow to test all tasks:

```yaml
# .github/workflows/test-tasks.yml
name: Test Taskfile
on: [push, pull_request]
jobs:
  test:
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v4
      - uses: arduino/setup-task@v1
      - run: task test
      - run: task build:hybrid:macos
      # etc.
```

### Remember

**Every task is a promise to users.** If `task run:hybrid` doesn't work, you've broken that promise.

Test your tasks. Keep them working. They're the user interface.
