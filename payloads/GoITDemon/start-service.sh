#!/bin/bash

# GoITDemon Service Agent Startup Script

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SERVICE_DIR="$SCRIPT_DIR/service"
BUILD_DIR="$SCRIPT_DIR/bin"

# Default configuration
TEAMSERVER_HOST="${TEAMSERVER_HOST:-127.0.0.1}"
TEAMSERVER_PORT="${TEAMSERVER_PORT:-40056}"
SERVICE_PASSWORD="${SERVICE_PASSWORD:-service-password}"
SERVICE_ENDPOINT="${SERVICE_ENDPOINT:-service-endpoint}"

echo "[+] Starting GoITDemon Service Agent"
echo "[+] Teamserver: $TEAMSERVER_HOST:$TEAMSERVER_PORT"

# Build service if not exists
if [ ! -f "$BUILD_DIR/GoITDemon-Service" ]; then
    echo "[+] Building service agent..."
    cd "$SCRIPT_DIR"
    make service
    if [ $? -ne 0 ]; then
        echo "[-] Failed to build service agent"
        exit 1
    fi
fi

# Start the service
echo "[+] Starting service agent..."
cd "$BUILD_DIR"
./GoITDemon-Service \
    -host "$TEAMSERVER_HOST" \
    -port "$TEAMSERVER_PORT" \
    -password "$SERVICE_PASSWORD" \
    -endpoint "$SERVICE_ENDPOINT"