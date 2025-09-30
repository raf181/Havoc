# GoITDemon Integration with Havoc Framework

This document describes the integration of GoITDemon with the Havoc Framework, enabling it to be built from the dashboard and managed through the client.

## Overview

GoITDemon is now fully integrated with the Havoc Framework as a service agent, providing:

- **Dashboard Building**: Build GoITDemon directly from the Havoc client interface
- **Client Management**: Manage GoITDemon agents through the Havoc client
- **Service Registration**: Automatic registration with the teamserver
- **Cross-Platform Support**: Windows, Linux, and macOS builds

## Architecture

The integration consists of several components:

### 1. Service Agent (`service/main.go`)
- Registers with the Havoc teamserver service endpoint
- Handles build requests from the dashboard
- Provides agent metadata and configuration options
- Builds GoITDemon binaries with custom configurations

### 2. Enhanced GoITDemon Agent (`main.go`, `pkg/agent/`)
- Compatible with Havoc packet format
- Proper authentication with teamserver
- Service agent registration support
- IT administration command handlers

### 3. Build Integration (`Makefile`, `start-service.sh`)
- Automated building of service agent
- Easy startup scripts
- Cross-platform compilation

## Setup Instructions

### 1. Start the Teamserver

Start the Havoc teamserver with the default profile:

```bash
cd /home/anoam/Havoc
./havoc server --profile profiles/havoc.yaotl --verbose
```

### 2. Start the Service Agent

```bash
cd /home/anoam/Havoc/payloads/GoITDemon
./start-service.sh
```

### 3. Configure a Listener

Before deploying GoITDemon agents, you need to configure a proper HTTP/HTTPS listener in the Havoc client:

1. Open the Havoc client and connect to the teamserver
2. Navigate to the Listeners section
3. Create a new HTTP listener with:
   - **Name**: GoITDemon-Listener
   - **Bind Address**: 0.0.0.0
   - **Port**: 8080 (or another available port)
   - **Host**: 127.0.0.1 (or your server IP)
   - **URIs**: `/fwlink`, `/pixel.gif`, `/updates`, `/download`

### 4. Verify Registration

Check the teamserver logs to confirm the service agent has registered:

```
[SERVICE] GoITDemon registered a new agent [Name: GoITDemon]
```

### 3. Build from Dashboard

1. Open the Havoc client
2. Navigate to the payload generation dialog
3. Select "GoITDemon" from the agent type dropdown
4. Configure the build options:
   - **Architecture**: x64, x86, arm64
   - **Format**: Linux Executable, Windows Executable, macOS Executable
   - **Sleep Time**: Agent callback interval (seconds)
   - **Jitter**: Randomization percentage
   - **User Agent**: HTTP User-Agent string
   - **Kill Date**: Optional termination date
   - **Working Hours**: Optional operational time window
5. Select a listener
6. Click "Generate"

### 4. Deploy and Execute

The generated binary will be saved with the format: `GoITDemon_[os]_[arch][.ext]`

Execute the agent on the target system. The agent will:
1. Connect to the specified teamserver
2. Authenticate using the GoITDemon magic value (`0x474f4954`)
3. Register with the teamserver
4. Appear in the client's session list

## Agent Features

### System Information Commands
- `sysinfo overview` - Complete system overview
- `sysinfo environment` - Environment variables and paths
- `sysinfo specs` - Detailed hardware specifications
- `sysinfo uptime` - System uptime information

### Service Management Commands
- `service list` - List all services
- `service start <name>` - Start a service
- `service stop <name>` - Stop a service
- `service restart <name>` - Restart a service
- `service status <name>` - Get service status

### Network Diagnostics Commands
- `netdiag ping <host>` - Ping a host
- `netdiag dns <host>` - Resolve hostname
- `netdiag interfaces` - List network interfaces

### Hardware Information Commands
- `hwinfo cpu` - CPU information
- `hwinfo memory` - Memory information
- `hwinfo disk` - Disk information
- `hwinfo all` - All hardware information

### System Health Commands
- `syshealth disk` - Disk usage information
- `syshealth memory` - Memory usage information
- `syshealth cpu` - CPU usage information

## Configuration Options

### Service Agent Configuration

Environment variables for the service agent:

- `TEAMSERVER_HOST`: Teamserver hostname (default: 127.0.0.1)
- `TEAMSERVER_PORT`: Teamserver port (default: 40056)
- `SERVICE_PASSWORD`: Service authentication password (default: password123)

### Agent Configuration

Environment variables for the GoITDemon agent:

- `GOITDEMON_SLEEP`: Sleep time in seconds (default: 5)
- `GOITDEMON_JITTER`: Jitter percentage (default: 10)
- `GOITDEMON_HOST`: Teamserver host (default: 127.0.0.1)
- `GOITDEMON_PORT`: Teamserver port (default: 40056)
- `GOITDEMON_SECURE`: Use HTTPS (default: false)
- `GOITDEMON_USERAGENT`: HTTP User-Agent string

## Troubleshooting

### Service Agent Issues

1. **Connection Failed**: Check teamserver status and port
2. **Authentication Failed**: Verify service password
3. **Build Failed**: Check Go installation and dependencies

```bash
# Check service agent logs
./start-service.sh

# Manual build test
cd service
go build -o test-service .
```

### Agent Issues

1. **Connection Failed**: Verify listener configuration
2. **Magic Value Mismatch**: Ensure consistent magic value usage
3. **Command Failures**: Check agent logs and command syntax

```bash
# Test agent build
make build

# Test with debug output
GOITDEMON_DEBUG=true ./bin/GoITDemon
```

### Common Solutions

1. **Port Conflicts**: Change default ports if in use
2. **Firewall Issues**: Ensure ports are accessible
3. **Certificate Issues**: Use HTTP for testing, HTTPS for production

## Development

### Building Components

```bash
# Build everything
make all

# Build service agent only
make service

# Build specific platform
make windows
make linux
make macos

# Clean build artifacts
make clean
```

### Customization

To add new commands:

1. Define command ID in `pkg/commands/handler.go`
2. Implement handler in appropriate file
3. Update service agent command list
4. Test with manual build

### Debugging

Enable debug output:

```bash
# Service agent debug
./bin/GoITDemon-Service -host 127.0.0.1 -port 40056 -password password123 -debug

# Agent debug (set in environment)
export GOITDEMON_DEBUG=true
./bin/GoITDemon
```

## Security Considerations

### Operational Security
- Use HTTPS in production environments
- Rotate service passwords regularly
- Monitor agent connections and commands
- Use jitter to avoid detection patterns

### Network Security
- Configure firewalls appropriately
- Use custom User-Agent strings
- Consider proxy support for corporate environments
- Implement certificate pinning for HTTPS

### Agent Security
- Avoid excessive logging in production
- Use kill dates for time-limited operations
- Implement working hours for stealth
- Monitor resource usage

## Integration Points

### Teamserver Integration
- Service endpoint: `/service/register`
- Agent communication: Standard Havoc HTTP/HTTPS
- Magic value: `0x474f4954` ("GOIT")
- Packet format: Compatible with Havoc standards

### Client Integration
- Agent type: "GoITDemon"
- Commands: IT administration focused
- UI: Standard Havoc payload generation dialog
- Management: Standard session management

## License and Credits

This integration maintains compatibility with the Havoc Framework license. GoITDemon provides IT administration capabilities while following Havoc's architecture and security model.