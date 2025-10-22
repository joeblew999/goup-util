# UTM Integration Guide

## Overview

goup-util integrates with UTM to enable **automated Windows testing from macOS**. This allows you (and Claude) to trigger Windows builds, run tests, and retrieve artifacts without manually switching to the VM.

**Key capability**: Execute `goup-util` commands inside a Windows VM from the macOS host.

## Requirements

### macOS Host

- UTM installed ([https://mac.getutm.app](https://mac.getutm.app))
- `utmctl` available (ships with UTM at `/Applications/UTM.app/Contents/MacOS/utmctl`)
- Symlink or add to PATH:
  ```bash
  ln -s /Applications/UTM.app/Contents/MacOS/utmctl /opt/homebrew/bin/utmctl
  ```

### Windows VM

1. **QEMU Guest Agent** (required for `utmctl` commands)
   - Usually installed automatically with UTM Windows VMs
   - Verify: Check Task Manager → Services → "QEMU Guest Agent"
   
2. **Development environment** (install via bootstrap):
   ```powershell
   # Run in Windows VM
   iwr https://raw.githubusercontent.com/joeblew999/goup-util/main/scripts/windows-bootstrap.ps1 -UseBasicParsing | iex
   ```
   
   This installs:
   - Git
   - Go
   - Task (Taskfile runner)
   - goup-util
   
3. **Clone repo in VM**:
   ```powershell
   cd $HOME
   git clone https://github.com/joeblew999/goup-util
   cd goup-util
   ```

## Quick Start

### List VMs

```bash
# From macOS
goup-util utm list
```

Output:
```
UUID                                 Status   Name
12345678-1234-1234-1234-123456789012 stopped  Windows 11
```

### Start/Stop VM

```bash
# Start
goup-util utm start "Windows 11"

# Check status
goup-util utm status "Windows 11"

# Stop
goup-util utm stop "Windows 11"
```

### Execute Commands in VM

```bash
# Build for Windows
goup-util utm exec "Windows 11" -- build windows examples/hybrid-dashboard

# Bundle for Windows
goup-util utm exec "Windows 11" -- bundle windows examples/hybrid-dashboard

# Check config
goup-util utm exec "Windows 11" -- config
```

**Note**: Commands after `--` are automatically prefixed with `goup-util`.

### Run Taskfile Tasks in VM

```bash
# Run a task
goup-util utm task "Windows 11" build:hybrid:windows

# Run tests
goup-util utm task "Windows 11" test:all
```

### File Transfer

```bash
# Pull MSIX from VM to macOS
goup-util utm pull "Windows 11" \
  "C:\\Users\\User\\goup-util\\examples\\hybrid-dashboard\\.bin\\hybrid-dashboard.msix" \
  ./artifacts/hybrid-dashboard.msix

# Push file to VM
goup-util utm push "Windows 11" \
  ./local-file.txt \
  "C:\\Users\\User\\Desktop\\file.txt"
```

## Using Taskfile Tasks

goup-util includes convenient Taskfile tasks for common UTM operations:

```bash
# List VMs
task utm:list

# Start/stop Windows VM
task utm:start
task utm:stop
task utm:status

# Test Windows build in VM
task utm:test:build

# Test Windows bundle in VM
task utm:test:bundle

# Run a task in VM
task utm:test:task

# Pull MSIX artifact
task utm:pull:msix
```

## How It Works

### Architecture

```
┌─────────────────────────────────────────┐
│  macOS Host                              │
│                                          │
│  ┌────────────────────┐                 │
│  │  goup-util utm     │                 │
│  │  (CLI wrapper)     │                 │
│  └─────────┬──────────┘                 │
│            │                             │
│            ▼                             │
│  ┌────────────────────┐                 │
│  │  utmctl            │                 │
│  │  (UTM CLI tool)    │                 │
│  └─────────┬──────────┘                 │
│            │                             │
└────────────┼─────────────────────────────┘
             │
             │ QEMU Guest Agent Protocol
             │
┌────────────▼─────────────────────────────┐
│  Windows 11 VM (UTM)                     │
│                                          │
│  ┌────────────────────┐                 │
│  │  QEMU Guest Agent  │                 │
│  └─────────┬──────────┘                 │
│            │                             │
│            ▼                             │
│  ┌────────────────────┐                 │
│  │  goup-util         │                 │
│  │  (Windows binary)  │                 │
│  └────────────────────┘                 │
│                                          │
│  Files, stdout, stderr, exit codes      │
│  returned to host                       │
└──────────────────────────────────────────┘
```

### Key Points

1. **No custom server needed** - Uses UTM's built-in guest agent
2. **Secure** - QEMU guest agent handles authentication
3. **Bidirectional file transfer** - Pull artifacts back to macOS
4. **Exit codes preserved** - Can detect build failures
5. **stdout/stderr captured** - See build output on macOS

## Testing Workflow

### Manual Testing

```bash
# 1. Start VM
task utm:start

# 2. Wait for boot (check status)
task utm:status

# 3. Run Windows build
task utm:test:build

# 4. Check if MSIX was created
goup-util utm exec "Windows 11" -- ls examples/hybrid-dashboard/.bin/

# 5. Pull MSIX to macOS
task utm:pull:msix

# 6. Verify artifact
ls -lh ./artifacts/hybrid-dashboard.msix

# 7. Stop VM
task utm:stop
```

### Automated Testing (Future)

Once this workflow is stable, we can automate it:

```bash
# Proposed: Run full Windows test suite
task test:windows:utm

# This would:
# 1. Start VM
# 2. Run all Windows builds
# 3. Run all Windows tests
# 4. Pull artifacts
# 5. Stop VM
# 6. Verify artifacts on macOS
```

## Troubleshooting

### QEMU Guest Agent Not Running

**Error**: `utmctl exec` fails with "guest agent not responding"

**Solution**:
1. Open Windows VM
2. Open Services (Win+R → `services.msc`)
3. Find "QEMU Guest Agent"
4. Right-click → Start
5. Set to "Automatic" startup

### Command Not Found in VM

**Error**: `goup-util: command not found`

**Solution**:
```bash
# Verify goup-util is installed in VM
goup-util utm exec "Windows 11" -- cmd /c where goup-util

# If not found, reinstall via bootstrap
# (Run this INSIDE the VM, not via utm exec)
```

### File Transfer Fails

**Error**: `failed to create output file`

**Solution**:
```bash
# Create artifacts directory first
mkdir -p ./artifacts

# Try again
task utm:pull:msix
```

### VM Name Has Spaces

Use quotes around VM name:

```bash
goup-util utm exec "Windows 11" -- build windows examples/hybrid-dashboard
```

## VM Path Conventions

**Windows paths in VM**:
- Home: `C:\Users\<Username>`
- goup-util repo: `C:\Users\<Username>\goup-util`
- Build outputs: `C:\Users\<Username>\goup-util\examples\<app>\.bin\`

**Windows path escaping**:
```bash
# Use double backslashes in shell
goup-util utm pull "Windows 11" "C:\\Users\\User\\file.txt" ./local/

# Or single backslash in quotes
goup-util utm pull "Windows 11" 'C:\Users\User\file.txt' ./local/
```

## Integration with CI/CD

### GitHub Actions (Future)

While GitHub Actions can't run UTM VMs directly, you could:

1. **Self-hosted runner** on your Mac with UTM
2. **Workflow** that uses `goup-util utm` commands
3. **Artifact upload** from pulled MSIX files

```yaml
# .github/workflows/test-windows-utm.yml (example)
name: Test Windows in UTM
on: [push]
jobs:
  test:
    runs-on: self-hosted  # Your Mac
    steps:
      - uses: actions/checkout@v4
      - name: Start Windows VM
        run: task utm:start
      - name: Test Windows build
        run: task utm:test:build
      - name: Pull MSIX
        run: task utm:pull:msix
      - name: Upload artifact
        uses: actions/upload-artifact@v3
        with:
          name: hybrid-dashboard-windows
          path: ./artifacts/hybrid-dashboard.msix
      - name: Stop VM
        run: task utm:stop
```

## Best Practices

1. **Start VM before testing** - `task utm:start` and wait for boot
2. **Use tasks for common operations** - More concise than full commands
3. **Pull artifacts after build** - Verify on macOS host
4. **Stop VM when done** - Save resources
5. **Keep VM updated** - Periodically run `goup-util self upgrade` in VM

## Next Steps

1. Set up Windows 11 VM in UTM
2. Install QEMU Guest Agent (usually automatic)
3. Run bootstrap script in VM
4. Clone repo in VM
5. Test with `task utm:list`
6. Try `task utm:test:build`

## Reference

- **UTM**: https://mac.getutm.app
- **utmctl docs**: https://docs.getutm.app/scripting/scripting/
- **QEMU Guest Agent**: https://wiki.qemu.org/Features/GuestAgent
