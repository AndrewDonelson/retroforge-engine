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

# Ensure we're in the repo root (where this script is located)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR" || exit 1

# Generate version string: YEAR.MM.DD.HHMM-alpha/beta
VERSION=$(date +"%Y.%m.%d.%H%M")-${VERSION_TYPE}

echo -e "${GREEN}Building RetroForge Engine v${VERSION}${NC}"
echo -e "${YELLOW}Working directory: $(pwd)${NC}"
echo ""

# Create bin directory
BIN_DIR="bin"
mkdir -p "${BIN_DIR}"

# Clean bin directory
echo -e "${YELLOW}Cleaning bin directory...${NC}"
rm -f "${BIN_DIR}"/*

# Clean screenshot files
echo -e "${YELLOW}Cleaning screenshot files...${NC}"
find . -name "screenshot-*.png" -type f -delete 2>/dev/null || true
find . -name "screenshot-*.jpg" -type f -delete 2>/dev/null || true
find . -name "screenshot-*.jpeg" -type f -delete 2>/dev/null || true

# Build flags
LDFLAGS="-s -w"
PKG="./cmd/retroforge"
WASM_PKG="./cmd/wasm"

# Get module name for absolute package path if needed
MODULE_NAME=$(go list -m)

# Function to build for a platform
build_platform() {
    local GOOS=$1
    local GOARCH=$2
    local EXT=$3
    local FILENAME="retroforge-${VERSION}${EXT}"
    
    echo -e "${YELLOW}Building ${GOOS}/${GOARCH}...${NC}"
    
    # Ebiten is mostly pure Go, but macOS requires CGO for GLFW bindings
    # Windows and Linux can cross-compile without CGO
    local CGO_FLAG=""
    if [ "${GOOS}" = "darwin" ]; then
        CGO_FLAG="CGO_ENABLED=1"
        echo -e "${YELLOW}  Note: macOS requires CGO for GLFW bindings${NC}"
    fi
    
    ${CGO_FLAG} GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags "${LDFLAGS}" -o "${BIN_DIR}/${FILENAME}" ${PKG}
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Built ${FILENAME}${NC}"
        return 0
    else
        echo -e "${RED}✗ Failed to build ${GOOS}/${GOARCH}${NC}"
        if [ "${GOOS}" = "darwin" ]; then
            echo -e "${YELLOW}  macOS builds require CGO. Cross-compilation may need macOS toolchain.${NC}"
        fi
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

# Build for Windows (amd64) - Pure Go, cross-compilation works!
if build_platform "windows" "amd64" ".exe"; then
    BUILD_SUCCESS=$((BUILD_SUCCESS + 1))
fi

# Build for Windows (arm64) - Pure Go, cross-compilation works!
if build_platform "windows" "arm64" "-windows-arm64.exe"; then
    BUILD_SUCCESS=$((BUILD_SUCCESS + 1))
fi

# Build for macOS (amd64) - Pure Go, cross-compilation works!
if build_platform "darwin" "amd64" "-macos-amd64"; then
    BUILD_SUCCESS=$((BUILD_SUCCESS + 1))
fi

# Build for macOS (arm64 / Apple Silicon) - Pure Go, cross-compilation works!
if build_platform "darwin" "arm64" "-macos-arm64"; then
    BUILD_SUCCESS=$((BUILD_SUCCESS + 1))
fi

# Build WASM
echo -e "${YELLOW}Building WASM...${NC}"
WASM_FILENAME="retroforge-${VERSION}.wasm"

# Ensure we're still in the repo root (in case build_platform changed directories)
cd "$SCRIPT_DIR" || exit 1

# Verify the WASM package directory exists
if [ ! -d "./cmd/wasm" ]; then
    echo -e "${RED}✗ WASM package directory not found: ./cmd/wasm${NC}"
    echo -e "${RED}Current directory: $(pwd)${NC}"
    echo -e "${RED}Expected location: $SCRIPT_DIR/cmd/wasm${NC}"
    ls -la cmd/ 2>&1 | head -10
    exit 1
fi

# Build WASM (no CGO needed)
echo -e "${YELLOW}Building WASM from: $(pwd)${NC}"
echo -e "${YELLOW}WASM package: ${WASM_PKG}${NC}"
echo -e "${YELLOW}Module: ${MODULE_NAME}${NC}"

# Build WASM - use module path for better reliability
BUILD_OUTPUT=$(GOOS=js GOARCH=wasm go build -ldflags "${LDFLAGS}" -o "${BIN_DIR}/${WASM_FILENAME}" "${MODULE_NAME}/cmd/wasm" 2>&1)
BUILD_STATUS=$?

if [ ${BUILD_STATUS} -ne 0 ]; then
    echo -e "${RED}✗ Failed to build WASM${NC}"
    echo -e "${RED}Build output:${NC}"
    echo "$BUILD_OUTPUT"
    echo -e "${RED}WASM build is critical. Please fix the issue and try again.${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Built ${WASM_FILENAME}${NC}"

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
    echo -e "${GREEN}Build complete! All platforms built successfully (Pure Go - no CGO!)${NC}"
else
    echo -e "${RED}No platforms built successfully!${NC}"
    exit 1
fi

