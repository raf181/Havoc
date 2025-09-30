package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os/exec"
	"runtime"
)

// NetDiagHandler handles network diagnostic commands
type NetDiagHandler struct{}

// NewNetDiagHandler creates a new network diagnostics handler
func NewNetDiagHandler() *NetDiagHandler {
	return &NetDiagHandler{}
}

// Handle processes network diagnostic commands
func (n *NetDiagHandler) Handle(data []byte) ([]byte, error) {
	if len(data) < 4 {
		return writeError(fmt.Errorf("invalid netdiag command data")), nil
	}

	subCommand := binary.LittleEndian.Uint32(data[:4])
	offset := 4

	switch subCommand {
	case NETDIAG_PING:
		hostname, _ := readString(data, offset)
		return n.pingHost(hostname)
	case NETDIAG_DNS_RESOLVE:
		hostname, _ := readString(data, offset)
		return n.resolveHost(hostname)
	case NETDIAG_INTERFACES:
		return n.getNetworkInterfaces()
	default:
		buf := new(bytes.Buffer)
		writeString(buf, fmt.Sprintf("Network diagnostic command %d not yet implemented", subCommand))
		return writeSuccess(buf.Bytes()), nil
	}
}

// pingHost pings a host
func (n *NetDiagHandler) pingHost(hostname string) ([]byte, error) {
	buf := new(bytes.Buffer)
	
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "4", hostname)
	} else {
		cmd = exec.Command("ping", "-c", "4", hostname)
	}
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return writeError(fmt.Errorf("ping failed: %v", err)), nil
	}
	
	writeString(buf, hostname)
	writeString(buf, string(output))
	return writeSuccess(buf.Bytes()), nil
}

// resolveHost resolves a hostname to IP addresses
func (n *NetDiagHandler) resolveHost(hostname string) ([]byte, error) {
	buf := new(bytes.Buffer)
	
	addrs, err := net.LookupHost(hostname)
	if err != nil {
		return writeError(fmt.Errorf("DNS resolution failed: %v", err)), nil
	}
	
	writeString(buf, hostname)
	binary.Write(buf, binary.LittleEndian, uint32(len(addrs)))
	for _, addr := range addrs {
		writeString(buf, addr)
	}
	
	return writeSuccess(buf.Bytes()), nil
}

// getNetworkInterfaces returns network interface information
func (n *NetDiagHandler) getNetworkInterfaces() ([]byte, error) {
	buf := new(bytes.Buffer)
	
	interfaces, err := net.Interfaces()
	if err != nil {
		return writeError(fmt.Errorf("failed to get interfaces: %v", err)), nil
	}
	
	binary.Write(buf, binary.LittleEndian, uint32(len(interfaces)))
	for _, iface := range interfaces {
		writeString(buf, iface.Name)
		writeString(buf, iface.HardwareAddr.String())
		binary.Write(buf, binary.LittleEndian, uint32(iface.Flags))
		
		addrs, _ := iface.Addrs()
		binary.Write(buf, binary.LittleEndian, uint32(len(addrs)))
		for _, addr := range addrs {
			writeString(buf, addr.String())
		}
	}
	
	return writeSuccess(buf.Bytes()), nil
}