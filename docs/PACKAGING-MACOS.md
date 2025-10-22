# macOS Packaging and Code Signing

This guide explains how to properly package goup-util for macOS with code signing to handle Screen Recording permissions.

## Why Packaging Matters

For screenshot functionality with robotgo, macOS requires:

1. ✅ **Proper app bundle structure** (`.app` folder)
2. ✅ **Code signing** (even ad-hoc for development)
3. ✅ **Consistent bundle ID** (so macOS remembers permissions)
4. ⚠️ **User grants Screen Recording permission** (no automatic entitlement)

## Quick Start

```bash
# Build and package as a signed macOS app
task package:macos

# The app bundle will be created at:
# .dist/goup-util.app

# Open it (this triggers the permission dialog)
open .dist/goup-util.app

# Grant permission when prompted:
# System Settings → Privacy & Security → Screen Recording
```

## Understanding macOS Screen Recording Permissions

### No Automatic Entitlement 🚫

**Important**: Unlike Camera or Microphone permissions, there is **NO entitlement** that allows automatic screen recording access for third-party apps.

The private entitlements like:
- `com.apple.private.screencapture.allow`
- `com.apple.private.tcc.allow`

...are **ONLY available to Apple's own applications**.

### What About Info.plist?

There is **NO `NSScreenCaptureUsageDescription`** key for screen recording permissions. This is different from other privacy permissions:

| Permission | Info.plist Key | Available? |
|------------|----------------|------------|
| Camera | `NSCameraUsageDescription` | ✅ Yes |
| Microphone | `NSMicrophoneUsageDescription` | ✅ Yes |
| Location | `NSLocationUsageDescription` | ✅ Yes |
| **Screen Recording** | ~~`NSScreenCaptureUsageDescription`~~ | ❌ **No** |

### The Process

1. **App is properly signed** (with any valid signature, even ad-hoc `-`)
2. **App attempts screen capture** (triggers permission check)
3. **System shows standard dialog** (cannot be customized)
4. **User manually grants permission** (in System Settings)
5. **Permission is remembered** (as long as bundle ID + signature stay same)

## Build and Packaging Options

### Option 1: Automated Script (Recommended)

Use the provided build script:

```bash
# Build and package with automatic signing detection
bash pkg/packaging/deprecated/build-macos.sh  # DEPRECATED - use 'goup-util bundle' instead

# Or via Taskfile
task package:macos
```

The script will:
- ✅ Build with CGO enabled (for robotgo)
- ✅ Create proper app bundle structure
- ✅ Copy Info.plist with bundle ID
- ✅ Auto-detect signing certificate
- ✅ Apply entitlements if certificate found
- ✅ Fall back to ad-hoc signature if needed

### Option 2: Manual Packaging

```bash
# 1. Build binary
CGO_ENABLED=1 go build -o goup-util .

# 2. Create app bundle
mkdir -p .dist/goup-util.app/Contents/MacOS
mkdir -p .dist/goup-util.app/Contents/Resources

# 3. Copy binary
cp goup-util .dist/goup-util.app/Contents/MacOS/goup-util
chmod +x .dist/goup-util.app/Contents/MacOS/goup-util

# 4. Copy Info.plist
cp pkg/packaging/deprecated/macos-info.plist .dist/goup-util.app/Contents/Info.plist  # OLD - now uses templates

# 5. Sign (ad-hoc)
codesign --force --deep --sign - .dist/goup-util.app

# 6. Verify
codesign --verify --deep --strict .dist/goup-util.app
```

## Code Signing Options

### Ad-Hoc Signature (Development)

**Pros**:
- ✅ No Apple Developer account needed
- ✅ Free
- ✅ Works for local testing
- ✅ Permissions are remembered (same machine)

**Cons**:
- ⚠️ Can't distribute to other machines
- ⚠️ Gatekeeper will block on other Macs

```bash
codesign --force --deep --sign - .dist/goup-util.app
```

### Apple ID Signature (Distribution)

**Pros**:
- ✅ Free Apple Developer account
- ✅ Can distribute to others
- ✅ Permissions persist correctly
- ✅ Works with notarization

**Cons**:
- ⚠️ Requires Apple ID
- ⚠️ Manual Gatekeeper approval by users

```bash
# Sign with your Apple ID
codesign --force --deep \
  --sign "YOUR_APPLE_ID_EMAIL" \
  --entitlements pkg/packaging/deprecated/entitlements.plist \
  --options runtime \
  .dist/goup-util.app
```

### Developer ID Application (Production)

**Pros**:
- ✅ Proper distribution
- ✅ Can be notarized by Apple
- ✅ No Gatekeeper warnings
- ✅ Professional

**Cons**:
- ⚠️ Requires paid Apple Developer Program ($99/year)

```bash
# Find your Developer ID
security find-identity -v -p codesigning | grep "Developer ID Application"

# Sign with Developer ID
codesign --force --deep \
  --sign "Developer ID Application: Your Name (TEAM_ID)" \
  --entitlements pkg/packaging/deprecated/entitlements.plist \
  --options runtime \
  .dist/goup-util.app

# Notarize (requires paid account)
xcrun notarytool submit .dist/goup-util.app \
  --apple-id "your@email.com" \
  --team-id "TEAM_ID" \
  --password "app-specific-password"
```

## Entitlements Explained

The `pkg/packaging/templates/macos-entitlements.plist.tmpl` includes:

```xml
<!-- Hardened Runtime -->
<key>com.apple.security.cs.allow-jit</key>
<key>com.apple.security.cs.allow-unsigned-executable-memory</key>
<key>com.apple.security.cs.disable-library-validation</key>

<!-- Network (for SDK downloads) -->
<key>com.apple.security.network.client</key>
<key>com.apple.security.network.server</key>

<!-- File Access -->
<key>com.apple.security.files.user-selected.read-write</key>
<key>com.apple.security.files.downloads.read-write</key>
```

**Note**: There is NO screen recording entitlement here because it doesn't exist for third-party apps!

## Bundle ID Importance

The bundle ID (`com.joeblew999.goup-util`) is **critical**:

- ✅ macOS uses it to remember permissions
- ✅ Changing it = losing all granted permissions
- ✅ Must be unique across your apps
- ✅ Should match code signing certificate

If you change the bundle ID or re-sign with a different certificate, users must grant permission again!

## Distribution Options

### 1. DMG Installer

```bash
# Create DMG
task package:macos:dmg

# Or manually
hdiutil create -volname "goup-util" \
  -srcfolder .dist \
  -ov -format UDZO \
  .dist/goup-util.dmg
```

Users can:
- Download the DMG
- Open it
- Drag `goup-util.app` to Applications
- Launch and grant permission

### 2. ZIP Archive

```bash
# Create ZIP
cd .dist
zip -r goup-util.zip goup-util.app
```

### 3. Homebrew Cask

```ruby
# homebrew-cask formula
cask "goup-util" do
  version "1.0.0"
  sha256 "..."

  url "https://github.com/joeblew999/goup-util/releases/download/v#{version}/goup-util.dmg"
  name "Goup Util"
  desc "Cross-platform Gio app build tool"
  homepage "https://github.com/joeblew999/goup-util"

  app "goup-util.app"
end
```

## Permission Workflow for End Users

### First Launch

1. User downloads and opens `goup-util.app`
2. User runs: `./goup-util screenshot test.png`
3. **System shows dialog**: "goup-util would like to record your screen"
4. User clicks "Allow"
5. Screenshot works! 🎉

### If Permission Denied

```bash
# User sees error:
# Error: screenshot failed: Capture image not found.
#
# Note: On macOS 10.15+, grant Screen Recording permission:
# System Settings → Privacy & Security → Screen Recording
```

User must:
1. Open **System Settings**
2. Go to **Privacy & Security** → **Screen Recording**
3. Check the box next to **goup-util**
4. Restart the app (if needed)

### Removing Permission

To revoke permission:
1. System Settings → Privacy & Security → Screen Recording
2. Uncheck goup-util
3. Or click the `-` button to remove it entirely

## Troubleshooting

### "goup-util.app is damaged and can't be opened"

**Cause**: Gatekeeper blocking unsigned or improperly signed apps

**Solution**:
```bash
# Remove quarantine attribute
xattr -dr com.apple.quarantine .dist/goup-util.app
```

### "Permission denied" even after granting

**Cause**: App needs to restart after permission granted

**Solution**: Fully quit and relaunch the app

### Permission resets after rebuild

**Cause**: Bundle ID or code signature changed

**Solution**: Use consistent signing identity:
```bash
# Always use the same identity
codesign --force --deep --sign "SAME_ID_EVERY_TIME" .dist/goup-util.app
```

### Can't find signing certificate

**Check available certificates**:
```bash
security find-identity -v -p codesigning
```

**Get a free Apple ID certificate**:
1. Sign in to Xcode with your Apple ID
2. Xcode → Settings → Accounts → Manage Certificates
3. Click `+` → Apple Development
4. Use this for signing

## CI/CD Considerations

### GitHub Actions

```yaml
name: Build and Package
on: [push]
jobs:
  build-macos:
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v4

      - name: Build and Package
        run: task package:macos

      - name: Ad-hoc sign (for CI)
        run: |
          codesign --force --deep --sign - .dist/goup-util.app

      - name: Upload DMG
        uses: actions/upload-artifact@v4
        with:
          name: goup-util-macos
          path: .dist/goup-util.dmg
```

### Notarization (Requires Paid Account)

```bash
# Build and sign
task package:macos

# Create DMG
task package:macos:dmg

# Notarize
xcrun notarytool submit .dist/goup-util.dmg \
  --apple-id "$APPLE_ID" \
  --team-id "$TEAM_ID" \
  --password "$APP_PASSWORD" \
  --wait

# Staple ticket
xcrun stapler staple .dist/goup-util.dmg
```

## Summary

✅ **DO**:
- Build as a proper `.app` bundle
- Sign with ANY signature (even ad-hoc)
- Keep bundle ID consistent
- Document permission requirements for users

❌ **DON'T**:
- Expect automatic screen recording access
- Try to use `NSScreenCaptureUsageDescription` (doesn't exist)
- Use private Apple entitlements (they won't work)
- Change bundle ID or signature unnecessarily

🎯 **Remember**: Screen recording permission is **always manual** on macOS. The best you can do is:
1. Proper packaging
2. Clear error messages
3. Good documentation

---

## Related Documentation

- [Screenshot Integration](SCREENSHOT.md) - Using the screenshot command
- [robotgo Guide](agents/robotgo.md) - Deep dive into robotgo
- [Apple Developer: Code Signing](https://developer.apple.com/documentation/security/code_signing_services)
- [Apple Developer: Hardened Runtime](https://developer.apple.com/documentation/security/hardened_runtime)
