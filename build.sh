#!/bin/bash

# Simple Compiler Script for Havoc Teamserver and Client

# Color codes
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color
AUTO_INSTALL_DEPS=0

for arg in "$@"; do
    case "$arg" in
        --auto-install-deps|--install-deps|-y|--yes)
            AUTO_INSTALL_DEPS=1
            ;;
    esac
done

echo -e "${BLUE}[*] Starting Havoc compilation...${NC}"

if ! command -v sudo &> /dev/null; then
    echo -e "${RED}[-] Error: sudo is required to install missing dependencies and set capabilities.${NC}"
    exit 1
fi

if ! command -v apt-get &> /dev/null; then
    echo -e "${RED}[-] Error: apt-get is required for automatic dependency installation on this system.${NC}"
    echo -e "${RED}    Install the missing build dependencies manually, then rerun this script.${NC}"
    exit 1
fi

MISSING_COMMANDS=()
for command_name in git go cmake nasm x86_64-w64-mingw32-gcc x86_64-w64-mingw32-g++ i686-w64-mingw32-gcc i686-w64-mingw32-g++; do
    if ! command -v "$command_name" &> /dev/null; then
        MISSING_COMMANDS+=("$command_name")
    fi
done

DEBIAN_PACKAGES=(
    build-essential
    python3-dev
    qtbase5-dev
    libqt5websockets5-dev
    libcap2-bin
    nasm
    gcc-mingw-w64-x86-64
    g++-mingw-w64-x86-64
    gcc-mingw-w64-i686
    g++-mingw-w64-i686
)

MISSING_PACKAGES=()
for package in "${DEBIAN_PACKAGES[@]}"; do
    if ! dpkg -s "$package" &> /dev/null; then
        MISSING_PACKAGES+=("$package")
    fi
done

if [ ${#MISSING_COMMANDS[@]} -gt 0 ] || [ ${#MISSING_PACKAGES[@]} -gt 0 ]; then
    if [ ${#MISSING_COMMANDS[@]} -gt 0 ]; then
        echo -e "${BLUE}[*] Missing commands detected:${NC} ${MISSING_COMMANDS[*]}"
    fi
    if [ ${#MISSING_PACKAGES[@]} -gt 0 ]; then
        echo -e "${BLUE}[*] Missing system packages detected:${NC} ${MISSING_PACKAGES[*]}"
    fi
    if [ "$AUTO_INSTALL_DEPS" -eq 0 ]; then
        read -r -p "Install missing dependencies now? [y/N] " install_missing_deps
        case "$install_missing_deps" in
            [yY]|[yY][eE][sS])
                AUTO_INSTALL_DEPS=1
                ;;
            *)
                echo -e "${RED}[-] Dependency installation skipped. Install the listed packages and rerun the script.${NC}"
                exit 1
                ;;
        esac
    fi

    if [ "$AUTO_INSTALL_DEPS" -eq 1 ]; then
        echo -e "${BLUE}[*] Installing missing dependencies with apt-get...${NC}"
        install_targets=("${MISSING_PACKAGES[@]}")
        for command_name in "${MISSING_COMMANDS[@]}"; do
            case "$command_name" in
                git)
                    install_targets+=(git)
                    ;;
                go)
                    install_targets+=(golang-go)
                    ;;
                cmake)
                    install_targets+=(cmake)
                    ;;
                nasm)
                    install_targets+=(nasm)
                    ;;
                x86_64-w64-mingw32-gcc|x86_64-w64-mingw32-g++)
                    install_targets+=(gcc-mingw-w64-x86-64 g++-mingw-w64-x86-64)
                    ;;
                i686-w64-mingw32-gcc|i686-w64-mingw32-g++)
                    install_targets+=(gcc-mingw-w64-i686 g++-mingw-w64-i686)
                    ;;
            esac
        done

        sudo apt-get update && sudo apt-get install -y "${install_targets[@]}"
        if [ $? -ne 0 ]; then
            echo -e "${RED}[-] Failed to install missing dependencies.${NC}"
            exit 1
        fi
    fi
fi

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
read -r -p "Install havoc to /usr/local/bin and set privileged-binding capability? [y/N] " install_havoc
case "$install_havoc" in
    [yY]|[yY][eE][sS])
        echo -e "${BLUE}[*] Setting up privileged-binding capability (requires sudo)...${NC}"
        sudo cp havoc /usr/local/bin/havoc && \
        sudo chmod 755 /usr/local/bin/havoc && \
        sudo setcap 'cap_net_bind_service=+ep' /usr/local/bin/havoc
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}[+] Installed to /usr/local/bin/havoc with privileged-binding capability!${NC}"
        else
            echo -e "${RED}[-] Warning: Failed to set capabilities. You might need to run the server with sudo if binding < 1024 ports.${NC}"
        fi
        ;;
    *)
        echo -e "${BLUE}[*] Skipping privileged install; the local build will remain in the workspace.${NC}"
        ;;
esac

# 4. Build Client
echo -e "${BLUE}[*] Building client...${NC}"
git submodule update --init --recursive
if [ -d "client/Modules/.git" ]; then
    git -C client/Modules pull --ff-only
else
    rm -rf client/Modules
    git clone https://github.com/HavocFramework/Modules client/Modules
fi
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
