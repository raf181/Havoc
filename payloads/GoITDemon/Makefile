# Makefile for GoITDemon

AGENT_NAME = GoITDemon
VERSION = 1.0.0
BUILD_DIR = bin
LDFLAGS = -ldflags '-s -w -X main.AGENT_VERSION=$(VERSION)'

.PHONY: all clean windows linux macos test deps service service-deps

all: windows linux macos service

# Install dependencies
deps:
	go mod tidy
	go mod download

# Install service dependencies
service-deps:
	cd service && go mod tidy && go mod download

# Build service agent
service: service-deps
	@echo "Building GoITDemon service agent..."
	@mkdir -p $(BUILD_DIR)
	cd service && go build -o ../$(BUILD_DIR)/GoITDemon-Service .

# Build for Windows
windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(AGENT_NAME)_windows_x64.exe .
	GOOS=windows GOARCH=386 go build $(LDFLAGS) -o $(BUILD_DIR)/$(AGENT_NAME)_windows_x86.exe .

# Build for Linux
linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(AGENT_NAME)_linux_x64 .
	GOOS=linux GOARCH=386 go build $(LDFLAGS) -o $(BUILD_DIR)/$(AGENT_NAME)_linux_x86 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(AGENT_NAME)_linux_arm64 .

# Build for macOS
macos:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(AGENT_NAME)_macos_x64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(AGENT_NAME)_macos_arm64 .

# Build current platform only
build:
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(AGENT_NAME) .

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	go clean

# Development build with debug info
debug:
	@mkdir -p $(BUILD_DIR)
	go build -gcflags="all=-N -l" -o $(BUILD_DIR)/$(AGENT_NAME)_debug .

# Run the agent locally for testing
run:
	go run .

# Format code
fmt:
	go fmt ./...

# Check for issues
vet:
	go vet ./...

# Install tools
tools:
	go install golang.org/x/tools/cmd/goimports@latest

# Show help
help:
	@echo "Available targets:"
	@echo "  all      - Build for all platforms (windows, linux, macos)"
	@echo "  windows  - Build for Windows (x64, x86)"
	@echo "  linux    - Build for Linux (x64, x86, arm64)"
	@echo "  macos    - Build for macOS (x64, arm64)"
	@echo "  build    - Build for current platform only"
	@echo "  test     - Run tests"
	@echo "  clean    - Clean build artifacts"
	@echo "  debug    - Build with debug information"
	@echo "  run      - Run agent locally"
	@echo "  fmt      - Format code"
	@echo "  vet      - Check for issues"
	@echo "  deps     - Install dependencies"
	@echo "  tools    - Install development tools"