# GoITDemon - Deployment and Usage Guide

## Overview

GoITDemon is a specialized Havoc Framework agent built in Go, designed specifically for IT administration tasks, system monitoring, and backup operations. It provides comprehensive cross-platform support for Windows and Linux environments.

## Key Features Implemented

### ✅ System Information & Monitoring
- **System Overview**: Complete system information including hardware, OS, and runtime details
- **Environment Information**: Environment variables, system paths, and configuration details  
- **System Specifications**: Detailed CPU, memory, and storage information
- **Uptime Tracking**: System uptime and boot time information

### ✅ Service Management
- **Cross-Platform Service Control**: Start, stop, restart services on Windows and Linux
- **Service Status Monitoring**: Check individual service status and health
- **Service Enumeration**: List all available services with current state
- **Platform-Specific Integration**: 
  - Windows: Uses `sc` command and PowerShell
  - Linux: Uses `systemctl` and legacy `service` commands

### ✅ Network Diagnostics
- **Ping Utility**: Test network connectivity to remote hosts
- **DNS Resolution**: Resolve hostnames to IP addresses
- **Network Interface Information**: List and monitor network interfaces
- **Cross-Platform Implementation**: Works on both Windows and Linux

### ✅ Hardware Information
- **CPU Information**: Detailed processor specifications and capabilities
- **Memory Information**: Physical and virtual memory statistics
- **Disk Information**: Storage device details and usage statistics
- **Comprehensive Hardware Inventory**: All hardware information in one command

### ✅ System Health Monitoring
- **Disk Usage Monitoring**: Real-time disk space utilization
- **Memory Usage Tracking**: Physical and swap memory consumption
- **CPU Usage Statistics**: Current processor utilization

### 🚧 Planned Features (Framework Ready)
- **Registry/Configuration Management**: Windows registry and Linux config file operations
- **Backup and Restoration**: File, directory, and system state backup capabilities
- **Event Log Monitoring**: System event analysis and reporting
- **Advanced Network Diagnostics**: Port scanning, traceroute, and network troubleshooting

## Build and Deployment

### Supported Platforms
- **Windows**: x64, x86 (Windows 10/11, Server 2016/2019/2022)
- **Linux**: x64, x86, ARM64 (Ubuntu, CentOS, RHEL, Debian)
- **macOS**: x64, ARM64 (for testing purposes)

### Build Instructions
```bash
# Clone and navigate to the project
cd /home/anoam/GitHub/Havoc/payloads/GoITDemon

# Install dependencies
go mod tidy

# Build for all platforms
./build.sh

# Or build for specific platform
go build -ldflags="-s -w" -o GoITDemon_linux_x64 .
```

### Generated Binaries
```
bin/
├── GoITDemon_windows_x64.exe    (6.7 MB)
├── GoITDemon_windows_x86.exe    (6.5 MB)
├── GoITDemon_linux_x64          (6.8 MB)
├── GoITDemon_linux_x86          (6.6 MB)
├── GoITDemon_linux_arm64        (6.4 MB)
├── GoITDemon_macos_x64          (6.7 MB)
└── GoITDemon_macos_arm64        (6.3 MB)
```

## Integration with Havoc Framework

### Command Structure
The GoITDemon uses the following command IDs that integrate with the Havoc teamserver:

```go
// IT Administration Commands
COMMAND_SYSINFO    = 3000  // System information and monitoring
COMMAND_SERVICE    = 3010  // Service management operations  
COMMAND_REGISTRY   = 3020  // Registry/configuration management
COMMAND_BACKUP     = 3030  // Backup and restoration operations
COMMAND_NETDIAG    = 3040  // Network diagnostics and troubleshooting
COMMAND_HWINFO     = 3050  // Hardware information gathering
COMMAND_SOFTINFO   = 3060  // Software inventory and management
COMMAND_SYSHEALTH  = 3070  // System health monitoring
```

### Communication Protocol
- **Transport**: HTTP/HTTPS (SMB support planned)
- **Metadata**: Comprehensive system information on initial checkin
- **Response Format**: Binary packets compatible with Havoc teamserver
- **Error Handling**: Graceful error responses with detailed messages

### Agent Identification
- **Agent Type**: "GoITDemon"
- **Magic Value**: 0x41424344 (ABCD)
- **Version**: 1.0.0
- **Platform Detection**: Automatic Windows/Linux/macOS detection

## Usage Examples

### Basic Operation
```bash
# Windows
.\GoITDemon_windows_x64.exe

# Linux  
./GoITDemon_linux_x64

# The agent will connect to 127.0.0.1:40056 by default
```

### Available Commands (via Havoc Client)
```
# System Information
sysinfo overview     - Complete system overview with hardware details
sysinfo environment  - Environment variables and system paths  
sysinfo specs       - Detailed hardware specifications
sysinfo uptime      - System uptime and boot information

# Service Management
service list        - List all system services
service start <name> - Start a specific service
service stop <name>  - Stop a specific service  
service restart <name> - Restart a service
service status <name> - Get detailed service status

# Network Diagnostics
netdiag ping <host>     - Ping a remote host
netdiag dns <host>      - Resolve hostname to IP
netdiag interfaces     - List network interfaces

# Hardware Information
hwinfo cpu         - CPU specifications and capabilities
hwinfo memory      - Memory information and usage
hwinfo disk        - Disk information and usage
hwinfo all         - Complete hardware inventory

# System Health
syshealth disk     - Current disk usage statistics  
syshealth memory   - Memory usage and availability
syshealth cpu      - Current CPU utilization
```

## Testing and Validation

### Capability Test Results
```
=== GoITDemon Capability Test ===

[+] Testing System Information...
[+] System overview: 264 bytes returned ✅
[+] Uptime info: 106 bytes returned ✅

[+] Testing Service Management...  
[+] Service list: 9746 bytes returned ✅

[+] Testing Network Diagnostics...
[+] Network interfaces: 256 bytes returned ✅

[+] Testing Hardware Information...
[+] CPU info: 1048 bytes returned ✅
[+] Memory info: 44 bytes returned ✅

[+] All tests completed! ✅
```

### Platform-Specific Features

#### Windows
- Service management via `sc` and PowerShell
- Windows-specific system information
- Registry operations (framework ready)
- Windows event log access (planned)

#### Linux  
- Service management via `systemctl` and `service`
- Package manager integration for software inventory
- Configuration file management (planned)
- System log analysis (planned)

## Security Considerations

### Operational Security
- **Legitimate IT Tool Appearance**: Designed to blend in with normal IT operations
- **Safe Operations**: Read-only operations by default, destructive operations require explicit confirmation
- **Error Handling**: Graceful error handling prevents crashes and maintains stealth
- **Resource Efficiency**: Minimal CPU and memory footprint

### Network Security
- **Encrypted Communication**: All teamserver communication uses HTTPS
- **Certificate Validation**: Supports custom CA certificates
- **Proxy Support**: Corporate proxy compatibility
- **Configurable Sleep**: Jitter and sleep patterns to avoid detection

## Roadmap and Future Development

### Phase 2 Enhancements
- **Registry Management**: Complete Windows registry backup and modification
- **Advanced Backup**: Scheduled backups with compression and encryption
- **Event Log Analysis**: Real-time event monitoring and alerting
- **Performance Monitoring**: Extended system performance metrics

### Phase 3 Features  
- **SMB Transport**: Named pipe communication for environments blocking HTTP
- **Enterprise Integration**: SIEM and monitoring tool integration
- **Compliance Reporting**: Automated compliance checking and reporting
- **GUI Enhancements**: Rich client interface for IT administrators

## Conclusion

GoITDemon represents a significant advancement in IT-focused red team capabilities, providing:

1. **Cross-Platform Compatibility**: Single codebase supporting Windows and Linux
2. **IT Administrator Focus**: Tools designed for legitimate IT operations
3. **Havoc Framework Integration**: Seamless integration with existing Havoc infrastructure
4. **Production-Ready**: Comprehensive error handling and robust operation
5. **Extensible Architecture**: Framework ready for additional capabilities

The agent successfully demonstrates that Go can be effectively used to create sophisticated, cross-platform agents that integrate seamlessly with the Havoc Framework while providing specialized capabilities for IT administration scenarios.

### Files Created
- **Main Agent**: `/home/anoam/GitHub/Havoc/payloads/GoITDemon/`
- **Binaries**: `/home/anoam/GitHub/Havoc/payloads/GoITDemon/bin/`
- **Documentation**: `README.md`, usage guides, and build instructions
- **Test Suite**: Capability validation and integration testing

This implementation provides a solid foundation for IT-focused red team operations and can be easily extended with additional capabilities as needed.