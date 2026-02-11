---
title: "User Guides"
date: 2025-12-21
draft: false
weight: 1
---

# User Guides

Everything you need to build, package, and distribute Gio applications with goup-util.

## Getting Started

1. **[Quick Start](/users/quickstart/)** -- Install goup-util and build your first app in 5 minutes
2. **[Platform Support](/users/platforms/)** -- What works on each platform, requirements, and known limitations

## Building and Distributing

3. **[Packaging](/users/packaging/)** -- The three-tier system: Build, Bundle, Package
4. **[Webviewer Shell](/users/webviewer-shell/)** -- Ship any website as a native desktop app with zero coding

## Command Reference

```bash
# Build for a platform
goup-util build <platform> <app-directory>

# Build and run immediately
goup-util run <platform> <app-directory>

# Create signed bundle for distribution
goup-util bundle <platform> <app-directory>

# Package into archive (tar.gz / zip)
goup-util package <platform> <app-directory>

# Generate platform icons from source image
goup-util icons <app-directory>

# Install platform SDKs
goup-util install <sdk-name>

# List available SDKs
goup-util list

# Full help
goup-util --help
```

**Supported platforms:** `macos`, `ios`, `ios-simulator`, `android`, `windows`
