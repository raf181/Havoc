#!/bin/bash

# GoITDemon Build Script
# Builds the IT-focused demon for multiple platforms

set -e

VERSION="1.0.0"
BUILD_DIR="bin"
AGENT_NAME="GoITDemon"

# Create build directory
mkdir -p $BUILD_DIR

echo "[+] Building $AGENT_NAME v$VERSION for multiple platforms..."

# Set build flags
BUILD_FLAGS="-ldflags=-s -ldflags=-w"

# Windows builds
echo "[+] Building for Windows..."
GOOS=windows GOARCH=amd64 go build $BUILD_FLAGS -o $BUILD_DIR/${AGENT_NAME}_windows_x64.exe .
GOOS=windows GOARCH=386 go build $BUILD_FLAGS -o $BUILD_DIR/${AGENT_NAME}_windows_x86.exe .

# Linux builds
echo "[+] Building for Linux..."
GOOS=linux GOARCH=amd64 go build $BUILD_FLAGS -o $BUILD_DIR/${AGENT_NAME}_linux_x64 .
GOOS=linux GOARCH=386 go build $BUILD_FLAGS -o $BUILD_DIR/${AGENT_NAME}_linux_x86 .
GOOS=linux GOARCH=arm64 go build $BUILD_FLAGS -o $BUILD_DIR/${AGENT_NAME}_linux_arm64 .

# macOS builds (for testing)
echo "[+] Building for macOS..."
GOOS=darwin GOARCH=amd64 go build $BUILD_FLAGS -o $BUILD_DIR/${AGENT_NAME}_macos_x64 .
GOOS=darwin GOARCH=arm64 go build $BUILD_FLAGS -o $BUILD_DIR/${AGENT_NAME}_macos_arm64 .

echo "[+] Build complete! Binaries are in the $BUILD_DIR directory:"
ls -la $BUILD_DIR/

echo ""
echo "[+] Usage examples:"
echo "  Windows: .\\${AGENT_NAME}_windows_x64.exe"
echo "  Linux:   ./${AGENT_NAME}_linux_x64"
echo "  macOS:   ./${AGENT_NAME}_macos_x64"
echo ""
echo "[+] The agent will attempt to connect to 127.0.0.1:40056 by default"
echo "    Make sure your Havoc teamserver is configured to accept this agent type"