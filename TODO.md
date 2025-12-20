# goup-util TODO

**Status**: Tool works! Builds hybrid apps for macOS, iOS, Android successfully.  
**Next**: Polish the experience and expand platform support.

See [docs/IMPROVEMENTS.md](docs/IMPROVEMENTS.md) for comprehensive improvement analysis.

---

## 🔥 High Priority (Do First)

### 1. Better Developer Experience
**Problem**: Builds are silent, errors are cryptic, no progress feedback  
**Impact**: Makes tool frustrating to use

- [ ] **Rich progress bars** with download speed, ETA, file sizes
- [ ] **Better error messages** with actionable suggestions and docs links
- [ ] **Structured logging** with --verbose mode
- [ ] **Build stage visibility** (Dependencies → Compilation → Packaging)

**Implementation**: Week 1-2 quick win from roadmap

---

### 2. Update Documentation
**Problem**: README doesn't reflect what tool actually does
**Impact**: People don't understand the value proposition

- [x] **Rewrite README** - show hybrid app capability ✅ (2025-10-21)
- [ ] **Quick start guide** that actually works end-to-end
- [ ] **Example showcase** - what you can build
- [x] **Platform support matrix** - what's tested vs what's possible ✅ (2025-10-21)

**Implementation**: 1-2 hours, do ASAP

---

### 2.5. Screenshot Command (NEW - High Priority)
**Problem**: Need automated screenshots for docs/marketing
**Impact**: Can't show what the tool produces visually

**Solution**: Add `goup-util screenshot` command using platform CLI tools

**Tasks**:
- [ ] Implement screenshot command in `cmd/screenshot.go`
  - Desktop: `screencapture` (macOS), PowerShell (Windows), `scrot` (Linux)
  - iOS Simulator: `xcrun simctl io booted screenshot`
  - Android: `adb exec-out screencap -p`
- [ ] Add Taskfile tasks: `screenshot:desktop`, `screenshot:ios`, `screenshot:android`, `screenshot:all`
- [ ] Generate screenshots for all examples to `docs/screenshots/`
- [ ] Update README with actual screenshots

**Strategy**: See [docs/SCREENSHOT-STRATEGY.md](docs/SCREENSHOT-STRATEGY.md)
**Approach**: Platform CLI tools (no CGO), optional robotgo integration later
**Implementation**: Week 1-2

**Note**: Considered go-vgo/robotgo (has screenshot + keyboard/mouse automation). Good for future, but starting simple with CLI tools for faster implementation and no CGO dependency.

---

### 3. Webview Integration Improvements
**Problem**: Core feature but not well documented/supported  
**Impact**: People can't build production hybrid apps

- [ ] **Go ↔ JavaScript bridge** - declarative API for function exposure
- [ ] **TypeScript definitions generator** - type-safe bridge from Go types
- [ ] **DevTools integration** - forward console.log, enable network inspection
- [ ] **Hot reload** for web content during development
- [ ] **Production example** - real hybrid app showing best practices

**Implementation**: Phase 3 of roadmap (Q3 2025)

---

### 4. Windows Testing Automation
**Problem**: Can't easily test Windows builds from macOS  
**Impact**: Windows support is untested, might be broken

**The UTM Vision** (from old TODO):

```
┌─────────────────────────────────────────┐
│     macOS Development Machine           │
│                                          │
│  ┌────────────────────────────────┐     │
│  │  UTM (Virtual Machine)         │     │
│  │                                │     │
│  │  ┌──────────────────────────┐  │     │
│  │  │  Windows 11              │  │     │
│  │  │                          │  │     │
│  │  │  - Git (via winget)      │  │     │
│  │  │  - Go (via winget)       │  │     │
│  │  │  - VSCode (via winget)   │  │     │
│  │  │                          │  │     │
│  │  │  → Run goup-util tests   │  │     │
│  │  │  → Build Windows apps    │  │     │
│  │  └──────────────────────────┘  │     │
│  └────────────────────────────────┘     │
└─────────────────────────────────────────┘

Automated via:
1. Packer + UTM plugin creates VM image
2. Go code provisions VM with winget packages
3. CI/CD runs tests in VM automatically
```

**Tasks**:
- [ ] **UTM automation** - Create Windows 11 VM from Go code
  - Use: https://github.com/naveenrajm7/packer-plugin-utm
- [ ] **Winget provisioning** - Install dev tools in VM
  - Use: https://github.com/mbarbita/go-winget
- [ ] **Test runner** - Execute goup-util tests in Windows VM
- [ ] **CI integration** - Run Windows tests on every commit

**Implementation**: Phase 4 (Q4 2025) or when Windows support needed

---

## 🚀 Medium Priority

### 5. Performance Improvements
- [ ] **Incremental builds** - hash-based caching, skip unchanged
- [ ] **Parallel builds** - build multiple platforms concurrently
- [ ] **Parallel icon generation** - 5-10x faster
- [ ] **Docker build cache** - consistent, fast CI/CD builds

**Implementation**: Phase 2 (Q2 2025)

---

### 6. Configuration System
- [ ] **goup.yaml** - project configuration file
- [ ] **Platform-specific settings** - bundle IDs, permissions, signing
- [ ] **Build profiles** - debug, release, staging
- [ ] **CLI overrides** - flags override config file

**Implementation**: Week 2 quick win + ongoing

---

### 7. Testing & Deployment
- [ ] **Simulator/emulator automation** - `goup-util test ios --simulator`
- [ ] **Device deployment** - `goup-util deploy android --device`
- [ ] **Store helpers** - `goup-util deploy appstore --testflight`
- [ ] **CI/CD templates** - GitHub Actions, CircleCI configs

**Implementation**: Phase 4 (Q4 2025)

---

## 🔮 Future (Nice to Have)

### 8. Cross-Compilation Fixes
- [ ] **Linux cross-compile** - Docker-based builds from macOS
- [ ] **Windows cross-compile** - Docker or remote builds
- [ ] **Better CGo detection** - warn early about cross-compile issues

### 9. Plugin System
- [ ] **Custom commands** - extend goup-util via plugins
- [ ] **Build hooks** - pre-build, post-build, pre-deploy
- [ ] **Plugin marketplace** - share community plugins

### 10. Enhanced Examples
- [ ] **Hybrid dashboard** - Gio UI + web charts/graphs
- [ ] **Offline-first app** - IndexedDB + Go backend sync
- [ ] **Camera integration** - Native camera + Go processing
- [ ] **Push notifications** - FCM/APNs integration
- [ ] **OAuth flow** - Authentication with webview

---

## ✅ Completed

- [x] **Core build system** - macOS, iOS, Android working
- [x] **Webviewer example** - Multi-tab browser builds successfully
- [x] **Icon generation** - All platforms supported
- [x] **SDK management** - Caching, auto-install works
- [x] **Workspace support** - Multi-module projects
- [x] **Documentation** - CLAUDE.md, agents/, IMPROVEMENTS.md, WEBVIEW-ANALYSIS.md
- [x] **Testing** - Validated on real platforms
- [x] **Deep analysis** - Know what needs improvement

---

## 📊 Progress Tracking

**Current Phase**: Proof of Concept → Production Ready  
**Next Milestone**: Week 1-4 Quick Wins (Better UX)  
**Long-term Goal**: Best Go hybrid app framework

---

## 🎯 This Week

Focus on **immediate impact**:

1. **Tonight**: Update README (1 hour)
2. **This week**: Better build feedback (2-3 days)
3. **Next week**: Production example app (3-4 days)

Small wins → momentum → adoption → ecosystem

---

## 💡 Ideas Parking Lot

Random ideas to evaluate later:

- **Winget MDM** (https://github.com/jantari/rewinged) - Host winget manifests for internal tools
- **Desktop PWA mode** - Gio app that IS a web browser for PWAs
- **Bridge tooling** - Auto-generate bridge code from OpenAPI specs
- **Visual builder** - GUI for designing Gio layouts
- **Hot reload for Go** - Recompile and restart on code changes
- **Remote build farm** - Build iOS apps without owning a Mac

---

**See also**:
- [docs/IMPROVEMENTS.md](docs/IMPROVEMENTS.md) - Comprehensive improvement roadmap
- [docs/WEBVIEW-ANALYSIS.md](docs/WEBVIEW-ANALYSIS.md) - Cross-platform webview deep dive
- [docs/agents/](docs/agents/) - AI assistant collaboration guides

---

## 🔗 Gio v0.9.1 + gogio Updates (2025-12-20)

New features from the latest Gio, gio-plugins, and gio-cmd updates.

### gogio New Features (gio-cmd PRs merged Dec 2025)

**PR #9: Deep Linking / Custom URI Schemes** ✨ IMPLEMENTED
- [x] `-schemes` flag merged into gogio (Dec 15, 2025)
- [x] Add `--schemes` flag to goup-util build command ✅
- [x] Support for Android, iOS, macOS, Windows deep links ✅
- [x] Example: `goup-util build android --schemes myapp://,https://example.com` ✅
- [ ] Integration with webviewer (pass URLs to webview)

**PR #23: Android App Queries** ✅ IMPLEMENTED
- [x] `-queries` flag merged (Dec 16, 2025)
- [x] Add `--queries` flag to goup-util build command ✅
- [x] Enables checking if apps are installed before launching intents ✅
- [x] Example: `goup-util build android --queries com.google.android.apps.maps` ✅

**PR #21: iOS App Store Compatibility Fixes**
- [x] All 6 validation issues fixed (Dec 16, 2025)
- [x] Bitcode stripping (Asset validation 90482)
- [x] 3-part version format (1.2.3 not 1.2.3.4)
- [x] iPad multitasking disabled by default
- [x] iPad 152x152 icon auto-generated
- [ ] Test App Store upload with goup-util built apps

**PR #22: WASM Go 1.23+ Compatibility**
- [x] Fixed WASM compatibility (Dec 15, 2025)
- [ ] Test `goup-util build wasm` with Go 1.25

**PR #19: Android 15+ / 16KB Page Size**
- [x] 64KB page-size for Android 15+ (May 2025)
- [x] Google Play requires 16KB-compatible by Nov 2025
- [ ] Verify goup-util Android builds are compatible

**PR #20: macOS/iOS Signing & Profiles** ✅ IMPLEMENTED
- [x] Custom profile support merged
- [x] Add `--signkey` flag (keystore, Keychain key, or provisioning profile) ✅

### gio-plugins Updates

**Auth Module (Issue #106)**
- [ ] Test OAuth flows (Apple, Google sign-in)
- [ ] Verify auth callbacks work correctly

### Platform Improvements
- [ ] Test Windows touch screen support in webviewer
- [ ] Test macOS fullscreen MaxSize with webviewer apps
- [ ] Verify Android text rendering fixes

**Current versions (2025-12-20):**
```bash
# Update gio-cmd (gogio) to get new features
go install gioui.org/cmd/gogio@latest

# In your project
go get gioui.org@7bcb315ee174  # v0.9.1-0.20251215212054-7bcb315ee174
go get github.com/gioui-plugins/gio-plugins@v0.9.1
```

---

## 🎯 Screenshot & Documentation Tasks

### Capture App Screenshots
Use Playwright MCP or native tools to create visual documentation:

- [ ] **macOS webviewer** - Running desktop app with tabs/browser
- [ ] **iOS simulator** - App running in iPhone simulator
- [ ] **Android emulator** - App running in Android emulator
- [ ] **All three side-by-side** - Show cross-platform capability

Save to `docs/screenshots/` and link in README.

### Create Complete Hybrid Example

**Problem**: Current webviewer just loads Google.com (external URL)

**Better**: `examples/hybrid-app-complete/` with:

```
hybrid-app-complete/
├── main.go              # Gio UI + embedded HTTP server
├── go.mod
├── icon-source.png
└── web/
    ├── index.html       # Landing page
    ├── app.js           # JavaScript with Go bridge calls
    ├── styles.css       # Styling
    └── assets/
        └── logo.png
```

**Features to demonstrate**:
- ✅ Embedded `//go:embed` web content (no external deps)
- ✅ Local HTTP server on localhost:8080
- ✅ Go ↔ JavaScript bridge (call Go functions from JS)
- ✅ Native Gio UI navigation (tabs, buttons)
- ✅ WebView displaying embedded content
- ✅ Offline-capable (all assets embedded)
- ✅ Works on all platforms (iOS, Android, macOS, Windows)

**This becomes THE showcase example** - proves the vision works end-to-end.

Priority: HIGH (After README update)
