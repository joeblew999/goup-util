# Windows Testing Checklist

Manual testing checklist for goup-util Windows functionality on UTM VM.

## VM Setup (One-Time)

### 1. Create Windows 11 VM in UTM

**Download Windows 11 ISO**:
- https://www.microsoft.com/software-download/windows11
- Select "Windows 11 (multi-edition ISO)"
- Download ARM64 version (for Apple Silicon Macs)

**Create VM in UTM**:
1. Open UTM
2. Create a New Virtual Machine
3. Virtualize (not Emulate)
4. Windows
5. Browse to Windows 11 ISO
6. Memory: 4GB minimum, 8GB recommended
7. Storage: 64GB minimum
8. Shared Directory: Optional (for file transfer)
9. Create

**Install Windows**:
1. Start VM
2. Follow Windows 11 setup
3. Skip Microsoft account (use local account)
4. Username: `developer`
5. Minimal privacy settings

### 2. Bootstrap Development Environment

**Run bootstrap script**:
```powershell
# Open PowerShell as Administrator
Set-ExecutionPolicy Bypass -Scope Process -Force

# Download and run bootstrap
iwr https://raw.githubusercontent.com/joeblew999/goup-util/main/scripts/windows-bootstrap.ps1 -UseBasicParsing | iex

# Or if offline, copy bootstrap script to VM and run
.\windows-bootstrap.ps1
```

**What gets installed**:
- ✅ Git
- ✅ Go
- ✅ Task (Taskfile runner)
- ✅ MSIX Packaging Tool
- ✅ VSCode (optional)

**Verify installation**:
```powershell
git --version
go version
task --version
msix --help
```

### 3. Clone and Build goup-util

```powershell
# Clone repository
cd ~\Documents
git clone https://github.com/joeblew999/goup-util.git
cd goup-util

# Build goup-util
go build .

# Verify
.\goup-util.exe --version
```

---

## Testing Workflow

### Before Each Test Session

1. **Update repository**:
```powershell
cd ~\Documents\goup-util
git fetch origin
git pull origin main
```

2. **Rebuild goup-util**:
```powershell
go build .
```

3. **Clean previous test artifacts**:
```powershell
rm -Recurse -Force examples\*\.bin\*
rm -Recurse -Force examples\*\.dist\*
rm -Recurse -Force examples\*\.staging\*
```

---

## Test Suite

### Test 1: Build Windows Binary

**Objective**: Verify cross-platform build produces working Windows .exe

```powershell
# Build hybrid-dashboard for Windows
.\goup-util.exe build windows examples\hybrid-dashboard

# Expected output:
# - Building hybrid-dashboard for windows...
# - ✓ Built successfully
# - Output: examples\hybrid-dashboard\.bin\hybrid-dashboard.exe

# Verify binary exists
Test-Path examples\hybrid-dashboard\.bin\hybrid-dashboard.exe

# Run the binary (should launch Gio app)
.\examples\hybrid-dashboard\.bin\hybrid-dashboard.exe
```

**Success criteria**:
- [ ] Binary builds without errors
- [ ] Binary runs and shows GUI window
- [ ] No runtime crashes

**Screenshot**: Capture running app

---

### Test 2: Bundle Windows Package (Structure Only)

**Objective**: Verify bundle command creates correct directory structure

```powershell
# Create bundle structure (no MSIX yet)
.\goup-util.exe bundle windows examples\hybrid-dashboard `
    --bundle-id hybrid-dashboard `
    --publisher "CN=TestPublisher" `
    --version 1.0.0.0

# Expected output:
# - Creating Windows bundle for hybrid-dashboard...
# - ✓ Binary copied: hybrid-dashboard.exe
# - ✓ Generated placeholder assets
# - ✓ AppxManifest.xml created
# - ⚠️  Skipping MSIX creation: requires Windows
# - Bundle structure created, ready to package on Windows

# Verify structure
Test-Path examples\hybrid-dashboard\.dist\.staging\hybrid-dashboard.exe
Test-Path examples\hybrid-dashboard\.dist\.staging\AppxManifest.xml
Test-Path examples\hybrid-dashboard\.dist\.staging\assets\logo.png
```

**Success criteria**:
- [ ] Bundle structure created
- [ ] AppxManifest.xml is valid XML
- [ ] Binary copied correctly
- [ ] Assets directory exists

**Verify AppxManifest.xml**:
```powershell
Get-Content examples\hybrid-dashboard\.dist\.staging\AppxManifest.xml

# Check for:
# - <Identity Name="hybrid-dashboard" Publisher="CN=TestPublisher" Version="1.0.0.0" />
# - <Application Executable="hybrid-dashboard.exe" ... />
```

---

### Test 3: Create MSIX Package

**Objective**: Verify MSIX packaging works on Windows

```powershell
# Create MSIX with --create-msix flag
.\goup-util.exe bundle windows examples\hybrid-dashboard `
    --bundle-id hybrid-dashboard `
    --publisher "CN=TestPublisher" `
    --version 1.0.0.0 `
    --create-msix

# Expected output:
# - Creating Windows bundle for hybrid-dashboard...
# - ✓ Binary copied: hybrid-dashboard.exe
# - ✓ Generated placeholder assets
# - ✓ AppxManifest.xml created
# - ✓ MSIX package created: examples\hybrid-dashboard\.dist\hybrid-dashboard.msix

# Verify MSIX exists
Test-Path examples\hybrid-dashboard\.dist\hybrid-dashboard.msix

# Check MSIX size (should be > 10MB)
(Get-Item examples\hybrid-dashboard\.dist\hybrid-dashboard.msix).Length / 1MB
```

**Success criteria**:
- [ ] MSIX file created
- [ ] File size is reasonable (> 10MB)
- [ ] No errors during packaging

---

### Test 4: Install MSIX Package

**Objective**: Verify MSIX can be installed on Windows

```powershell
# Install the MSIX package
Add-AppxPackage examples\hybrid-dashboard\.dist\hybrid-dashboard.msix

# Expected output:
# - (silent success)

# Verify installation
Get-AppxPackage | Where-Object { $_.Name -like "*hybrid-dashboard*" }

# Expected output:
# - Name: hybrid-dashboard
# - Publisher: CN=TestPublisher
# - Version: 1.0.0.0
# - InstallLocation: C:\Program Files\WindowsApps\...
```

**Success criteria**:
- [ ] Package installs without errors
- [ ] Package appears in installed apps list
- [ ] Package has correct metadata

**Launch installed app**:
```powershell
# Find Start Menu entry
# Start → All Apps → hybrid-dashboard

# Or launch directly
Start-Process "shell:AppsFolder\$(Get-AppxPackage | Where-Object { $_.Name -like "*hybrid-dashboard*" } | Select-Object -ExpandProperty PackageFamilyName)!App"
```

**Success criteria**:
- [ ] App launches from Start Menu
- [ ] App functions correctly
- [ ] No permission errors

---

### Test 5: Uninstall MSIX Package

**Objective**: Verify clean uninstall

```powershell
# Uninstall package
Get-AppxPackage | Where-Object { $_.Name -like "*hybrid-dashboard*" } | Remove-AppxPackage

# Verify removal
Get-AppxPackage | Where-Object { $_.Name -like "*hybrid-dashboard*" }

# Expected output:
# - (nothing - package removed)
```

**Success criteria**:
- [ ] Package uninstalls cleanly
- [ ] No orphaned files
- [ ] No registry errors

---

### Test 6: Package Command (Distribution Archive)

**Objective**: Verify package command creates distribution files

```powershell
# Create distribution package
.\goup-util.exe package windows examples\hybrid-dashboard

# Expected output:
# - Packaging hybrid-dashboard for Windows distribution...
# - ✓ Packaged hybrid-dashboard for Windows: examples\hybrid-dashboard\.dist\hybrid-dashboard-windows.zip

# Verify ZIP exists
Test-Path examples\hybrid-dashboard\.dist\hybrid-dashboard-windows.zip

# Extract and verify contents
Expand-Archive examples\hybrid-dashboard\.dist\hybrid-dashboard-windows.zip -DestinationPath .\test-extract
Test-Path test-extract\hybrid-dashboard.exe
```

**Success criteria**:
- [ ] ZIP file created
- [ ] ZIP contains .exe
- [ ] .exe is executable

**Cleanup**:
```powershell
rm -Recurse -Force test-extract
```

---

### Test 7: Idempotency

**Objective**: Verify build cache works on Windows

```powershell
# First build (full compilation)
Measure-Command { .\goup-util.exe build windows examples\hybrid-dashboard }

# Second build (should skip)
Measure-Command { .\goup-util.exe build windows examples\hybrid-dashboard }

# Expected output:
# - ✓ hybrid-dashboard for windows is up-to-date (use --force to rebuild)
# - Should be MUCH faster

# Force rebuild
.\goup-util.exe build --force windows examples\hybrid-dashboard
```

**Success criteria**:
- [ ] Second build is significantly faster
- [ ] Second build shows "up-to-date" message
- [ ] --force flag triggers rebuild

---

### Test 8: Screenshot Functionality (CGO)

**Objective**: Verify robotgo screenshot support on Windows

**Note**: This requires CGO compilation which is complex on Windows.

```powershell
# Build with CGO (requires MinGW-w64)
$env:CGO_ENABLED = "1"
go build -tags screenshot .

# Test screenshot
.\goup-util.exe screenshot test-screenshot.png

# Verify screenshot created
Test-Path test-screenshot.png

# View screenshot
Invoke-Item test-screenshot.png
```

**Success criteria**:
- [ ] Screenshot command works
- [ ] Image is valid PNG
- [ ] Image shows desktop content

**Note**: If CGO compilation fails:
- Screenshot functionality may need to be disabled on Windows
- Or build on Windows with proper MinGW-w64 setup
- See: docs/WINDOWS-PACKAGING-STRATEGY.md

---

## Reporting Results

### Screenshot Collection

Take screenshots of:
1. ✅ Successful build output
2. ✅ Running Gio app (hybrid-dashboard.exe)
3. ✅ MSIX installation confirmation
4. ✅ Installed app in Start Menu
5. ❌ Any errors encountered

Save to: `docs/screenshots/windows/`

### Test Report Template

```markdown
## Windows Testing Report

**Date**: YYYY-MM-DD
**VM**: UTM Windows 11 ARM64
**goup-util version**: (git commit hash)

### Test Results

- [ ] Test 1: Build Windows Binary
- [ ] Test 2: Bundle Windows Package
- [ ] Test 3: Create MSIX Package
- [ ] Test 4: Install MSIX Package
- [ ] Test 5: Uninstall MSIX Package
- [ ] Test 6: Package Command
- [ ] Test 7: Idempotency
- [ ] Test 8: Screenshot Functionality

### Issues Found

1. (Issue description)
   - Error: (error message)
   - Workaround: (if any)

### Notes

(Any observations, performance issues, etc.)
```

Save report to: `docs/test-reports/windows-YYYY-MM-DD.md`

---

## Troubleshooting

### Issue: msix command not found

```powershell
# Reinstall MSIX Packaging Tool
winget install Microsoft.MsixPackagingTool

# Verify PATH
$env:Path -split ";" | Select-String "msix"
```

### Issue: CGO compilation fails

```powershell
# Install MinGW-w64
winget install mingw-w64

# Or disable screenshot feature
go build -tags="!screenshot" .
```

### Issue: MSIX installation fails

```powershell
# Enable developer mode
Start-Process ms-settings:developers

# Or sign the MSIX with a test certificate
# (future implementation)
```

### Issue: Git clone fails

```powershell
# Configure Git
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"

# Use HTTPS instead of SSH
git clone https://github.com/joeblew999/goup-util.git
```

---

## Cleanup

After testing:

```powershell
# Remove test artifacts
cd ~\Documents\goup-util
rm -Recurse -Force examples\*\.bin
rm -Recurse -Force examples\*\.dist
rm -Recurse -Force *.png

# Uninstall test packages
Get-AppxPackage | Where-Object { $_.Name -like "*hybrid*" -or $_.Name -like "*test*" } | Remove-AppxPackage
```

---

## Future: Automated Testing

See `docs/WINDOWS-PACKAGING-STRATEGY.md` for plans to automate this with:
- HTTP API server in VM
- Test runner on macOS host
- CI/CD integration

**Status**: Manual testing only for now (Q4 2025 for automation)
