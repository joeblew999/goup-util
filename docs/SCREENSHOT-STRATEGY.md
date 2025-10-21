# Screenshot Strategy for goup-util

## Goal

Add screenshot capabilities to goup-util to capture running apps on:
1. **Desktop** - macOS, Windows, Linux
2. **Simulators** - iOS Simulator, Android Emulator

## Use Cases

- **Documentation**: Capture app screenshots for README/docs
- **Testing**: Visual regression testing
- **Demos**: Generate marketing materials
- **CI/CD**: Automated screenshot generation in pipelines

## Evaluation: go-vgo/robotgo

### Pros ✅

1. **Pure Go API** - Clean Go interface with screenshot methods:
   ```go
   robotgo.CaptureScreen(x, y, width, height)  // Capture region
   robotgo.CaptureImg()                        // Full screen as image.Image
   robotgo.SaveCapture("file.png")            // Direct save to file
   robotgo.SaveJpeg(img, "file.jpeg", quality)
   ```

2. **Multi-display support** - Works with multiple monitors:
   ```go
   num := robotgo.DisplaysNum()
   for i := 0; i < num; i++ {
       robotgo.DisplayID = i
       img, _ := robotgo.CaptureImg()
       robotgo.Save(img, fmt.Sprintf("screen_%d.png", i))
   }
   ```

3. **Cross-platform** - Supports macOS, Windows, Linux (X11)
4. **Active development** - Well maintained, good documentation
5. **Screenshot-only subset** - We only need capture features, not keyboard/mouse automation

### Cons ❌

1. **CGO dependency** - Requires:
   - GCC on all platforms
   - X11/Xtst libraries on Linux
   - libpng for bitmap operations
   - Platform-specific build requirements

2. **Cross-compilation challenges** - Cannot easily cross-compile from macOS to Windows/Linux due to CGO

3. **Binary size** - Adds ~2-5MB to binary (includes C libraries)

4. **Simulator support** - Does NOT capture iOS Simulator or Android Emulator directly
   - Would need platform-specific tools (see below)

## Alternative Approaches

### For Desktop Screenshots

#### Option A: Use robotgo (Recommended for desktop)
- Best cross-platform Go solution
- Works well when compiling natively on each platform
- Good for local development and native CI runners

#### Option B: Shell out to platform tools
```go
// macOS
screencapture -x -t png output.png

// Linux
import (DISPLAY=:0.0)
gnome-screenshot -f output.png
# or scrot, maim, etc.

// Windows
# PowerShell
[System.Windows.Forms.Screen]::PrimaryScreen.CaptureScreen()
```
**Pros**: No CGO, lighter weight
**Cons**: Platform-specific code, external dependencies

### For iOS Simulator

Use `xcrun simctl`:
```bash
# List running simulators
xcrun simctl list devices booted

# Take screenshot
xcrun simctl io booted screenshot output.png

# Or specific device
xcrun simctl io <device-uuid> screenshot output.png
```

### For Android Emulator/Device

Use `adb`:
```bash
# Find devices
adb devices

# Take screenshot
adb exec-out screencap -p > screenshot.png

# Or pull from device
adb shell screencap -p /sdcard/screenshot.png
adb pull /sdcard/screenshot.png
```

## Recommended Architecture

### Phase 1: Platform-Specific CLI (Simple & Reliable)

Create `goup-util screenshot` command that wraps platform tools:

```go
// cmd/screenshot.go
type ScreenshotCmd struct {
    Target   string  // "desktop", "ios-sim", "android-emu"
    Output   string  // Output file path
    Device   string  // Device ID (for simulators)
    Region   string  // x,y,w,h for partial capture
}
```

**Implementation**:
```bash
# Desktop screenshot (current window/full screen)
goup-util screenshot --target desktop --output screenshot.png

# iOS Simulator (auto-detect running simulator)
goup-util screenshot --target ios-sim --output ios-screenshot.png

# Android Emulator
goup-util screenshot --target android --output android-screenshot.png

# Specific device
goup-util screenshot --target ios-sim --device <uuid> --output screenshot.png
```

**Platform tools used**:
- **macOS desktop**: `screencapture` (built-in)
- **iOS Simulator**: `xcrun simctl` (built-in with Xcode)
- **Android**: `adb` (already managed by goup-util)
- **Windows**: PowerShell + .NET (future)
- **Linux**: `gnome-screenshot` or `scrot` (future)

### Phase 2: Optional robotgo Integration (Parallel with Phase 1)

Add robotgo as **optional dependency** via build tags for advanced features:
- Multi-monitor support
- Precise region capture
- Cross-platform consistency
- Window detection and activation
- Pixel color detection
- Find/match screen regions

**Build tag implementation**:

```go
//go:build screenshot
// +build screenshot

// pkg/screenshot/robotgo.go
package screenshot

import "github.com/go-vgo/robotgo"

type RobotgoCapturer struct{}

func (c *RobotgoCapturer) CaptureDesktop(output string) error {
    img, err := robotgo.CaptureImg()
    if err != nil {
        return err
    }
    return robotgo.Save(img, output)
}
```

**Without build tag** (default, no CGO):

```go
// pkg/screenshot/cli.go (no build tag)
package screenshot

import "os/exec"

type CLICapturer struct{}

func (c *CLICapturer) CaptureDesktop(output string) error {
    // Use platform CLI tools (screencapture, scrot, etc.)
    return exec.Command("screencapture", "-x", output).Run()
}
```

**Building**:
```bash
# Default: CLI tools only (no CGO)
go build .
goup-util screenshot --output test.png  # Uses screencapture/scrot/etc

# With robotgo: Advanced features
go build -tags screenshot .
goup-util screenshot --output test.png --all-displays  # Uses robotgo
```

**Benefits of this approach**:
- ✅ Default build has no CGO dependency (works everywhere)
- ✅ Power users can opt-in to robotgo features
- ✅ Both implementations exist in codebase simultaneously
- ✅ Runtime detection: If built with -tags screenshot, use robotgo; otherwise use CLI
- ✅ Documented in `goup-util screenshot --help`

**Reference**: See [docs/agents/robotgo.md](agents/robotgo.md) for full implementation guide

## Implementation Plan

### Week 1: CLI Screenshot Command (High Priority)

1. Add `cmd/screenshot.go` with Cobra command
2. Implement desktop screenshot (macOS only initially)
   - Use `exec.Command("screencapture", ...)`
   - Support `-x` (no sound), `-t png` (format)
   - Support region with `-R x,y,w,h`

3. Implement iOS Simulator screenshot
   - Auto-detect booted simulators with `xcrun simctl list`
   - Use `xcrun simctl io booted screenshot`
   - Handle multiple running simulators

4. Implement Android screenshot
   - Reuse existing ADB integration from installer
   - Use `adb exec-out screencap -p`
   - Handle multiple connected devices

### Week 2: Taskfile Integration

Add screenshot tasks to Taskfile.yml:
```yaml
screenshot:desktop:
  desc: Take desktop screenshot
  cmds:
    - "{{.GOUP}} screenshot --target desktop --output docs/screenshots/desktop-{{OS}}-{{ARCH}}.png"

screenshot:ios:
  desc: Screenshot iOS Simulator
  cmds:
    - "{{.GOUP}} screenshot --target ios-sim --output docs/screenshots/ios-sim.png"

screenshot:android:
  desc: Screenshot Android Emulator
  cmds:
    - "{{.GOUP}} screenshot --target android --output docs/screenshots/android-emu.png"

screenshots:all:
  desc: Generate all screenshots for docs
  deps:
    - screenshot:desktop
    - screenshot:ios
    - screenshot:android
```

### Week 3: Documentation & CI Integration

1. Add `docs/screenshots/` folder
2. Create screenshot workflow for CI:
   - GitHub Actions: Build app → Launch simulator → Screenshot → Upload artifacts
3. Update README with screenshots
4. Document usage in CLAUDE.md

## Decision: Start Simple, Add Complexity Later

**Recommendation**: Start with Phase 1 (CLI wrappers)
- ✅ No CGO dependency
- ✅ Uses proven platform tools
- ✅ Works in CI environments
- ✅ Easy to maintain
- ✅ Fast to implement

**Future**: Add robotgo as optional build tag
- For users who need advanced features
- For cross-platform desktop automation
- When CGO is acceptable

## Code Example

```go
// pkg/screenshot/screenshot.go
package screenshot

import (
    "fmt"
    "os/exec"
    "runtime"
)

type Capturer interface {
    Capture(output string) error
}

func NewCapturer(target string) (Capturer, error) {
    switch target {
    case "desktop":
        return newDesktopCapturer()
    case "ios-sim":
        return newIOSSimCapturer()
    case "android":
        return newAndroidCapturer()
    default:
        return nil, fmt.Errorf("unknown target: %s", target)
    }
}

// macOS desktop screenshot
type macOSCapturer struct{}

func (c *macOSCapturer) Capture(output string) error {
    cmd := exec.Command("screencapture", "-x", "-t", "png", output)
    return cmd.Run()
}

// iOS Simulator screenshot
type iosSimCapturer struct {
    deviceID string
}

func (c *iosSimCapturer) Capture(output string) error {
    device := c.deviceID
    if device == "" {
        device = "booted" // Use currently running simulator
    }
    cmd := exec.Command("xcrun", "simctl", "io", device, "screenshot", output)
    return cmd.Run()
}

// Android screenshot
type androidCapturer struct {
    deviceID string
}

func (c *androidCapturer) Capture(output string) error {
    args := []string{"exec-out", "screencap", "-p"}
    if c.deviceID != "" {
        args = append([]string{"-s", c.deviceID}, args...)
    }

    cmd := exec.Command("adb", args...)
    outFile, err := os.Create(output)
    if err != nil {
        return err
    }
    defer outFile.Close()

    cmd.Stdout = outFile
    return cmd.Run()
}
```

## Summary

**Start with**: Platform CLI wrappers (Week 1)
**Add later**: robotgo as optional (Phase 2)
**Focus on**: Desktop + iOS Simulator + Android Emulator
**Keep**: CGO-free main build, robotgo optional via build tags

This gives us:
1. Fast implementation (1-2 weeks)
2. Zero CGO dependencies by default
3. Works in all CI environments
4. Room to grow with robotgo later
