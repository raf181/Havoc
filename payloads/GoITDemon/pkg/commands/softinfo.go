package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// SoftInfoHandler handles software information commands
type SoftInfoHandler struct{}

// NewSoftInfoHandler creates a new software info handler
func NewSoftInfoHandler() *SoftInfoHandler {
	return &SoftInfoHandler{}
}

// Handle processes software info commands
func (s *SoftInfoHandler) Handle(data []byte) ([]byte, error) {
	if len(data) < 4 {
		return writeError(fmt.Errorf("invalid softinfo command data")), nil
	}

	subCommand := binary.LittleEndian.Uint32(data[:4])

	switch subCommand {
	case SOFTINFO_INSTALLED:
		return s.getInstalledSoftware()
	case SOFTINFO_RUNNING:
		return s.getRunningProcesses()
	default:
		buf := new(bytes.Buffer)
		writeString(buf, fmt.Sprintf("Software info command %d not yet implemented", subCommand))
		return writeSuccess(buf.Bytes()), nil
	}
}

// getInstalledSoftware returns installed software list
func (s *SoftInfoHandler) getInstalledSoftware() ([]byte, error) {
	buf := new(bytes.Buffer)
	
	if runtime.GOOS == "windows" {
		return s.getWindowsInstalledSoftware(buf)
	} else {
		return s.getLinuxInstalledSoftware(buf)
	}
}

// getWindowsInstalledSoftware gets installed software on Windows
func (s *SoftInfoHandler) getWindowsInstalledSoftware(buf *bytes.Buffer) ([]byte, error) {
	// Use PowerShell to get installed software
	cmd := exec.Command("powershell", "-Command", "Get-WmiObject -Class Win32_Product | Select-Object Name, Version, Vendor")
	output, err := cmd.Output()
	if err != nil {
		return writeError(fmt.Errorf("failed to get Windows software list: %v", err)), nil
	}
	
	lines := strings.Split(string(output), "\n")
	var software []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "Name") && !strings.HasPrefix(line, "----") {
			software = append(software, line)
		}
	}
	
	binary.Write(buf, binary.LittleEndian, uint32(len(software)))
	for _, item := range software {
		writeString(buf, item)
	}
	
	return writeSuccess(buf.Bytes()), nil
}

// getLinuxInstalledSoftware gets installed software on Linux
func (s *SoftInfoHandler) getLinuxInstalledSoftware(buf *bytes.Buffer) ([]byte, error) {
	var cmd *exec.Cmd
	
	// Try different package managers
	if _, err := exec.LookPath("dpkg"); err == nil {
		// Debian/Ubuntu
		cmd = exec.Command("dpkg", "-l")
	} else if _, err := exec.LookPath("rpm"); err == nil {
		// Red Hat/CentOS
		cmd = exec.Command("rpm", "-qa")
	} else if _, err := exec.LookPath("pacman"); err == nil {
		// Arch Linux
		cmd = exec.Command("pacman", "-Q")
	} else {
		return writeError(fmt.Errorf("no supported package manager found")), nil
	}
	
	output, err := cmd.Output()
	if err != nil {
		return writeError(fmt.Errorf("failed to get Linux software list: %v", err)), nil
	}
	
	lines := strings.Split(string(output), "\n")
	var software []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			software = append(software, line)
		}
	}
	
	binary.Write(buf, binary.LittleEndian, uint32(len(software)))
	for _, item := range software {
		writeString(buf, item)
	}
	
	return writeSuccess(buf.Bytes()), nil
}

// getRunningProcesses returns currently running processes
func (s *SoftInfoHandler) getRunningProcesses() ([]byte, error) {
	buf := new(bytes.Buffer)
	
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("tasklist", "/fo", "csv")
	} else {
		cmd = exec.Command("ps", "aux")
	}
	
	output, err := cmd.Output()
	if err != nil {
		return writeError(fmt.Errorf("failed to get process list: %v", err)), nil
	}
	
	lines := strings.Split(string(output), "\n")
	var processes []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			processes = append(processes, line)
		}
	}
	
	binary.Write(buf, binary.LittleEndian, uint32(len(processes)))
	for _, proc := range processes {
		writeString(buf, proc)
	}
	
	return writeSuccess(buf.Bytes()), nil
}