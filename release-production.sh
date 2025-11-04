#!/bin/bash

# Production Release Script for RetroForge
# This script:
# 1. Builds all platform binaries
# 2. Updates webapp WASM files
# 3. Provides instructions for GitHub release

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Get version type (alpha or beta)
VERSION_TYPE=${1:-alpha}

if [[ "$VERSION_TYPE" != "alpha" && "$VERSION_TYPE" != "beta" ]]; then
    echo -e "${RED}Error: Version type must be 'alpha' or 'beta'${NC}"
    exit 1
fi

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}RetroForge Production Release${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Step 1: Clean screenshot files
echo -e "${YELLOW}Step 1: Cleaning screenshot files...${NC}"
find . -name "screenshot-*.png" -type f -delete 2>/dev/null || true
find . -name "screenshot-*.jpg" -type f -delete 2>/dev/null || true
find . -name "screenshot-*.jpeg" -type f -delete 2>/dev/null || true
echo -e "${GREEN}✓ Screenshots cleaned${NC}"
echo ""

# Step 2: Build all platforms
echo -e "${YELLOW}Step 2: Building all platforms...${NC}"
./build-production.sh "$VERSION_TYPE"

if [ $? -ne 0 ]; then
    echo -e "${RED}Build failed! Aborting release.${NC}"
    exit 1
fi

# Get version from VERSION.txt
VERSION=$(cat bin/VERSION.txt)
TAG="v${VERSION}"

echo ""
echo -e "${GREEN}✓ Build completed successfully${NC}"
echo ""

# Step 3: Update webapp WASM files
echo -e "${YELLOW}Step 3: Updating webapp WASM files...${NC}"

WEBAPP_DIR="../retroforge-webapp"
WASM_FILE=$(ls bin/retroforge-*.wasm | head -n1)

if [ ! -d "$WEBAPP_DIR" ]; then
    echo -e "${YELLOW}Warning: Webapp directory not found at $WEBAPP_DIR${NC}"
    echo -e "${YELLOW}Skipping WASM update.${NC}"
else
    if [ -f "$WASM_FILE" ]; then
        # Create engine directory if it doesn't exist
        mkdir -p "$WEBAPP_DIR/public/engine"
        
        # Copy WASM file
        cp "$WASM_FILE" "$WEBAPP_DIR/public/engine/retroforge.wasm"
        echo -e "${GREEN}✓ Copied WASM to webapp${NC}"
        
        # Copy wasm_exec.js if it exists
        if [ -f "bin/wasm_exec.js" ]; then
            cp "bin/wasm_exec.js" "$WEBAPP_DIR/public/engine/wasm_exec.js"
            echo -e "${GREEN}✓ Copied wasm_exec.js to webapp${NC}"
        fi
    else
        echo -e "${RED}Error: WASM file not found${NC}"
    fi
fi

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}Release Summary${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo -e "Version: ${GREEN}${VERSION}${NC}"
echo -e "Tag: ${GREEN}${TAG}${NC}"
echo -e "Type: ${GREEN}${VERSION_TYPE}${NC}"
echo ""
echo -e "Binaries built:"
ls -lh bin/retroforge-* bin/wasm_exec.js 2>/dev/null | awk '{print "  " $9 " (" $5 ")"}'
echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${YELLOW}Next Steps:${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo "1. Review the built binaries in: bin/"
echo ""
echo "2. Create GitHub Release:"
echo "   Option A - Automated (via GitHub Actions):"
echo "     - Go to: Actions > Build and Release"
echo "     - Click 'Run workflow'"
echo "     - Select version type: ${VERSION_TYPE}"
echo "     - Click 'Run workflow'"
echo ""
echo "   Option B - Manual:"
echo "     - Go to: https://github.com/[repo]/releases/new"
echo "     - Tag: ${TAG}"
echo "     - Title: RetroForge Engine ${VERSION}"
echo "     - Upload all files from bin/ directory"
echo "     - Mark as ${VERSION_TYPE} (prerelease)"
echo ""
echo "3. Commit webapp changes (if WASM was updated):"
if [ -d "$WEBAPP_DIR" ]; then
    echo "   cd $WEBAPP_DIR"
    echo "   git add public/engine/"
    echo "   git commit -m 'Update WASM to ${VERSION}'"
    echo "   git push"
fi
echo ""
echo -e "${GREEN}Release preparation complete!${NC}"

