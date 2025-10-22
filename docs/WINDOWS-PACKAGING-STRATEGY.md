# Windows Packaging Strategy

## Current State Analysis

### Existing MSIX Code

goup-util currently has **two separate MSIX commands**:

1. **`msix-manifest`** (`cmd/msix_manifest.go`)
   - Generates `AppxManifest.xml` from template
   - Takes YAML data or CLI flags
   - **Pure Go** - works on any platform
   - **Template-based** - similar to our new macOS approach
   - **Status**: ✅ Already follows our new pattern!

2. **`msix-pack`** (`cmd/msix_pack.go`)
   - Creates MSIX package from directory
   - **Windows-only** - checks `runtime.GOOS`
   - Requires external `msix` toolkit command
   - **Problem**: Can't test from macOS

### The Problem

```
┌─────────────────────────────────────────┐
│  macOS Development Machine              │
│                                         │
│  ✅ Can test: macOS, iOS, Android      │
│  ❌ Cannot test: Windows MSIX           │
│                                         │
│  Options:                               │
│  1. Trust it works (risky)             │
│  2. VM with manual testing (slow)      │
│  3. Automated VM testing (complex)     │
└─────────────────────────────────────────┘
```

---

## Proposed Refactoring

### Phase 1: Align with New Packaging System (Now)

**Goal**: Make MSIX follow the same pattern as macOS (build → bundle → package)

#### Current Flow (Fragmented)
```bash
# Build Windows binary
goup-util build windows examples/hybrid-dashboard

# Generate manifest (separate command)
goup-util msix-manifest \
  --output AppxManifest.xml \
  --name MyApp \
  --publisher "CN=Company"

# Pack MSIX (separate command, Windows-only)
goup-util msix-pack \
  --directory ./build \
  --package myapp.msix
```

#### Proposed Flow (Unified)
```bash
# 1. Build (idempotent, works now)
goup-util build windows examples/hybrid-dashboard

# 2. Bundle (NEW - integrate msix-manifest logic)
goup-util bundle windows examples/hybrid-dashboard \
  --bundle-id MyApp \
  --publisher "CN=Company" \
  --version 1.0.0

# 3. Package (integrate msix-pack, Windows-only for now)
goup-util package windows examples/hybrid-dashboard
```

### Implementation Plan

#### 1. Create `pkg/packaging/windows.go`

```go
package packaging

import (
    _ "embed"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "runtime"
    "text/template"
)

//go:embed templates/windows-appxmanifest.xml.tmpl
var windowsManifestTemplate string

// WindowsBundleConfig contains configuration for creating a Windows MSIX bundle
type WindowsBundleConfig struct {
    // App identity
    Name                   string // App name
    Publisher              string // Publisher (e.g., "CN=MyCompany")
    PublisherDisplayName   string // Display name for publisher
    DisplayName            string // Display name
    Description            string // App description
    Version                string // Version (e.g., "1.0.0.0")

    // Paths
    BinaryPath  string // Path to the .exe
    OutputDir   string // Where to create the MSIX
    AssetsDir   string // Path to logo assets (optional)

    // Signing (for future)
    SigningCertificate string // Path to .pfx file
    CertificatePassword string // Password for certificate
}

// CreateWindowsBundle creates a Windows MSIX bundle
func CreateWindowsBundle(config WindowsBundleConfig) error {
    // 1. Create bundle directory structure
    // 2. Copy binary
    // 3. Generate AppxManifest.xml from template
    // 4. Copy/generate assets
    // 5. Call msix pack (if on Windows)
    // 6. Handle signing (future)

    return nil
}
```

#### 2. Move Template to `pkg/packaging/templates/`

Create `pkg/packaging/templates/windows-appxmanifest.xml.tmpl`:
- Move the template from `cmd/msix_manifest.go`
- Update to use `WindowsBundleConfig` struct
- Embed with `//go:embed`

#### 3. Update `cmd/bundle.go`

Add Windows case to bundle command:
```go
case "windows":
    return bundleWindows(proj, bundleID, version, publisher, outputDir)
```

#### 4. Deprecate Old Commands

Mark `msix-manifest` and `msix-pack` as deprecated:
- Add deprecation warnings
- Point to new `bundle` command
- Keep for backward compatibility (like bash script)

---

## Testing Strategy

### The Challenge

**We can't test Windows MSIX creation from macOS** because:
1. `msix pack` command only exists on Windows
2. MSIX file format is Windows-specific
3. Code signing requires Windows certificates

### Solution Tiers

#### Tier 1: ✅ What We Can Test Now (macOS)

**Template Generation**:
```go
// Test that AppxManifest.xml generates correctly
func TestWindowsManifestGeneration(t *testing.T) {
    config := WindowsBundleConfig{
        Name: "TestApp",
        Publisher: "CN=TestPublisher",
        // ...
    }

    manifest, err := generateWindowsManifest(config)
    // Assert XML is valid
    // Assert all fields populated correctly
}
```

**Bundle Structure**:
```go
// Test that bundle directory is created correctly
func TestWindowsBundleStructure(t *testing.T) {
    // Create bundle (will fail at msix pack step, that's OK)
    // Verify directory structure
    // Verify AppxManifest.xml exists
    // Verify binary copied
}
```

**Platform Detection**:
```go
// Test that Windows-only steps are skipped on macOS
func TestWindowsBundleGracefulSkip(t *testing.T) {
    if runtime.GOOS != "windows" {
        // Should skip msix pack gracefully
        // Should still create directory structure
    }
}
```

#### Tier 2: 🔄 Manual VM Testing (Short-term)

**Setup**:
1. Create Windows 11 VM in UTM
2. Install: Go, Git, MSIX toolkit
3. Clone goup-util
4. Build and test manually

**Test Checklist**:
```bash
# Inside Windows VM
go build .
.\goup-util.exe bundle windows examples\hybrid-dashboard
.\goup-util.exe package windows examples\hybrid-dashboard

# Verify MSIX created
dir examples\hybrid-dashboard\.dist\*.msix

# Test installation
Add-AppxPackage examples\hybrid-dashboard\.dist\hybrid-dashboard.msix
```

**Frequency**: Before releases, when Windows code changes

#### Tier 3: 🚀 Automated VM Testing (Future - Q4 2025)

**The UTM Automation Vision** (from TODO.md):

```
┌──────────────────────────────────────────────┐
│  macOS Host                                  │
│  ┌────────────────────────────────────┐     │
│  │  UTM (Virtual Machine)             │     │
│  │  ┌──────────────────────────────┐  │     │
│  │  │  Windows 11                  │  │     │
│  │  │                              │  │     │
│  │  │  Provisioned via:            │  │     │
│  │  │  - Git (via winget)          │  │     │
│  │  │  - Go (via winget)           │  │     │
│  │  │  - MSIX toolkit (via winget) │  │     │
│  │  │                              │  │     │
│  │  │  Tests Run:                  │  │     │
│  │  │  → Build Windows apps        │  │     │
│  │  │  → Package MSIX              │  │     │
│  │  │  → Install and verify        │  │     │
│  │  └──────────────────────────────┘  │     │
│  └────────────────────────────────────┘     │
└──────────────────────────────────────────────┘
```

**Implementation Steps**:

1. **VM Management in Go**:
```go
// pkg/vm/utm.go
type UTMManager struct {
    VMPath string
}

func (m *UTMManager) Start() error { ... }
func (m *UTMManager) Stop() error { ... }
func (m *UTMManager) Execute(cmd string) (string, error) { ... }
```

2. **WinRM Communication**:
```go
// pkg/vm/winrm.go - Communicate with Windows VM
// Use: github.com/masterzen/winrm
```

3. **Test Runner**:
```go
// cmd/test_windows.go
// Starts VM, runs tests, collects results
```

4. **CI Integration**:
```yaml
# .github/workflows/test-windows.yml
- name: Test Windows Packaging
  run: |
    task test:windows:vm
```

---

## Communication Layer Design (Future)

### The Problem

**Host (macOS) needs to:**
1. Start/stop Windows VM
2. Copy files to VM
3. Execute commands in VM
4. Retrieve results from VM

### Option 1: WinRM (Recommended)

**Pros**:
- Native Windows remote management
- Well-established protocol
- Go library available: `github.com/masterzen/winrm`

**Cons**:
- Requires WinRM enabled in VM
- Network configuration needed

**Implementation**:
```go
import "github.com/masterzen/winrm"

client, err := winrm.NewClient(&winrm.Endpoint{
    Host: "localhost",
    Port: 5985,
}, "Administrator", "password")

stdout, stderr, exitCode, err := client.Run(
    "goup-util.exe bundle windows examples\\hybrid-dashboard",
)
```

### Option 2: SSH

**Pros**:
- Universal protocol
- Good Go libraries

**Cons**:
- Need to install OpenSSH Server on Windows
- Less "native" than WinRM

### Option 3: Shared Folders + File Polling

**Pros**:
- Simple, no network config
- Works with UTM shared folders

**Cons**:
- Hacky, polling-based
- Slow, unreliable

**Implementation**:
```go
// 1. Host writes command to shared folder
os.WriteFile("/Volumes/Shared/command.txt", []byte("bundle windows ..."), 0644)

// 2. Windows agent polls and executes
// (Windows Go daemon)

// 3. Host polls for result
result, _ := os.ReadFile("/Volumes/Shared/result.txt")
```

**Recommendation**: Start with Option 3 for MVP, migrate to WinRM for production.

---

## Migration Path

### Now (Phase 1) ✅ Ready to Implement

**Goals**:
- Refactor MSIX into `pkg/packaging/windows.go`
- Integrate with `bundle` command
- Template-based AppxManifest.xml
- Deprecate old commands
- **Testing**: Go unit tests only (template generation, structure)

**Work Required**:
1. Create `pkg/packaging/windows.go` (2-3 hours)
2. Move template to `templates/` (30 min)
3. Update `cmd/bundle.go` (1 hour)
4. Write tests (2 hours)
5. Deprecate old commands (30 min)
6. Update docs (1 hour)

**Total**: ~7 hours of work

**Risk**: Low - all Go code, testable on macOS

### Later (Phase 2) - Manual VM Testing

**Goals**:
- Verify MSIX actually works
- Test installation on Windows
- Catch Windows-specific issues

**Work Required**:
1. Create Windows 11 UTM VM (2 hours)
2. Provision with winget packages (1 hour)
3. Write test checklist (1 hour)
4. Execute tests manually (1 hour per release)

**Total**: 5 hours setup + 1 hour per release

**Risk**: Medium - manual process, can be forgotten

### Future (Phase 3) - Automated Testing

**Goals**:
- Fully automated Windows testing
- CI/CD integration
- Fast feedback on PRs

**Work Required**:
1. VM management code (8 hours)
2. WinRM communication layer (4 hours)
3. Test runner integration (4 hours)
4. CI/CD pipeline (2 hours)
5. Documentation (2 hours)

**Total**: ~20 hours of work

**Risk**: High - complex, many moving parts

**Timeline**: Q4 2025 or when Windows support becomes critical

---

## Recommendations

### Immediate Actions (Now)

1. ✅ **Refactor MSIX into new packaging system**
   - Follows same pattern as macOS
   - Pure Go, template-based
   - Testable on macOS (structure/templates)
   - Deprecate old commands

2. ✅ **Write Go unit tests**
   - Template generation
   - Bundle structure
   - Platform detection

3. ✅ **Document limitations**
   - Clear warnings: "Actual MSIX creation only works on Windows"
   - Manual testing guide
   - CI skips Windows packaging tests for now

### Short-term (Next 1-2 months)

1. ⏳ **Manual VM testing**
   - Create Windows 11 UTM VM
   - Test checklist before releases
   - Document any Windows-specific issues

2. ⏳ **Improve error messages**
   - Graceful failures on macOS
   - Clear instructions for Windows users

### Long-term (Q4 2025)

1. 🚀 **Automated VM testing**
   - Only if Windows support becomes critical
   - Only if we have multiple Windows users
   - Cost/benefit analysis needed

---

## Questions for You

1. **Priority**: How important is Windows support right now?
   - If low: Just refactor, skip VM testing for now
   - If high: Set up manual VM testing workflow

2. **VM Setup**: Do you already have a Windows 11 UTM VM?
   - If yes: I can write the manual test checklist
   - If no: We can skip Windows testing entirely for now

3. **Communication Layer**: When we do automated testing, prefer:
   - WinRM (more native, requires network setup)
   - SSH (more universal, requires OpenSSH on Windows)
   - Shared folders (simpler, less reliable)

4. **MSIX Toolkit**: Where does the `msix` command come from?
   - Official Microsoft tool?
   - Third-party package?
   - Need to document installation for Windows users

5. **Code Signing**: Do you have/need Windows code signing certificates?
   - If yes: Need to handle .pfx files and passwords
   - If no: MSIX will be unsigned (OK for testing/internal use)

6. **Screenshot Support on Windows**: robotgo works on Windows!
   - Requires CGO (same as macOS)
   - Cross-compilation will be tricky (CGO + Windows target)
   - Need MinGW-w64 for cross-compiling from macOS
   - Alternative: Build on Windows directly (in VM)

---

## Summary

**Current MSIX code is actually pretty good** - already template-based and follows good patterns!

**Main issue is testing** - can't verify on macOS.

**Recommendation**:
1. ✅ Refactor into new packaging system (7 hours, low risk)
2. ⏸️ Defer VM automation until Windows becomes critical
3. 📝 Document limitations clearly
4. 🧪 Write Go tests for what we can test (templates, structure)

**Decision needed**: How much do we care about Windows right now?
