---
title: "Platform Support"
date: 2025-12-21
draft: false
weight: 2
---

# Platform Support

goup-util supports building native applications for all major platforms from a single codebase.

## Mobile Platforms

### 📱 iOS
- **Output**: Native iOS applications (.app bundles)
- **Distribution**: App Store ready with proper signing
- **Features**:
  - Automatic icon generation for all device sizes
  - Proper Info.plist configuration
  - Asset catalog management
  - Simulator and device testing

```bash
# Build for iOS
goup-util build ios ./my-app

# Generate iOS-specific assets
goup-util icons ios ./my-app
```

### 🤖 Android
- **Output**: Android Package files (.apk)
- **Distribution**: Google Play Store compatible
- **Features**:
  - Multi-density icon generation
  - Proper manifest configuration
  - Resource optimization
  - APK signing support

```bash
# Build for Android
goup-util build android ./my-app

# Generate Android assets
goup-util icons android ./my-app
```

## Desktop Platforms

### 🍎 macOS
- **Output**: Native macOS applications (.app bundles)
- **Distribution**: Mac App Store ready or direct distribution
- **Features**:
  - Proper app bundle structure
  - Code signing integration
  - DMG creation for distribution
  - Native system integration

```bash
# Build for macOS
goup-util build macos ./my-app
```

### 🪟 Windows
- **Output**: Windows executables (.exe) and MSIX packages
- **Distribution**: Microsoft Store compatible
- **Features**:
  - MSIX package creation
  - Windows 10/11 compatibility
  - Proper manifest generation
  - Code signing support

```bash
# Build Windows executable
goup-util build windows ./my-app

# Create MSIX package
goup-util build windows-msix ./my-app
```

### 🐧 Linux
- **Output**: Native Linux binaries
- **Distribution**: AppImage, Flatpak, or traditional packages
- **Features**:
  - Multiple architecture support
  - Desktop integration
  - Package format flexibility

```bash
# Build for Linux
goup-util build linux ./my-app
```

## Web Platform

### 🌐 Progressive Web Apps (PWA)
- **Output**: Modern web applications
- **Distribution**: Web deployment or app store submission
- **Features**:
  - Service worker generation
  - Web app manifest
  - Responsive design
  - Offline functionality

```bash
# Build for web
goup-util build web ./my-app
```

## Cross-Platform Features

### 🎨 Asset Management
All platforms benefit from automatic asset generation:
- **Icons**: Platform-specific sizes and formats
- **Splash Screens**: Proper dimensions for each platform
- **Resources**: Optimized for each target platform

### 📦 Package Management
- **Dependencies**: Automatically managed for each platform
- **SDKs**: Isolated, version-controlled environments
- **Build Tools**: Platform-specific toolchains

### 🔧 Configuration
- **Project-aware**: Understands your project structure
- **Platform-specific**: Customizable per platform
- **Build optimization**: Tailored for each target

## Platform Requirements

| Platform | Host OS | Additional Tools |
|----------|---------|------------------|
| iOS | macOS | Xcode Command Line Tools |
| Android | Any | Android SDK (auto-installed) |
| macOS | macOS | Xcode Command Line Tools |
| Windows | Any | Windows SDK (auto-installed) |
| Linux | Any | Standard build tools |
| Web | Any | Modern web tools |

## Build Matrix

Build for multiple platforms simultaneously:

```bash
# Build for all platforms
goup-util build all ./my-app

# Build for mobile only
goup-util build mobile ./my-app

# Build for desktop only
goup-util build desktop ./my-app

# Custom combinations
goup-util build ios,android,web ./my-app
```
