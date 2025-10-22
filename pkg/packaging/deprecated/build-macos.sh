#!/bin/bash
# DEPRECATED: This bash script is replaced by the pure Go 'bundle' command
#
# Migration:
#   Old: bash pkg/packaging/build-macos.sh
#   New: go run . bundle macos examples/hybrid-dashboard
#
# The new bundle command provides:
# - Pure Go implementation (cross-platform)
# - Template-based Info.plist generation
# - Auto-detection of signing certificates
# - Better error messages
# - Consistent CLI interface
#
# This script is kept for reference only.
# See: docs/PACKAGING.md for complete documentation

set -e

echo "⚠️  WARNING: This script is DEPRECATED"
echo "   Use: go run . bundle macos <app-directory>"
echo "   See: docs/PACKAGING.md"
echo ""
sleep 2

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
APP_NAME="goup-util"
BUNDLE_ID="com.joeblew999.goup-util"
VERSION="1.0.0"

# Output directories
BUILD_DIR="$PROJECT_ROOT/.build"
DIST_DIR="$PROJECT_ROOT/.dist"
APP_BUNDLE="$DIST_DIR/$APP_NAME.app"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "🚀 Building goup-util for macOS..."
echo ""

# Clean previous builds
echo "🧹 Cleaning previous builds..."
rm -rf "$BUILD_DIR"
rm -rf "$DIST_DIR"
mkdir -p "$BUILD_DIR"
mkdir -p "$DIST_DIR"

# Build the binary with CGO enabled
echo "🔨 Building binary with CGO (for robotgo screenshot support)..."
cd "$PROJECT_ROOT"
CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o "$BUILD_DIR/$APP_NAME" .

if [ ! -f "$BUILD_DIR/$APP_NAME" ]; then
    echo -e "${RED}❌ Build failed!${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Binary built successfully${NC}"

# Create app bundle structure
echo "📦 Creating app bundle..."
mkdir -p "$APP_BUNDLE/Contents/MacOS"
mkdir -p "$APP_BUNDLE/Contents/Resources"

# Copy binary
cp "$BUILD_DIR/$APP_NAME" "$APP_BUNDLE/Contents/MacOS/$APP_NAME"
chmod +x "$APP_BUNDLE/Contents/MacOS/$APP_NAME"

# Copy Info.plist
cp "$SCRIPT_DIR/macos-info.plist" "$APP_BUNDLE/Contents/Info.plist"

# Update Info.plist with current version
sed -i '' "s/<string>1.0.0<\/string>/<string>$VERSION<\/string>/g" "$APP_BUNDLE/Contents/Info.plist"

echo -e "${GREEN}✓ App bundle created${NC}"

# Check for code signing certificate
echo "🔐 Checking for code signing certificate..."

# Try to find a Developer ID Application certificate
SIGNING_IDENTITY=$(security find-identity -v -p codesigning | grep "Developer ID Application" | head -1 | awk -F'"' '{print $2}')

if [ -z "$SIGNING_IDENTITY" ]; then
    # Fall back to any valid signing identity
    SIGNING_IDENTITY=$(security find-identity -v -p codesigning | grep "Apple Development" | head -1 | awk -F'"' '{print $2}')
fi

if [ -z "$SIGNING_IDENTITY" ]; then
    echo -e "${YELLOW}⚠️  No code signing certificate found${NC}"
    echo -e "${YELLOW}   The app will work but require manual permission granting${NC}"
    echo ""
    echo "To improve the user experience, you can:"
    echo "1. Sign with ad-hoc signature (for local testing):"
    echo "   codesign --force --deep --sign - \"$APP_BUNDLE\""
    echo ""
    echo "2. Get a free Apple Developer ID:"
    echo "   https://developer.apple.com/account/"
    echo ""
    echo "3. Sign with your Apple ID:"
    echo "   codesign --force --deep --sign \"YOUR_APPLE_ID\" \"$APP_BUNDLE\""
    echo ""

    # Ad-hoc signing
    echo "Applying ad-hoc signature..."
    codesign --force --deep --sign - "$APP_BUNDLE"
else
    echo -e "${GREEN}✓ Found signing identity: $SIGNING_IDENTITY${NC}"
    echo "🔏 Signing app bundle..."

    # Sign with entitlements
    codesign --force --deep \
        --sign "$SIGNING_IDENTITY" \
        --entitlements "$SCRIPT_DIR/entitlements.plist" \
        --options runtime \
        "$APP_BUNDLE"

    echo -e "${GREEN}✓ App signed successfully${NC}"
fi

# Verify signature
echo "🔍 Verifying signature..."
codesign --verify --deep --strict "$APP_BUNDLE"
codesign -dv "$APP_BUNDLE"

echo ""
echo -e "${GREEN}✅ Build complete!${NC}"
echo ""
echo "📍 Location: $APP_BUNDLE"
echo ""
echo "🎯 To install and test:"
echo "   1. Open: open \"$APP_BUNDLE\""
echo "   2. Grant Screen Recording permission:"
echo "      System Settings → Privacy & Security → Screen Recording"
echo "   3. Test: \"$APP_BUNDLE/Contents/MacOS/$APP_NAME\" screenshot test.png"
echo ""
echo "📦 To distribute:"
echo "   1. Create DMG: hdiutil create -volname \"$APP_NAME\" -srcfolder \"$DIST_DIR\" -ov -format UDZO \"$DIST_DIR/$APP_NAME.dmg\""
echo "   2. Or ZIP: cd \"$DIST_DIR\" && zip -r \"$APP_NAME.zip\" \"$APP_NAME.app\""
echo ""
