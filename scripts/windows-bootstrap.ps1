# Windows VM Bootstrap Script
# Sets up a fresh Windows 11 machine for goup-util development and testing
#
# Usage:
#   1. Inside Windows VM, open PowerShell as Administrator
#   2. Set-ExecutionPolicy Bypass -Scope Process -Force
#   3. .\windows-bootstrap.ps1
#
# Or one-liner from internet:
#   iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/joeblew999/goup-util/main/scripts/windows-bootstrap.ps1'))

Write-Host "============================================" -ForegroundColor Cyan
Write-Host " goup-util Windows VM Bootstrap" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan
Write-Host ""

# Check if running as administrator
$currentPrincipal = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
$isAdmin = $currentPrincipal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)

if (-not $isAdmin) {
    Write-Host "ERROR: This script must be run as Administrator" -ForegroundColor Red
    Write-Host "Right-click PowerShell and select 'Run as Administrator'" -ForegroundColor Yellow
    exit 1
}

Write-Host "[1/6] Checking winget..." -ForegroundColor Green

# Check if winget is installed
try {
    $wingetVersion = winget --version
    Write-Host "  winget is installed: $wingetVersion" -ForegroundColor Gray
} catch {
    Write-Host "  ERROR: winget not found" -ForegroundColor Red
    Write-Host "  winget should be pre-installed on Windows 11" -ForegroundColor Yellow
    Write-Host "  Install from: https://aka.ms/getwinget" -ForegroundColor Yellow
    exit 1
}

Write-Host ""
Write-Host "[2/6] Installing Git..." -ForegroundColor Green

# Install Git
Write-Host "  Installing Git.Git..." -ForegroundColor Gray
winget install --id Git.Git --exact --accept-source-agreements --accept-package-agreements --silent

# Refresh PATH
$env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")

Write-Host "  Git installed" -ForegroundColor Gray

Write-Host ""
Write-Host "[3/6] Installing Go..." -ForegroundColor Green

# Install Go
Write-Host "  Installing GoLang.Go..." -ForegroundColor Gray
winget install --id GoLang.Go --exact --accept-source-agreements --accept-package-agreements --silent

# Refresh PATH
$env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")

Write-Host "  Go installed" -ForegroundColor Gray

Write-Host ""
Write-Host "[4/6] Installing Task (Taskfile)..." -ForegroundColor Green

# Install Task
Write-Host "  Installing Task.Task..." -ForegroundColor Gray
winget install --id Task.Task --exact --accept-source-agreements --accept-package-agreements --silent

# Refresh PATH
$env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")

Write-Host "  Task installed" -ForegroundColor Gray

Write-Host ""
Write-Host "[5/6] Installing MSIX Packaging Tool..." -ForegroundColor Green

# Install MSIX Packaging Tool
Write-Host "  Installing Microsoft.MsixPackagingTool..." -ForegroundColor Gray
winget install --id Microsoft.MsixPackagingTool --exact --accept-source-agreements --accept-package-agreements --silent

Write-Host "  MSIX Packaging Tool installed" -ForegroundColor Gray

Write-Host ""
Write-Host "[6/6] Installing VSCode (optional)..." -ForegroundColor Green

# Ask if user wants VSCode
$installVSCode = Read-Host "  Install VSCode? (y/n)"
if ($installVSCode -eq "y" -or $installVSCode -eq "Y") {
    Write-Host "  Installing Microsoft.VisualStudioCode..." -ForegroundColor Gray
    winget install --id Microsoft.VisualStudioCode --exact --accept-source-agreements --accept-package-agreements --silent
    Write-Host "  VSCode installed" -ForegroundColor Gray
} else {
    Write-Host "  Skipping VSCode" -ForegroundColor Gray
}

Write-Host ""
Write-Host "============================================" -ForegroundColor Cyan
Write-Host " Installation Complete!" -ForegroundColor Green
Write-Host "============================================" -ForegroundColor Cyan
Write-Host ""

# Verify installations
Write-Host "Verifying installations:" -ForegroundColor Yellow
Write-Host ""

# Refresh PATH one more time
$env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")

# Git
try {
    $gitVersion = git --version
    Write-Host "  Git: $gitVersion" -ForegroundColor Green
} catch {
    Write-Host "  Git: NOT FOUND (restart terminal)" -ForegroundColor Red
}

# Go
try {
    $goVersion = go version
    Write-Host "  Go: $goVersion" -ForegroundColor Green
} catch {
    Write-Host "  Go: NOT FOUND (restart terminal)" -ForegroundColor Red
}

# Task
try {
    $taskVersion = task --version
    Write-Host "  Task: $taskVersion" -ForegroundColor Green
} catch {
    Write-Host "  Task: NOT FOUND (restart terminal)" -ForegroundColor Red
}

# MSIX
try {
    $msixHelp = msix --help 2>&1
    if ($msixHelp -match "msix") {
        Write-Host "  MSIX: Installed" -ForegroundColor Green
    } else {
        Write-Host "  MSIX: NOT FOUND" -ForegroundColor Red
    }
} catch {
    Write-Host "  MSIX: NOT FOUND" -ForegroundColor Red
}

Write-Host ""
Write-Host "============================================" -ForegroundColor Cyan
Write-Host " Installing goup-util" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan
Write-Host ""

# Determine architecture
$arch = if ([System.Environment]::Is64BitOperatingSystem) { "amd64" } else { "arm64" }
$binaryName = "goup-util-windows-$arch.exe"

Write-Host "Detecting platform: windows-$arch" -ForegroundColor Gray
Write-Host "Downloading latest goup-util release..." -ForegroundColor Green

try {
    # Get latest release info from GitHub
    $releaseUrl = "https://api.github.com/repos/joeblew999/goup-util/releases/latest"
    $release = Invoke-RestMethod -Uri $releaseUrl

    # Find the Windows binary
    $asset = $release.assets | Where-Object { $_.name -eq $binaryName }

    if (-not $asset) {
        Write-Host "  ERROR: Binary not found for windows-$arch" -ForegroundColor Red
        Write-Host "  Available assets:" -ForegroundColor Yellow
        $release.assets | ForEach-Object { Write-Host "    - $($_.name)" -ForegroundColor Gray }
        Write-Host ""
        Write-Host "  Fallback: Clone and build from source:" -ForegroundColor Yellow
        Write-Host "    git clone https://github.com/joeblew999/goup-util.git" -ForegroundColor Gray
        Write-Host "    cd goup-util" -ForegroundColor Gray
        Write-Host "    go build ." -ForegroundColor Gray
        exit 1
    }

    # Download to local directory
    $installPath = "$env:USERPROFILE\goup-util.exe"
    Write-Host "  Downloading $binaryName..." -ForegroundColor Gray
    Invoke-WebRequest -Uri $asset.browser_download_url -OutFile $installPath

    # Verify download
    if (Test-Path $installPath) {
        $fileSize = (Get-Item $installPath).Length / 1MB
        Write-Host "  Downloaded: $([math]::Round($fileSize, 2)) MB" -ForegroundColor Gray
        Write-Host "  Installed to: $installPath" -ForegroundColor Green

        # Handle Windows SmartScreen (unsigned binary)
        Write-Host ""
        Write-Host "  Windows Security Notice:" -ForegroundColor Yellow
        Write-Host "  The binary is not code-signed yet (no signing certificate)." -ForegroundColor Gray
        Write-Host "  Unblocking downloaded file..." -ForegroundColor Gray

        # Unblock the file
        Unblock-File -Path $installPath

        # Add to PATH if not already there
        $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
        $installDir = Split-Path $installPath
        if ($userPath -notlike "*$installDir*") {
            Write-Host "  Adding to PATH..." -ForegroundColor Gray
            [Environment]::SetEnvironmentVariable("Path", "$userPath;$installDir", "User")
            $env:Path = "$env:Path;$installDir"
        }

        # Test the binary
        Write-Host ""
        Write-Host "  Testing goup-util..." -ForegroundColor Green

        try {
            & $installPath --help 2>&1 | Select-Object -First 5
            Write-Host ""
            Write-Host "  goup-util successfully installed and working!" -ForegroundColor Green
        } catch {
            Write-Host ""
            Write-Host "  If you see Windows SmartScreen warning:" -ForegroundColor Yellow
            Write-Host "    1. Click 'More info'" -ForegroundColor Gray
            Write-Host "    2. Click 'Run anyway'" -ForegroundColor Gray
            Write-Host "    3. Or right-click the file → Properties → Unblock" -ForegroundColor Gray
            Write-Host ""
            Write-Host "  To bypass SmartScreen for this file:" -ForegroundColor Gray
            Write-Host "    Unblock-File -Path $installPath" -ForegroundColor Gray
            Write-Host ""
            Write-Host "  Note: Future releases will be code-signed." -ForegroundColor Yellow
        }
    } else {
        Write-Host "  ERROR: Download failed" -ForegroundColor Red
        exit 1
    }

} catch {
    Write-Host "  ERROR: Failed to download goup-util" -ForegroundColor Red
    Write-Host "  Error: $_" -ForegroundColor Red
    Write-Host ""
    Write-Host "  Fallback: Clone and build from source:" -ForegroundColor Yellow
    Write-Host "    git clone https://github.com/joeblew999/goup-util.git" -ForegroundColor Gray
    Write-Host "    cd goup-util" -ForegroundColor Gray
    Write-Host "    go build ." -ForegroundColor Gray
    exit 1
}

Write-Host ""
Write-Host "============================================" -ForegroundColor Cyan
Write-Host " Setup Complete!" -ForegroundColor Green
Write-Host "============================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Cyan
Write-Host "  1. Test goup-util:"
Write-Host "     goup-util --help"
Write-Host ""
Write-Host "  2. Clone example projects:"
Write-Host "     git clone https://github.com/joeblew999/goup-util.git"
Write-Host "     cd goup-util"
Write-Host ""
Write-Host "  3. Test Windows bundling:"
Write-Host "     goup-util build windows examples\hybrid-dashboard"
Write-Host "     goup-util bundle windows examples\hybrid-dashboard --create-msix"
Write-Host ""
Write-Host "  4. Update goup-util anytime:"
Write-Host "     goup-util self upgrade"
Write-Host ""
Write-Host "For automated testing, see: docs/WINDOWS-PACKAGING-STRATEGY.md"
Write-Host ""
