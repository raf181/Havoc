package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os/exec"
	"strings"
)

// ServiceHandler handles service management commands
type ServiceHandler struct{}

// NewServiceHandler creates a new service handler
func NewServiceHandler() *ServiceHandler {
	return &ServiceHandler{}
}

// Handle processes service management commands
func (s *ServiceHandler) Handle(data []byte) ([]byte, error) {
	if len(data) < 4 {
		return writeError(fmt.Errorf("invalid service command data")), nil
	}

	subCommand := binary.LittleEndian.Uint32(data[:4])
	offset := 4

	switch subCommand {
	case SERVICE_LIST:
		return s.listServices()
	case SERVICE_START:
		serviceName, _ := readString(data, offset)
		return s.startService(serviceName)
	case SERVICE_STOP:
		serviceName, _ := readString(data, offset)
		return s.stopService(serviceName)
	case SERVICE_RESTART:
		serviceName, _ := readString(data, offset)
		return s.restartService(serviceName)
	case SERVICE_STATUS:
		serviceName, _ := readString(data, offset)
		return s.getServiceStatus(serviceName)
	default:
		return writeError(fmt.Errorf("unknown service subcommand: %d", subCommand)), nil
	}
}

// listServices returns a list of all services
func (s *ServiceHandler) listServices() ([]byte, error) {
	buf := new(bytes.Buffer)

	if isWindows() {
		return s.listWindowsServices(buf)
	} else {
		return s.listLinuxServices(buf)
	}
}

// listWindowsServices lists Windows services using sc command
func (s *ServiceHandler) listWindowsServices(buf *bytes.Buffer) ([]byte, error) {
	cmd := exec.Command("sc", "query", "state=", "all")
	output, err := cmd.Output()
	if err != nil {
		return writeError(fmt.Errorf("failed to list Windows services: %v", err)), nil
	}

	services := s.parseWindowsServices(string(output))
	binary.Write(buf, binary.LittleEndian, uint32(len(services)))

	for _, service := range services {
		writeString(buf, service.Name)
		writeString(buf, service.DisplayName)
		writeString(buf, service.State)
		writeString(buf, service.Type)
	}

	return writeSuccess(buf.Bytes()), nil
}

// listLinuxServices lists Linux services using systemctl
func (s *ServiceHandler) listLinuxServices(buf *bytes.Buffer) ([]byte, error) {
	cmd := exec.Command("systemctl", "list-units", "--type=service", "--all", "--no-pager")
	output, err := cmd.Output()
	if err != nil {
		// Fallback to service command
		cmd = exec.Command("service", "--status-all")
		output, err = cmd.Output()
		if err != nil {
			return writeError(fmt.Errorf("failed to list Linux services: %v", err)), nil
		}
		return s.parseLinuxServicesLegacy(buf, string(output))
	}

	services := s.parseLinuxServices(string(output))
	binary.Write(buf, binary.LittleEndian, uint32(len(services)))

	for _, service := range services {
		writeString(buf, service.Name)
		writeString(buf, service.DisplayName)
		writeString(buf, service.State)
		writeString(buf, service.Type)
	}

	return writeSuccess(buf.Bytes()), nil
}

// Service represents a system service
type Service struct {
	Name        string
	DisplayName string
	State       string
	Type        string
}

// parseWindowsServices parses sc query output
func (s *ServiceHandler) parseWindowsServices(output string) []Service {
	var services []Service
	lines := strings.Split(output, "\n")
	
	var currentService Service
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "SERVICE_NAME:") {
			if currentService.Name != "" {
				services = append(services, currentService)
			}
			currentService = Service{
				Name: strings.TrimSpace(strings.TrimPrefix(line, "SERVICE_NAME:")),
			}
		} else if strings.HasPrefix(line, "DISPLAY_NAME:") {
			currentService.DisplayName = strings.TrimSpace(strings.TrimPrefix(line, "DISPLAY_NAME:"))
		} else if strings.HasPrefix(line, "STATE") {
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				currentService.State = parts[3]
			}
		} else if strings.HasPrefix(line, "TYPE") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				currentService.Type = strings.Join(parts[2:], " ")
			}
		}
	}
	if currentService.Name != "" {
		services = append(services, currentService)
	}

	return services
}

// parseLinuxServices parses systemctl output
func (s *ServiceHandler) parseLinuxServices(output string) []Service {
	var services []Service
	lines := strings.Split(output, "\n")
	
	// Skip header lines
	for i, line := range lines {
		if i < 1 || strings.TrimSpace(line) == "" {
			continue
		}
		
		fields := strings.Fields(line)
		if len(fields) >= 4 {
			service := Service{
				Name:        fields[0],
				DisplayName: strings.Join(fields[4:], " "),
				State:       fields[2],
				Type:        "service",
			}
			
			// Remove .service suffix if present
			if strings.HasSuffix(service.Name, ".service") {
				service.Name = strings.TrimSuffix(service.Name, ".service")
			}
			
			services = append(services, service)
		}
	}

	return services
}

// parseLinuxServicesLegacy parses legacy service command output
func (s *ServiceHandler) parseLinuxServicesLegacy(buf *bytes.Buffer, output string) ([]byte, error) {
	lines := strings.Split(output, "\n")
	var services []Service

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Parse [ + ] servicename or [ - ] servicename format
		if strings.Contains(line, "]") {
			parts := strings.SplitN(line, "]", 2)
			if len(parts) == 2 {
				status := "unknown"
				if strings.Contains(parts[0], "+") {
					status = "running"
				} else if strings.Contains(parts[0], "-") {
					status = "stopped"
				}
				
				serviceName := strings.TrimSpace(parts[1])
				service := Service{
					Name:        serviceName,
					DisplayName: serviceName,
					State:       status,
					Type:        "service",
				}
				services = append(services, service)
			}
		}
	}

	binary.Write(buf, binary.LittleEndian, uint32(len(services)))
	for _, service := range services {
		writeString(buf, service.Name)
		writeString(buf, service.DisplayName)
		writeString(buf, service.State)
		writeString(buf, service.Type)
	}

	return writeSuccess(buf.Bytes()), nil
}

// startService starts a service
func (s *ServiceHandler) startService(serviceName string) ([]byte, error) {
	var cmd *exec.Cmd
	
	if isWindows() {
		cmd = exec.Command("sc", "start", serviceName)
	} else {
		// Try systemctl first, then service command
		cmd = exec.Command("systemctl", "start", serviceName)
		if err := cmd.Run(); err != nil {
			cmd = exec.Command("service", serviceName, "start")
		}
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return writeError(fmt.Errorf("failed to start service %s: %v - %s", serviceName, err, string(output))), nil
	}

	buf := new(bytes.Buffer)
	writeString(buf, fmt.Sprintf("Service %s started successfully", serviceName))
	writeString(buf, string(output))
	
	return writeSuccess(buf.Bytes()), nil
}

// stopService stops a service
func (s *ServiceHandler) stopService(serviceName string) ([]byte, error) {
	var cmd *exec.Cmd
	
	if isWindows() {
		cmd = exec.Command("sc", "stop", serviceName)
	} else {
		cmd = exec.Command("systemctl", "stop", serviceName)
		if err := cmd.Run(); err != nil {
			cmd = exec.Command("service", serviceName, "stop")
		}
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return writeError(fmt.Errorf("failed to stop service %s: %v - %s", serviceName, err, string(output))), nil
	}

	buf := new(bytes.Buffer)
	writeString(buf, fmt.Sprintf("Service %s stopped successfully", serviceName))
	writeString(buf, string(output))
	
	return writeSuccess(buf.Bytes()), nil
}

// restartService restarts a service
func (s *ServiceHandler) restartService(serviceName string) ([]byte, error) {
	var cmd *exec.Cmd
	
	if isWindows() {
		// Windows doesn't have a direct restart, so stop then start
		stopCmd := exec.Command("sc", "stop", serviceName)
		stopCmd.Run() // Ignore error if service wasn't running
		
		cmd = exec.Command("sc", "start", serviceName)
	} else {
		cmd = exec.Command("systemctl", "restart", serviceName)
		if err := cmd.Run(); err != nil {
			cmd = exec.Command("service", serviceName, "restart")
		}
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return writeError(fmt.Errorf("failed to restart service %s: %v - %s", serviceName, err, string(output))), nil
	}

	buf := new(bytes.Buffer)
	writeString(buf, fmt.Sprintf("Service %s restarted successfully", serviceName))
	writeString(buf, string(output))
	
	return writeSuccess(buf.Bytes()), nil
}

// getServiceStatus gets the status of a specific service
func (s *ServiceHandler) getServiceStatus(serviceName string) ([]byte, error) {
	var cmd *exec.Cmd
	
	if isWindows() {
		cmd = exec.Command("sc", "query", serviceName)
	} else {
		cmd = exec.Command("systemctl", "status", serviceName, "--no-pager")
		if err := cmd.Run(); err != nil {
			cmd = exec.Command("service", serviceName, "status")
		}
	}

	output, err := cmd.CombinedOutput()
	
	buf := new(bytes.Buffer)
	writeString(buf, serviceName)
	
	if err != nil {
		writeString(buf, "unknown")
		writeString(buf, fmt.Sprintf("Error getting status: %v", err))
	} else {
		status := s.parseServiceStatus(string(output), isWindows())
		writeString(buf, status)
		writeString(buf, string(output))
	}
	
	return writeSuccess(buf.Bytes()), nil
}

// parseServiceStatus extracts service status from command output
func (s *ServiceHandler) parseServiceStatus(output string, isWindows bool) string {
	output = strings.ToLower(output)
	
	if isWindows {
		if strings.Contains(output, "running") {
			return "running"
		} else if strings.Contains(output, "stopped") {
			return "stopped"
		} else if strings.Contains(output, "paused") {
			return "paused"
		}
	} else {
		if strings.Contains(output, "active (running)") {
			return "running"
		} else if strings.Contains(output, "inactive") || strings.Contains(output, "dead") {
			return "stopped"
		} else if strings.Contains(output, "failed") {
			return "failed"
		} else if strings.Contains(output, "active (exited)") {
			return "exited"
		}
	}
	
	return "unknown"
}