#!/bin/bash

# Build Production Script for RetroForge Engine
# Builds binaries for all platforms including WASM
# Version format: YEAR.MM.DD.HHMM-{alpha|beta}

set -e  # Exit on error for critical operations

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
    
    # CGO is required for SDL (native audio). 
    # Cross-compilation to Windows/macOS from Linux may fail.
    # Linux (amd64/arm64) and WASM builds work reliably.
    CGO_ENABLED=1 GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags "${LDFLAGS}" -o "${BIN_DIR}/${FILENAME}" ${PKG}
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Built ${FILENAME}${NC}"
        return 0
    else
        echo -e "${YELLOW}⚠ Failed to build ${GOOS}/${GOARCH}${NC}"
        echo -e "${YELLOW}  Note: SDL/CGO cross-compilation may not work.${NC}"
        echo -e "${YELLOW}  This is expected for Windows/macOS from Linux.${NC}"
        return 1
    fi
}

# Track build success
BUILD_SUCCESS=0

# Build for Linux (amd64) - always succeeds on Linux runner
if build_platform "linux" "amd64" ""; then
    BUILD_SUCCESS=$((BUILD_SUCCESS + 1))
fi

# Build for Linux (arm64) - may work with proper cross-compilation setup
if build_platform "linux" "arm64" "-linux-arm64"; then
    BUILD_SUCCESS=$((BUILD_SUCCESS + 1))
else
    echo -e "${YELLOW}Skipping Linux ARM64 (requires cross-compilation setup)${NC}"
fi

# Build for Windows (amd64) - cross-compilation with CGO is complex
if build_platform "windows" "amd64" ".exe"; then
    BUILD_SUCCESS=$((BUILD_SUCCESS + 1))
else
    echo -e "${YELLOW}Skipping Windows amd64 (CGO cross-compilation not supported)${NC}"
fi

# Build for Windows (arm64) - cross-compilation with CGO is complex
if build_platform "windows" "arm64" "-windows-arm64.exe"; then
    BUILD_SUCCESS=$((BUILD_SUCCESS + 1))
else
    echo -e "${YELLOW}Skipping Windows ARM64 (CGO cross-compilation not supported)${NC}"
fi

# Build for macOS (amd64) - cross-compilation with CGO is complex
if build_platform "darwin" "amd64" "-macos-amd64"; then
    BUILD_SUCCESS=$((BUILD_SUCCESS + 1))
else
    echo -e "${YELLOW}Skipping macOS amd64 (CGO cross-compilation not supported)${NC}"
fi

# Build for macOS (arm64 / Apple Silicon) - cross-compilation with CGO is complex
if build_platform "darwin" "arm64" "-macos-arm64"; then
    BUILD_SUCCESS=$((BUILD_SUCCESS + 1))
else
    echo -e "${YELLOW}Skipping macOS ARM64 (CGO cross-compilation not supported)${NC}"
fi

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
    echo -e "${RED}WASM build is critical. Please fix the issue and try again.${NC}"
    exit 1
fi

# Create version file
echo "${VERSION}" > "${BIN_DIR}/VERSION.txt"

# Summary
echo ""
echo -e "${GREEN}=== Build Summary ===${NC}"
echo "Version: ${VERSION}"
echo "Output directory: ${BIN_DIR}"
echo "Platforms built: ${BUILD_SUCCESS}/6 (+ WASM)"
echo ""
echo "Built binaries:"
ls -lh "${BIN_DIR}"/* 2>/dev/null | awk '{print "  " $9 " (" $5 ")"}'
echo ""
if [ ${BUILD_SUCCESS} -gt 0 ]; then
    echo -e "${GREEN}Build complete! (Some platforms may have been skipped due to CGO limitations)${NC}"
else
    echo -e "${RED}No platforms built successfully!${NC}"
    exit 1
fi

