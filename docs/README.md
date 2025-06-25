# goup-util: Cross-Platform Development Made Simple

**One tool. Every platform. Zero hassle.**

goup-util is a developer-focused CLI that eliminates the complexity of cross-platform mobile and desktop development. Build for iOS, Android, macOS, Windows, and web from a single codebase with automated SDK management and asset generation.

## Why goup-util?

### 🎯 One Command, All Platforms
```bash
# Build for every platform in seconds
goup-util build all ./my-app
```

### 📱 Mobile-First Design
- **iOS**: Native apps with proper signing and assets
- **Android**: APKs ready for Play Store distribution
- **Automatic icon generation** for all screen densities

### 🖥️ Desktop Ready
- **macOS**: Native .app bundles with proper metadata
- **Windows**: MSIX packages for Microsoft Store
- **Linux**: Optimized binaries

### 🌐 Web Compatible
- Progressive Web Apps (PWA) support
- Modern web standards compliance
- Responsive design automation

## Core Benefits

### ✅ Zero Configuration SDK Management
No more wrestling with Android Studio, Xcode versions, or missing dependencies. goup-util automatically downloads and manages all required SDKs in isolated environments.

```bash
# Sets up complete Android development environment
goup-util setup android

# Installs iOS development tools  
goup-util setup ios
```

### ✅ Automated Asset Pipeline
Icons, splash screens, and platform-specific assets are generated automatically from your source files.

```bash
# Generates all icon sizes for all platforms
goup-util icons all ./my-app
```

### ✅ Reproducible Builds
Every build is identical across machines. Perfect for CI/CD pipelines and team development.

### ✅ Project-Aware Architecture
All commands understand your project structure and handle paths, dependencies, and outputs intelligently.

## Perfect For

### 🚀 Startups & Indie Developers
- Launch on all platforms simultaneously
- Minimize time-to-market
- Focus on features, not toolchain complexity

### 🏢 Enterprise Teams
- Standardized development environments
- Consistent CI/CD across platforms
- Reduced onboarding time for new developers

### 🔬 Rapid Prototyping
- Test ideas across platforms instantly
- Quick feedback loops
- Easy platform-specific customization

## How It Works

### 1. Write Once
Create your application using modern development practices with a single codebase.

### 2. Build Everywhere
```bash
goup-util build ios ./my-app        # iOS app
goup-util build android ./my-app    # Android APK
goup-util build windows ./my-app    # Windows MSIX
goup-util build web ./my-app        # Progressive Web App
```

### 3. Ship Fast
All outputs are production-ready with proper signing, metadata, and platform conventions.

## Key Features

- **🔧 SDK Management**: Automatic download and setup of platform SDKs
- **🎨 Asset Generation**: Icons, splash screens, and resources for all platforms
- **📦 Package Creation**: Store-ready packages (APK, MSIX, DMG, etc.)
- **🔄 CI/CD Ready**: Perfect for automated build pipelines
- **📱 Mobile Optimized**: Handles signing, provisioning, and store requirements
- **🖥️ Desktop Native**: Platform-specific features and integrations
- **🌐 Web Standards**: Modern PWA capabilities

## Getting Started

1. **Install Prerequisites**
   ```bash
   # Install Go, Task, and platform tools
   go install gioui.org/cmd/gogio@latest
   ```

2. **Setup Development Environment**
   ```bash
   # One command sets up everything
   goup-util setup all
   ```

3. **Build Your First App**
   ```bash
   # Create and build for all platforms
   goup-util build all ./my-first-app
   ```

## Architecture

goup-util is designed around three core principles:

1. **Project-Aware**: Every command understands your project structure
2. **Idempotent**: Running commands multiple times produces identical results
3. **Service-Ready**: Built for both CLI usage and future API integration

This makes it perfect for both interactive development and automated pipelines.

---

*Ready to build everywhere? Get started with goup-util today.*
