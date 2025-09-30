# GoITDemon - IT Administration Focused Havoc Agent

GoITDemon is a specialized Havoc Framework agent written in Go, designed specifically for IT administration tasks, system monitoring, and backup operations. It provides cross-platform support for Windows and Linux environments.

## Features

### System Information & Monitoring
- **System Overview**: Comprehensive system information including hardware, OS, and runtime details
- **Environment Info**: Environment variables, system paths, and configuration details
- **Hardware Specs**: Detailed CPU, memory, and storage information
- **System Health**: Real-time monitoring of CPU, memory, and disk usage
- **Uptime Tracking**: System uptime and boot time information

### Service Management
- **Cross-Platform Service Control**: Start, stop, restart services on both Windows and Linux
- **Service Status Monitoring**: Check service status and health
- **Service Enumeration**: List all available services with their current state
- **Service Dependencies**: View service dependency information

### Network Diagnostics
- **Ping Utility**: Test network connectivity to hosts
- **DNS Resolution**: Resolve hostnames to IP addresses
- **Network Interfaces**: List and monitor network interface information
- **Port Scanning**: Basic port availability testing
- **Network Health**: Monitor network connectivity and performance

### Backup & Restoration (Planned)
- **File Backup**: Backup critical files and directories
- **Registry Backup**: Windows registry backup and restoration
- **System State Backup**: Complete system state preservation
- **Incremental Backups**: Space-efficient backup solutions

### Security & Compliance
- **Safe Operations**: All operations designed with IT safety in mind
- **Audit Logging**: Comprehensive logging of all operations
- **Non-Destructive**: Read-only operations by default
- **Permission Checking**: Verify adequate permissions before operations

## Supported Platforms

### Windows
- Windows 10/11 (x64, x86)
- Windows Server 2016/2019/2022
- Service management via `sc` command
- PowerShell integration for advanced operations
- Windows Registry operations (planned)

### Linux
- Ubuntu/Debian (x64, x86, ARM64)
- CentOS/RHEL (x64, x86, ARM64)
- Service management via `systemctl` and `service`
- Package manager integration
- Configuration file management

## Building

### Prerequisites
- Go 1.21 or later
- Make (optional)

### Quick Build
```bash
# Build for all platforms
make all

# Or use the build script
chmod +x build.sh
./build.sh

# Build for specific platform
make windows  # Windows binaries
make linux    # Linux binaries
make macos    # macOS binaries (for testing)
```

### Manual Build
```bash
# Windows x64
GOOS=windows GOARCH=amd64 go build -o GoITDemon_windows_x64.exe .

# Linux x64
GOOS=linux GOARCH=amd64 go build -o GoITDemon_linux_x64 .
```

## Usage

### Basic Operation
```bash
# Windows
.\\GoITDemon_windows_x64.exe

# Linux
./GoITDemon_linux_x64

# The agent will attempt to connect to 127.0.0.1:40056 by default
```

### Available Commands

#### System Information
```
sysinfo overview     - Complete system overview
sysinfo environment  - Environment variables and paths
sysinfo specs       - Detailed hardware specifications
sysinfo uptime      - System uptime information
```

#### Service Management
```
service list        - List all services
service start <name> - Start a service
service stop <name>  - Stop a service
service restart <name> - Restart a service
service status <name> - Get service status
```

#### Network Diagnostics
```
netdiag ping <host>     - Ping a host
netdiag dns <host>      - Resolve hostname
netdiag interfaces     - List network interfaces
```

#### Hardware Information
```
hwinfo cpu         - CPU information
hwinfo memory      - Memory information
hwinfo disk        - Disk information
hwinfo all         - All hardware information
```

#### System Health
```
syshealth disk     - Disk usage information
syshealth memory   - Memory usage information
syshealth cpu      - CPU usage information
```

## Integration with Havoc Framework

### Teamserver Configuration
The GoITDemon integrates with the existing Havoc teamserver infrastructure:

1. **Agent Registration**: Automatically registers as "GoITDemon" agent type
2. **Command Compatibility**: Uses Havoc-compatible command IDs and packet formats
3. **Transport Layer**: Supports HTTP/HTTPS communication (SMB planned)
4. **Metadata Exchange**: Provides comprehensive system metadata on checkin

### Client Integration
To add support for GoITDemon in the Havoc client:

1. Add new command definitions for IT administration
2. Update command handlers to process GoITDemon responses
3. Add UI elements for IT-specific operations
4. Implement result visualization for system health data

## Security Considerations

### Operational Security
- **Stealth Operations**: Designed to blend in with normal IT activities
- **Minimal Footprint**: Efficient memory and CPU usage
- **Safe Defaults**: Read-only operations unless explicitly requested
- **Error Handling**: Graceful error handling to prevent crashes

### Network Security
- **Encrypted Communication**: All communication with teamserver is encrypted
- **Certificate Validation**: Supports custom CA certificates
- **Proxy Support**: Can operate through corporate proxies
- **Jitter and Sleep**: Configurable communication patterns

## Development

### Adding New Commands
1. Define command ID in `pkg/commands/handler.go`
2. Create handler file in `pkg/commands/`
3. Implement command logic with proper error handling
4. Add cross-platform support where applicable
5. Update documentation

### Testing
```bash
# Run tests
make test

# Run with debug information
make debug

# Local testing
make run
```

### Code Style
- Follow Go conventions and best practices
- Use `go fmt` for formatting
- Add comprehensive error handling
- Include cross-platform compatibility checks

## Roadmap

### Phase 1 (Current)
- [x] Basic system information gathering
- [x] Cross-platform service management
- [x] Network diagnostic tools
- [x] Hardware monitoring
- [x] System health monitoring

### Phase 2 (Planned)
- [ ] Registry/configuration management
- [ ] File and system backup operations
- [ ] Event log monitoring and analysis
- [ ] Performance monitoring and alerts
- [ ] Software inventory and updates

### Phase 3 (Future)
- [ ] SMB transport support
- [ ] Advanced backup scheduling
- [ ] Compliance checking and reporting
- [ ] Integration with enterprise tools
- [ ] GUI client enhancements

## License

This project is part of the Havoc Framework and follows the same licensing terms.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes with proper testing
4. Submit a pull request with detailed description
5. Ensure all tests pass and code is properly formatted

## Support

For issues and questions:
1. Check the existing documentation
2. Search for similar issues in the repository
3. Create a detailed issue report with:
   - Operating system and version
   - Go version used
   - Exact error messages or unexpected behavior
   - Steps to reproduce the issue

## Disclaimer

This tool is designed for legitimate IT administration and system monitoring purposes. Users are responsible for ensuring compliance with applicable laws and regulations in their jurisdiction. The authors are not responsible for any misuse of this software.