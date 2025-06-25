# Quick Start Guide

Get up and running with goup-util in under 5 minutes.

## Prerequisites

- [Go](https://golang.org/) (1.18+)
- [Task](https://taskfile.dev/) (build automation)

## Installation

```bash
# Clone and build
git clone <repository-url>
cd goup-util
task build
```

## Setup Your Development Environment

### Option 1: Setup Everything
```bash
# Installs all SDKs and tools for every platform
task setup:all
```

### Option 2: Platform-Specific Setup
```bash
# Android development
task setup:android

# iOS development (macOS only)
task setup:ios

# Windows development
task setup:windows
```

## Build Your First App

### Using the Example App
```bash
# Build for your current platform
task build:example:macos     # macOS
task build:example:android   # Android
task build:example:ios       # iOS
task build:example:windows   # Windows

# Build for all platforms
task build:example:all
```

### Using Your Own App
```bash
# Create a new project
mkdir my-awesome-app
cd my-awesome-app

# Initialize with Go modules
go mod init my-awesome-app

# Build for specific platforms
goup-util build android ./
goup-util build ios ./
goup-util build windows ./

# Or build for everything
goup-util build all ./
```

## Next Steps

- [Platform-Specific Features](platforms.md)
- [Asset Management](assets.md)
- [CI/CD Integration](cicd.md)
- [Advanced Configuration](advanced.md)
