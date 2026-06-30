#!/bin/bash

# Simple Compiler Script for Havoc Teamserver and Client

# Color codes
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}[*] Starting Havoc compilation...${NC}"

# 1. Resolve Go path
if ! command -v go &> /dev/null; then
    if [ -x "/home/anoam/.local/go/bin/go" ]; then
        export PATH="/home/anoam/.local/go/bin:$PATH"
    elif [ -d "/usr/local/go/bin" ]; then
        export PATH="/usr/local/go/bin:$PATH"
    fi
fi

if ! command -v go &> /dev/null; then
    echo -e "${RED}[-] Error: Go is not installed or not in PATH.${NC}"
    exit 1
fi

# 2. Build Teamserver
echo -e "${BLUE}[*] Compiling teamserver (CGO-free)...${NC}"
cd teamserver
rm -f ../havoc
CGO_ENABLED=0 go build -ldflags="-s -w -X cmd.VersionCommit=$(git rev-parse HEAD 2>/dev/null || echo 'unknown')" -o ../havoc main.go
if [ $? -eq 0 ]; then
    echo -e "${GREEN}[+] Teamserver compiled successfully as ./havoc${NC}"
else
    echo -e "${RED}[-] Failed to compile teamserver${NC}"
    exit 1
fi
cd ..

# 3. Setup privileges for privileged ports (like 443)
echo -e "${BLUE}[*] Setting up privileged-binding capability (requires sudo)...${NC}"
sudo cp havoc /usr/local/bin/havoc && \
sudo chmod 755 /usr/local/bin/havoc && \
sudo setcap 'cap_net_bind_service=+ep' /usr/local/bin/havoc
if [ $? -eq 0 ]; then
    echo -e "${GREEN}[+] Installed to /usr/local/bin/havoc with privileged-binding capability!${NC}"
else
    echo -e "${RED}[-] Warning: Failed to set capabilities. You might need to run the server with sudo if binding < 1024 ports.${NC}"
fi

# 4. Build Client
echo -e "${BLUE}[*] Building client...${NC}"
git submodule update --init --recursive
rm -rf client/Build
mkdir -p client/Build
cd client/Build
cmake ..
if [ $? -eq 0 ]; then
    cmake --build . -- -j2
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}[+] Client compiled successfully!${NC}"
    else
        echo -e "${RED}[-] Failed to compile client${NC}"
        exit 1
    fi
else
    echo -e "${RED}[-] Failed to configure client build using CMake${NC}"
    exit 1
fi

echo -e "${GREEN}[+] Havoc build process complete!${NC}"
