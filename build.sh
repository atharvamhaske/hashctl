#!/bin/bash
# Build script for hashctl

VERSION=${1:-v0.1.0}
BUILD_DATE=$(date -u +%Y-%m-%d)

echo "Building hashctl $VERSION..."

# Linux AMD64
echo "Building Linux AMD64..."
GOOS=linux GOARCH=amd64 go build -ldflags "-X github.com/atharvamhaske/hashctl/cmd.Version=$VERSION -X github.com/atharvamhaske/hashctl/cmd.BuildDate=$BUILD_DATE" -o hashctl-linux-amd64 .

# Linux ARM64
echo "Building Linux ARM64..."
GOOS=linux GOARCH=arm64 go build -ldflags "-X github.com/atharvamhaske/hashctl/cmd.Version=$VERSION -X github.com/atharvamhaske/hashctl/cmd.BuildDate=$BUILD_DATE" -o hashctl-linux-arm64 .

# macOS AMD64
echo "Building macOS AMD64..."
GOOS=darwin GOARCH=amd64 go build -ldflags "-X github.com/atharvamhaske/hashctl/cmd.Version=$VERSION -X github.com/atharvamhaske/hashctl/cmd.BuildDate=$BUILD_DATE" -o hashctl-darwin-amd64 .

# macOS ARM64
echo "Building macOS ARM64..."
GOOS=darwin GOARCH=arm64 go build -ldflags "-X github.com/atharvamhaske/hashctl/cmd.Version=$VERSION -X github.com/atharvamhaske/hashctl/cmd.BuildDate=$BUILD_DATE" -o hashctl-darwin-arm64 .

# Windows AMD64
echo "Building Windows AMD64..."
GOOS=windows GOARCH=amd64 go build -ldflags "-X github.com/atharvamhaske/hashctl/cmd.Version=$VERSION -X github.com/atharvamhaske/hashctl/cmd.BuildDate=$BUILD_DATE" -o hashctl-windows-amd64.exe .

# Windows ARM64
echo "Building Windows ARM64..."
GOOS=windows GOARCH=arm64 go build -ldflags "-X github.com/atharvamhaske/hashctl/cmd.Version=$VERSION -X github.com/atharvamhaske/hashctl/cmd.BuildDate=$BUILD_DATE" -o hashctl-windows-arm64.exe .

# Generate checksums
echo "Generating checksums..."
sha256sum hashctl-* > checksums.txt

echo "Done! Binaries built:"
ls -lh hashctl-*


