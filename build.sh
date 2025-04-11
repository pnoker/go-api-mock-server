#!/bin/bash

# Set color output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Create output directory
mkdir -p build

# Get current time as version number
VERSION=$(date +"%Y%m%d%H%M%S")

echo -e "${YELLOW}Building Go API Mock Server v${VERSION}...${NC}"

# Build for Windows
echo -e "${GREEN}Building for Windows...${NC}"
GOOS=windows GOARCH=amd64 go build -ldflags="-X main.Version=${VERSION}" -o build/go-api-mock-server-windows-amd64.exe

# Build for macOS (Intel/AMD)
echo -e "${GREEN}Building for macOS (Intel/AMD)...${NC}"
GOOS=darwin GOARCH=amd64 go build -ldflags="-X main.Version=${VERSION}" -o build/go-api-mock-server-darwin-amd64

# Build for macOS (Apple Silicon)
echo -e "${GREEN}Building for macOS (Apple Silicon)...${NC}"
GOOS=darwin GOARCH=arm64 go build -ldflags="-X main.Version=${VERSION}" -o build/go-api-mock-server-darwin-arm64

# Build for Linux (Intel/AMD)
echo -e "${GREEN}Building for Linux (Intel/AMD)...${NC}"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-X main.Version=${VERSION}" -o build/go-api-mock-server-linux-amd64

# Build for Linux (ARM64)
echo -e "${GREEN}Building for Linux (ARM64)...${NC}"
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-X main.Version=${VERSION}" -o build/go-api-mock-server-linux-arm64

echo -e "${GREEN}Build completed!${NC}"
echo -e "${YELLOW}Executables are in the build directory:${NC}"
echo -e "  - Windows (64-bit): ${GREEN}go-api-mock-server-windows-amd64.exe${NC}"
echo -e "  - macOS (Intel/AMD): ${GREEN}go-api-mock-server-darwin-amd64${NC}"
echo -e "  - macOS (Apple Silicon): ${GREEN}go-api-mock-server-darwin-arm64${NC}"
echo -e "  - Linux (Intel/AMD): ${GREEN}go-api-mock-server-linux-amd64${NC}"
echo -e "  - Linux (ARM64): ${GREEN}go-api-mock-server-linux-arm64${NC}"
echo -e "\n${YELLOW}To run on Linux:${NC}"
echo -e "  chmod +x go-api-mock-server-linux-amd64"
echo -e "  ./go-api-mock-server-linux-amd64" 