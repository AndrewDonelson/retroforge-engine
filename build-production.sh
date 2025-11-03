#!/bin/bash

# Build Production Script for RetroForge Engine
# Builds binaries for all platforms including WASM
# Version format: YEAR.MM.DD.HHMM-{alpha|beta}

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Determine version type (alpha or beta) - default to alpha
VERSION_TYPE=${1:-alpha}

# Generate version string: YEAR.MM.DD.HHMM-alpha/beta
VERSION=$(date +"%Y.%m.%d.%H%M")-${VERSION_TYPE}

echo -e "${GREEN}Building RetroForge Engine v${VERSION}${NC}"
echo ""

# Create bin directory
BIN_DIR="bin"
mkdir -p "${BIN_DIR}"

# Clean bin directory
echo -e "${YELLOW}Cleaning bin directory...${NC}"
rm -f "${BIN_DIR}"/*

# Build flags
LDFLAGS="-s -w"
PKG="./cmd/retroforge"
WASM_PKG="./cmd/wasm"

# Function to build for a platform
build_platform() {
    local GOOS=$1
    local GOARCH=$2
    local EXT=$3
    local FILENAME="retroforge-${VERSION}${EXT}"
    
    echo -e "${YELLOW}Building ${GOOS}/${GOARCH}...${NC}"
    GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags "${LDFLAGS}" -o "${BIN_DIR}/${FILENAME}" ${PKG}
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Built ${FILENAME}${NC}"
    else
        echo -e "${RED}✗ Failed to build ${GOOS}/${GOARCH}${NC}"
        return 1
    fi
}

# Build for Linux (amd64)
build_platform "linux" "amd64" ""

# Build for Linux (arm64)
build_platform "linux" "arm64" "-linux-arm64"

# Build for Windows (amd64)
build_platform "windows" "amd64" ".exe"

# Build for Windows (arm64)
build_platform "windows" "arm64" "-windows-arm64.exe"

# Build for macOS (amd64)
build_platform "darwin" "amd64" "-macos-amd64"

# Build for macOS (arm64 / Apple Silicon)
build_platform "darwin" "arm64" "-macos-arm64"

# Build WASM
echo -e "${YELLOW}Building WASM...${NC}"
WASM_FILENAME="retroforge-${VERSION}.wasm"
GOOS=js GOARCH=wasm go build -ldflags "${LDFLAGS}" -o "${BIN_DIR}/${WASM_FILENAME}" ${WASM_PKG}

# Copy wasm_exec.js if available
GOROOT=$(go env GOROOT)
if [ -f "${GOROOT}/misc/wasm/wasm_exec.js" ]; then
    cp "${GOROOT}/misc/wasm/wasm_exec.js" "${BIN_DIR}/wasm_exec.js"
elif [ -f "${GOROOT}/lib/wasm/wasm_exec.js" ]; then
    cp "${GOROOT}/lib/wasm/wasm_exec.js" "${BIN_DIR}/wasm_exec.js"
else
    echo -e "${YELLOW}Warning: wasm_exec.js not found${NC}"
fi

if [ -f "${BIN_DIR}/${WASM_FILENAME}" ]; then
    echo -e "${GREEN}✓ Built ${WASM_FILENAME}${NC}"
else
    echo -e "${RED}✗ Failed to build WASM${NC}"
    exit 1
fi

# Create version file
echo "${VERSION}" > "${BIN_DIR}/VERSION.txt"

# Summary
echo ""
echo -e "${GREEN}=== Build Summary ===${NC}"
echo "Version: ${VERSION}"
echo "Output directory: ${BIN_DIR}"
echo ""
echo "Built binaries:"
ls -lh "${BIN_DIR}"/* | awk '{print "  " $9 " (" $5 ")"}'
echo ""
echo -e "${GREEN}Build complete!${NC}"

