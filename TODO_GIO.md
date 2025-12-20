# Gio & gio-plugins TODO

New features and investigations from the v0.9.1 update (2025-12-20).

## High Priority - Deep Linking / Custom URI Schemes

### Gio UI - Custom URI Scheme Support (NEW!)

**What it is:** Apps can now be launched via custom URI schemes like `gio://some/data` on Android, iOS, macOS, and Windows.

**Why it matters for goup-util:**
- Enable deep linking for hybrid apps
- Support `myapp://dashboard/stats` URLs
- Universal links for marketing/sharing
- Critical for wellknown integration

**Investigate:**
- [ ] How to configure custom URI scheme in app bundle
- [ ] How to receive and parse URI data in Gio app
- [ ] Integration with gio-plugins webviewer (pass URLs to webview)
- [ ] Platform-specific setup (Info.plist, AndroidManifest, etc.)
- [ ] Add `goup-util init --uri-scheme myapp://` support

**Branch to watch:** `deeplink2025` in gio-plugins

### gio-plugins - Deep Linking Branch

**What it is:** Active development on `deeplink2025` branch for improved deep linking.

**Investigate:**
- [ ] What's in the deeplink2025 branch?
- [ ] Universal links support?
- [ ] App links (Android) support?
- [ ] Integration with wellknown system

## Medium Priority - Platform Improvements

### Windows Touch Screen Support

**What it is:** Windows Pointer API for detecting touch screen interactions. Fourth and fifth mouse button support.

**Why it matters:**
- Better tablet/touch support on Windows
- More input options for hybrid apps

**Investigate:**
- [ ] Test touch events in webviewer on Windows tablets
- [ ] Verify mouse button support works through webview
- [ ] Document touch-specific UI patterns

### macOS Fullscreen MaxSize

**What it is:** macOS fullscreen mode now respects MaxSize window constraint.

**Why it matters:**
- Better control over app window sizing
- Important for fixed-layout hybrid apps

**Investigate:**
- [ ] Test MaxSize with webviewer apps
- [ ] Document fullscreen behavior

## Bug Fixes to Verify

### Android Text Rendering

**What was fixed:** GPU changes that broke text rendering on some Android devices were reverted.

**Verify:**
- [ ] Test webviewer on various Android devices
- [ ] Check text rendering in native Gio UI portions

### GPU Clipping

**What was fixed:** Vertex corner positions causing 1px overlaps and skewed rendering.

**Verify:**
- [ ] Visual inspection of hybrid app UI
- [ ] Check for rendering artifacts at element boundaries

## gio-plugins Specific

### Auth Global Event Listener (#106)

**What was fixed:** Global event listener issues in authentication module.

**Investigate:**
- [ ] Test OAuth flows (Apple, Google sign-in)
- [ ] Verify auth callbacks work correctly
- [ ] Update auth example if needed

### Updated Dependencies (#105)

**What changed:** gio-plugins now requires Gio v0.9.1.

**Done:**
- [x] Updated examples/gio-plugin-webviewer/go.mod
- [x] Updated CLAUDE.md version compatibility section
- [x] Updated Taskfile.yml version variables

## Future Integration Ideas

### URI Scheme + HTMX Pattern

```go
// App launched via: myapp://dashboard?view=sales
func handleDeepLink(uri string) {
    // Parse URI
    // Load corresponding HTMX view
    // webview.Navigate("http://localhost:8080" + uri.Path)
}
```

### wellknown Integration

```
/.well-known/apple-app-site-association → iOS universal links
/.well-known/assetlinks.json → Android app links
```

goup-util could auto-generate these files during build.

### goup-util Commands to Add

```bash
# Initialize project with URI scheme
goup-util init myapp --uri-scheme myapp://

# Generate universal/app links files
goup-util wellknown generate --domain example.com

# Validate deep linking setup
goup-util doctor --check-deeplinks
```

## Resources

- Gio commits: https://github.com/gioui/gio/commits/main
- gio-plugins v0.9.1: https://github.com/gioui-plugins/gio-plugins/releases/tag/v0.9.1
- deeplink2025 branch: https://github.com/gioui-plugins/gio-plugins/tree/deeplink2025
- Issue #104 (version compat): https://github.com/gioui-plugins/gio-plugins/issues/104
- Issue #106 (auth fix): https://github.com/gioui-plugins/gio-plugins/issues/106

## Version Reference

**Current compatible versions (2025-12-20):**
```bash
go get gioui.org@7bcb315ee174
go get github.com/gioui-plugins/gio-plugins@v0.9.1
go mod tidy
```

Resulting in:
- `gioui.org v0.9.1-0.20251215212054-7bcb315ee174`
- `github.com/gioui-plugins/gio-plugins v0.9.1`
