# Packaging System

goup-util provides a three-tier packaging system for distributing Gio applications:

## The Three Tiers

### 1. Build - Compile the Application

```bash
goup-util build <platform> <app-directory>
```

**What it does:**
- Compiles Go source code to platform-specific binaries
- Creates basic app structures (.app bundles, APKs, etc.)
- **Idempotent**: Skips rebuild if source hasn't changed
- **Fast**: Uses build cache to avoid unnecessary work

**Output locations:**
- macOS: `<app>/.bin/<name>.app` (basic bundle)
- Android: `<app>/.bin/<name>.apk`
- iOS: `<app>/.bin/<name>.app`
- Windows: `<app>/.bin/<name>.exe`

**Flags:**
- `--force` - Force rebuild even if up-to-date
- `--check` - Check if rebuild needed (exit 0=no, 1=yes)

**Examples:**
```bash
# Build for macOS (skip if up-to-date)
goup-util build macos examples/hybrid-dashboard

# Force rebuild
goup-util build --force macos examples/hybrid-dashboard

# Check if rebuild needed
goup-util build --check macos examples/hybrid-dashboard
echo $?  # 0=up-to-date, 1=needs rebuild
```

---

### 2. Bundle - Create Signed App Bundles

```bash
goup-util bundle <platform> <app-directory> [flags]
```

**What it does:**
- Creates properly structured app bundles for distribution
- Generates platform-specific metadata (Info.plist, manifests)
- **Code signing** with auto-detection or specified identity
- Adds entitlements for macOS hardened runtime
- **Pure Go**: No bash scripts, cross-platform tool

**Output locations:**
- macOS: `<app>/.dist/<name>.app` (signed bundle)
- Android: `<app>/.dist/<name>.apk` (signed)
- iOS: `<app>/.dist/<name>.ipa` (coming soon)
- Windows: `<app>/.dist/<name>.exe` (signed, coming soon)

**Flags:**
- `--bundle-id` - Bundle identifier (default: com.example.<name>)
- `--version` - Version string (default: 1.0.0)
- `--sign` - Code signing identity (empty for auto-detect)
- `--entitlements` - Use entitlements (default: true)
- `--output` - Output directory (default: .dist/)

**Examples:**
```bash
# Create signed macOS bundle (auto-detect certificate)
goup-util bundle macos examples/hybrid-dashboard

# Use specific bundle ID
goup-util bundle macos examples/hybrid-dashboard --bundle-id com.mycompany.app

# Skip entitlements
goup-util bundle macos examples/hybrid-dashboard --entitlements=false

# Use specific signing identity
goup-util bundle macos examples/hybrid-dashboard --sign "Developer ID Application: Company Name"
```

**Code Signing:**
- Automatically detects "Developer ID Application" certificates
- Falls back to "Apple Development" certificates
- Uses ad-hoc signature (`-`) if no certificate found
- Ad-hoc is suitable for local testing, not distribution

---

### 3. Package - Create Distribution Archives

```bash
goup-util package <platform> <app-directory>
```

**What it does:**
- Creates compressed archives of signed bundles
- Ready for upload to app stores or direct distribution
- Uses pure Go archiving (no external tools)

**Output locations:**
- macOS: `<app>/.dist/<name>-macos.tar.gz`
- Android: `<app>/.dist/<name>-android.apk` (copy)
- iOS: `<app>/.dist/<name>-ios.tar.gz`
- Windows: `<app>/.dist/<name>-windows.zip`

**Examples:**
```bash
# Package macOS app for distribution
goup-util package macos examples/hybrid-dashboard

# Package Android app
goup-util package android examples/hybrid-dashboard
```

---

## Complete Workflow

### Local Development

```bash
# 1. Build and test (idempotent, fast iterations)
goup-util build macos examples/hybrid-dashboard
open examples/hybrid-dashboard/.bin/hybrid-dashboard.app

# Make changes to code...

# 2. Rebuild (automatic skip if no changes)
goup-util build macos examples/hybrid-dashboard
```

### Release Distribution

```bash
# 1. Build the app
goup-util build macos examples/hybrid-dashboard

# 2. Create signed bundle
goup-util bundle macos examples/hybrid-dashboard \
  --bundle-id com.mycompany.myapp \
  --version 1.0.0

# 3. Test the signed bundle
open examples/hybrid-dashboard/.dist/hybrid-dashboard.app

# 4. Package for distribution
goup-util package macos examples/hybrid-dashboard

# 5. Upload the tar.gz to your distribution channel
ls examples/hybrid-dashboard/.dist/*.tar.gz
```

### CI/CD Pipeline

```bash
# Check if rebuild needed (fast!)
if goup-util build --check macos examples/hybrid-dashboard; then
  echo "Up to date, skipping build"
else
  echo "Building..."
  goup-util build macos examples/hybrid-dashboard
  goup-util bundle macos examples/hybrid-dashboard
  goup-util package macos examples/hybrid-dashboard
fi
```

---

## Platform-Specific Notes

### macOS

**Build output:** Basic .app bundle in `.bin/`
- Contains: Binary + minimal Info.plist
- Not signed
- Works for local testing

**Bundle output:** Signed .app bundle in `.dist/`
- Contains: Binary + Info.plist + Entitlements
- Code signed (ad-hoc or certificate)
- Ready for distribution
- Includes hardened runtime entitlements

**Package output:** tar.gz archive
- Compressed bundle ready for upload
- Preserves code signature

**Code Signing Options:**
1. **Ad-hoc** (default if no certificate): `-`
   - Good for: Local testing
   - Not for: Distribution outside your organization

2. **Apple Development**: Auto-detected
   - Good for: Testing on your devices
   - Not for: Public distribution

3. **Developer ID Application**: Auto-detected (preferred)
   - Good for: Public distribution outside Mac App Store
   - Requires: Paid Apple Developer account

4. **Mac App Store**: Specify manually
   - Good for: Mac App Store submission
   - Requires: App Store Connect setup

### Android

**Build output:** APK in `.bin/`
- Debug-signed APK
- Works on emulators and test devices

**Bundle output:** (Coming soon)
- Release-signed APK
- Ready for Play Store

**Package output:** APK copy
- Same as bundle output

### iOS

**Build output:** .app in `.bin/`
- Unsigned bundle
- Works in iOS Simulator only

**Bundle output:** (Coming soon)
- Signed .app
- Ready for device testing or .ipa creation

**Package output:** tar.gz
- Archive of .app bundle

### Windows

**Build output:** .exe in `.bin/`
- Unsigned executable

**Bundle output:** (Coming soon)
- Signed .exe with manifest

**Package output:** zip
- Compressed executable

---

## Taskfile Integration

All packaging operations have corresponding Taskfile tasks:

### Build Tasks
```bash
task build:hybrid:macos          # Build (idempotent)
task build:hybrid:macos:force    # Force rebuild
task build:hybrid:macos:check    # Check if rebuild needed
```

### Bundle Tasks
```bash
task bundle:hybrid:macos         # Create signed bundle
```

### Package Tasks
```bash
task package:macos:dmg           # Create DMG (macOS only)
```

### CI Tasks
```bash
task ci:check                    # Check all examples
task ci:build                    # Build all (idempotent)
task ci:build:force             # Force rebuild all
```

---

## Implementation Details

### Pure Go Packaging

All packaging operations are implemented in pure Go:

- **pkg/packaging/macos.go** - macOS bundle creation
- **pkg/packaging/archive.go** - tar.gz and zip creation
- **pkg/packaging/templates/** - Info.plist and entitlements templates

No bash scripts, no external tools (except platform SDKs).

### Template System

Bundle metadata is generated from Go templates:

```go
// pkg/packaging/templates/macos-info.plist.tmpl
<key>CFBundleIdentifier</key>
<string>{{.BundleID}}</string>
```

Templates are embedded using `//go:embed` for distribution.

### Idempotency

Build operations are idempotent via `pkg/buildcache/`:
- SHA256 hashing of source files
- Timestamp tracking
- Output validation
- Smart rebuild decisions

---

## Migration from Old System

### Old Way (bash script)
```bash
bash pkg/packaging/build-macos.sh
```

### New Way (pure Go)
```bash
goup-util bundle macos examples/hybrid-dashboard
```

**The bash script is deprecated** but kept for reference in `pkg/packaging/deprecated/build-macos.sh`.

---

## Troubleshooting

### "Binary not found" error
```
Error: binary not found in .bin/myapp or .bin/myapp.app/Contents/MacOS/myapp
```
**Solution:** Run `goup-util build <platform> <app>` first

### Code signing failed
```
Error: codesign failed: errSecInternalComponent
```
**Solution:**
1. Check Keychain Access for valid certificates
2. Use `security find-identity -v -p codesigning` to list available identities
3. Try ad-hoc signing: `--sign -`

### "Up to date" but I changed code
```
✓ myapp for macos is up-to-date (use --force to rebuild)
```
**Solution:** The build cache doesn't detect all changes yet. Use `--force`:
```bash
goup-util build --force macos examples/myapp
```

---

## See Also

- [IDEMPOTENCY.md](IDEMPOTENCY.md) - Build caching details
- [SCREENSHOT.md](SCREENSHOT.md) - Screenshot feature setup
- [PACKAGING-MACOS.md](PACKAGING-MACOS.md) - macOS-specific guide
- [../CLAUDE.md](../CLAUDE.md) - Full project documentation
