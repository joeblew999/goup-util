# Screenshot Integration

The `goup-util screenshot` command provides cross-platform screenshot capabilities for testing and documentation.

## Quick Start

```bash
# Take a screenshot (default mode - no CGO)
go run . screenshot output.png

# Or via Taskfile
task screenshot
```

## Two Modes

### 1. CLI Tools (Default - Recommended)

**No CGO required, lighter binaries**

- ✅ **macOS**: `screencapture` (built-in)
- ✅ **Linux**: `scrot` (install: `sudo apt install scrot`)
- ✅ **Windows**: PowerShell (built-in)

Build normally:
```bash
go build .
```

### 2. Robotgo (Optional - Advanced Features)

**Requires CGO, full multi-display support**

Features:
- Multi-display capture
- Display information
- Precise region capture
- Window detection

Build with:
```bash
go build -tags screenshot .
# or
task build:screenshot
```

## macOS Permission Setup

**IMPORTANT**: macOS 10.15+ requires Screen Recording permission.

### Grant Permission

1. **System Settings → Privacy & Security → Screen Recording**
2. Add your terminal app (Terminal.app, iTerm2, VSCode, etc.)
3. Restart the terminal app

### Check Permission

```bash
# Try a screenshot
go run . screenshot test.png

# If you see "could not create image from display"
# → You need to grant permission
```

## Usage Examples

### Basic Screenshots

```bash
# Full screen
go run . screenshot fullscreen.png

# With Taskfile
task screenshot
```

### Region Capture

```bash
# Capture specific region
go run . screenshot --x 100 --y 100 -w 800 -H 600 region.png

# With Taskfile
task screenshot:region
```

### Delayed Capture

Useful for capturing menus, tooltips, or hover states:

```bash
# 3 second delay
go run . screenshot --delay 3000 screenshot.png

# With Taskfile
task screenshot:delay
```

### Multi-Display (Robotgo Only)

```bash
# Build with robotgo first
task build:screenshot

# Capture all displays
./goup-util-screenshot screenshot --all --prefix display

# Get display info
./goup-util-screenshot screenshot --info
```

## Taskfile Commands

```bash
# Basic screenshot
task screenshot

# Region capture
task screenshot:region

# Delayed capture
task screenshot:delay

# Build with robotgo support
task build:screenshot

# Test multi-display features
task test:screenshot

# Capture for documentation
task docs:screenshots
```

## Documentation Workflow

For capturing app screenshots for README.md:

```bash
# 1. Build and launch the app
task run:hybrid

# 2. Wait for app to be ready

# 3. Capture screenshot with delay
go run . screenshot --delay 2000 docs/screenshots/app-macos.png

# 4. Update README.md with screenshot path
```

Or use the automated task:
```bash
task docs:screenshots
```

## Platform-Specific Notes

### macOS

- ✅ Built-in `screencapture` command
- ⚠️ Requires Screen Recording permission (macOS 10.15+)
- ✅ Region capture supported
- ✅ Silent mode (`-x` flag disables camera sound)

**Permission Check**:
```bash
# This will fail if permission not granted
screencapture test.png

# Success = no error, file created
ls -lh test.png
```

### Linux

- 📦 Install scrot: `sudo apt install scrot`
- ✅ Region capture supported
- ⚠️ Requires X11 (Wayland not supported via CLI)

**Install scrot**:
```bash
sudo apt install scrot

# Test
scrot test.png
```

### Windows

- ✅ Built-in PowerShell support
- ⚠️ Region capture not supported in CLI mode
- ℹ️ For advanced features, use robotgo build

## Troubleshooting

### macOS: "could not create image from display"

**Solution**: Grant Screen Recording permission

1. System Settings → Privacy & Security → Screen Recording
2. Add your terminal app
3. Restart terminal

### Linux: "scrot not found"

**Solution**: Install scrot
```bash
sudo apt install scrot
```

### Linux: "Can't open X display"

**Solution**: Set DISPLAY variable (usually `:0`)
```bash
export DISPLAY=:0
go run . screenshot test.png
```

### Windows: PowerShell execution errors

**Solution**: Check PowerShell execution policy
```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

## Building with Robotgo

### Requirements

**macOS**:
- Xcode Command Line Tools: `xcode-select --install`
- Screen Recording permission

**Linux**:
```bash
# GCC
sudo apt install gcc libc6-dev

# X11 dependencies
sudo apt install libx11-dev xorg-dev libxtst-dev libpng++-dev
```

**Windows**:
```bash
# Install MinGW-w64
winget install MartinStorsjo.LLVM-MinGW.UCRT

# Add to PATH
C:\mingw64\bin
```

### Build Command

```bash
# Enable CGO and screenshot tag
CGO_ENABLED=1 go build -tags screenshot -o goup-util-screenshot .

# Test
./goup-util-screenshot screenshot --info
```

### Verify Build

```bash
# CLI mode (default)
go run . screenshot --help
# Shows: "Current mode: CLI tools on darwin"

# Robotgo mode
./goup-util-screenshot screenshot --help
# Shows: "Current mode: robotgo (CGO) on darwin"
```

## Integration in CI/CD

For automated documentation screenshots:

```yaml
# .github/workflows/screenshots.yml
name: Capture Screenshots

on: [push]

jobs:
  screenshots:
    runs-on: macos-latest  # macOS has built-in screencapture
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v4

      - name: Build app
        run: task build:hybrid:macos

      - name: Launch app
        run: |
          open examples/hybrid-dashboard/.bin/hybrid-dashboard.app &
          sleep 3

      - name: Capture screenshot
        run: task screenshot docs/screenshots/hybrid-macos.png

      - name: Upload screenshots
        uses: actions/upload-artifact@v4
        with:
          name: screenshots
          path: docs/screenshots/
```

## API Reference

### Command Structure

```bash
goup-util screenshot [output-file] [flags]
```

### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--x` | - | int | 0 | X coordinate of region |
| `--y` | - | int | 0 | Y coordinate of region |
| `--width` | `-w` | int | 0 | Width of region |
| `--height` | `-H` | int | 0 | Height of region |
| `--all` | - | bool | false | Capture all displays (robotgo) |
| `--prefix` | - | string | "display" | Prefix for multi-display |
| `--delay` | `-d` | int | 0 | Delay before capture (ms) |
| `--quality` | `-q` | int | 90 | JPEG quality (1-100) |
| `--info` | - | bool | false | Show display info (robotgo) |

### Exit Codes

- `0`: Success
- `1`: Screenshot failed (permission, missing tool, etc.)

## Architecture

The screenshot package uses a two-tier approach:

```
cmd/screenshot.go
    ↓
pkg/screenshot/
    ├── screenshot.go    (interface)
    ├── cli.go          (default - platform CLI tools)
    ├── robotgo.go      (optional - build tag)
    └── robotgo_stub.go (stub when not built with tag)
```

### Build Tags

```go
//go:build screenshot
// +build screenshot
```

Files with this tag only compile when `-tags screenshot` is used.

## Future Enhancements

- [ ] Video recording support
- [ ] Window-specific capture (by PID/title)
- [ ] Automatic screenshot comparison (visual testing)
- [ ] Screenshot annotations
- [ ] GIF animation support

## Related Documentation

- [robotgo agent guide](agents/robotgo.md) - Deep dive into robotgo
- [WEBVIEW-ANALYSIS.md](WEBVIEW-ANALYSIS.md) - Screenshot webviews
- [Taskfile.yml](../Taskfile.yml) - All screenshot tasks

## References

- **robotgo**: https://github.com/go-vgo/robotgo
- **macOS screencapture**: `man screencapture`
- **Linux scrot**: `man scrot`
- **Build tags**: https://pkg.go.dev/cmd/go#hdr-Build_constraints
