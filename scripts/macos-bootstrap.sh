#!/bin/bash
# macOS Bootstrap Script for goup-util
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/joeblew999/goup-util/main/scripts/macos-bootstrap.sh | bash
#
# Or download and run:
#   chmod +x macos-bootstrap.sh
#   ./macos-bootstrap.sh

set -e

echo "============================================"
echo " goup-util macOS Bootstrap"
echo "============================================"
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
GRAY='\033[0;90m'
NC='\033[0m' # No Color

# Detect architecture
ARCH=$(uname -m)
if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" = "arm64" ]; then
    ARCH="arm64"
else
    echo -e "${RED}ERROR: Unsupported architecture: $ARCH${NC}"
    exit 1
fi

BINARY_NAME="goup-util-darwin-$ARCH"
echo -e "${GRAY}Detected platform: darwin-$ARCH${NC}"

# Check prerequisites
echo -e "${GREEN}[1/4] Checking prerequisites...${NC}"

# Check if Homebrew is installed
if ! command -v brew &> /dev/null; then
    echo -e "${YELLOW}  Homebrew not found. Installing Homebrew...${NC}"
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
else
    echo -e "${GRAY}  ✓ Homebrew installed${NC}"
fi

# Check if Git is installed
if ! command -v git &> /dev/null; then
    echo -e "${YELLOW}  Git not found. Installing Git...${NC}"
    brew install git
else
    echo -e "${GRAY}  ✓ Git installed: $(git --version)${NC}"
fi

# Check if Go is installed
echo ""
echo -e "${GREEN}[2/4] Installing Go...${NC}"
if ! command -v go &> /dev/null; then
    echo -e "${GRAY}  Installing Go via Homebrew...${NC}"
    brew install go
else
    echo -e "${GRAY}  ✓ Go installed: $(go version)${NC}"
fi

# Check if Task is installed
echo ""
echo -e "${GREEN}[3/4] Installing Task (Taskfile)...${NC}"
if ! command -v task &> /dev/null; then
    echo -e "${GRAY}  Installing Task via Homebrew...${NC}"
    brew install go-task/tap/go-task
else
    echo -e "${GRAY}  ✓ Task installed: $(task --version)${NC}"
fi

# Download and install goup-util
echo ""
echo -e "${GREEN}[4/4] Installing goup-util...${NC}"

# Create temporary directory
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

echo -e "${GRAY}  Fetching latest release...${NC}"

# Get latest release info from GitHub
RELEASE_URL="https://api.github.com/repos/joeblew999/goup-util/releases/latest"
RELEASE_JSON=$(curl -fsSL "$RELEASE_URL")

# Extract download URL for the binary
DOWNLOAD_URL=$(echo "$RELEASE_JSON" | grep -o "\"browser_download_url\": \"[^\"]*$BINARY_NAME\"" | sed 's/"browser_download_url": "//' | tr -d '"')

if [ -z "$DOWNLOAD_URL" ]; then
    echo -e "${RED}  ERROR: Binary not found for darwin-$ARCH${NC}"
    echo -e "${YELLOW}  Available assets:${NC}"
    echo "$RELEASE_JSON" | grep -o "\"name\": \"goup-util-[^\"]*\"" | sed 's/"name": "/    - /' | tr -d '"'
    echo ""
    echo -e "${YELLOW}  Fallback: Clone and build from source:${NC}"
    echo -e "${GRAY}    git clone https://github.com/joeblew999/goup-util.git${NC}"
    echo -e "${GRAY}    cd goup-util${NC}"
    echo -e "${GRAY}    go build .${NC}"
    rm -rf "$TMP_DIR"
    exit 1
fi

echo -e "${GRAY}  Downloading $BINARY_NAME...${NC}"
curl -fsSL -o goup-util "$DOWNLOAD_URL"

# Make executable
chmod +x goup-util

# Install to /usr/local/bin
INSTALL_PATH="/usr/local/bin/goup-util"

if [ -w "/usr/local/bin" ]; then
    mv goup-util "$INSTALL_PATH"
else
    echo -e "${YELLOW}  Need sudo to install to /usr/local/bin${NC}"
    sudo mv goup-util "$INSTALL_PATH"
fi

# Verify installation
if [ -f "$INSTALL_PATH" ]; then
    FILE_SIZE=$(du -h "$INSTALL_PATH" | cut -f1)
    echo -e "${GRAY}  Downloaded: $FILE_SIZE${NC}"
    echo -e "${GREEN}  ✓ Installed to: $INSTALL_PATH${NC}"

    # Handle macOS Gatekeeper (unsigned binary)
    echo ""
    echo -e "${YELLOW}⚠️  macOS Security Notice:${NC}"
    echo -e "${GRAY}  The binary is not code-signed yet (no signing certificate).${NC}"
    echo -e "${GRAY}  Removing quarantine attribute...${NC}"

    # Remove quarantine attribute
    xattr -d com.apple.quarantine "$INSTALL_PATH" 2>/dev/null || true

    # Test the binary
    echo ""
    echo -e "${GREEN}  Testing goup-util...${NC}"

    if "$INSTALL_PATH" --help > /dev/null 2>&1; then
        "$INSTALL_PATH" --help | head -5
        echo ""
        echo -e "${GREEN}  ✓ goup-util successfully installed and working!${NC}"
    else
        echo -e "${YELLOW}  If you see a security warning:${NC}"
        echo -e "${GRAY}    1. Open System Settings → Privacy & Security${NC}"
        echo -e "${GRAY}    2. Look for blocked app message${NC}"
        echo -e "${GRAY}    3. Click 'Allow Anyway'${NC}"
        echo -e "${GRAY}    4. Run: goup-util --help${NC}"
        echo ""
        echo -e "${GRAY}  Or bypass Gatekeeper for this binary:${NC}"
        echo -e "${GRAY}    sudo spctl --add $INSTALL_PATH${NC}"
        echo -e "${GRAY}    xattr -d com.apple.quarantine $INSTALL_PATH${NC}"
        echo ""
        echo -e "${YELLOW}  Note: Future releases will be code-signed.${NC}"
    fi
else
    echo -e "${RED}  ERROR: Installation failed${NC}"
    rm -rf "$TMP_DIR"
    exit 1
fi

# Cleanup
cd ~
rm -rf "$TMP_DIR"

echo ""
echo "============================================"
echo -e " ${GREEN}Setup Complete!${NC}"
echo "============================================"
echo ""
echo -e "${CYAN}Next steps:${NC}"
echo "  1. Test goup-util:"
echo "     goup-util --help"
echo ""
echo "  2. Clone example projects:"
echo "     git clone https://github.com/joeblew999/goup-util.git"
echo "     cd goup-util"
echo ""
echo "  3. Build an example:"
echo "     goup-util build macos examples/hybrid-dashboard"
echo "     goup-util bundle macos examples/hybrid-dashboard"
echo ""
echo "  4. Update goup-util anytime:"
echo "     goup-util self upgrade"
echo ""
echo "  5. Install Android SDK/NDK (for mobile builds):"
echo "     goup-util install android-sdk"
echo "     goup-util install android-ndk"
echo ""
echo "For more info: https://github.com/joeblew999/goup-util"
echo ""
