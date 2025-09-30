package main

import (
	"fmt"
	"log"

	"github.com/HavocFramework/GoITDemon/pkg/commands"
)

// Simple test program to demonstrate the IT demon capabilities
func main() {
	fmt.Println("=== GoITDemon Capability Test ===")
	
	// Create command handler
	handler := commands.NewHandler()
	
	// Test system information
	fmt.Println("\n[+] Testing System Information...")
	testSysInfo(handler)
	
	// Test service management
	fmt.Println("\n[+] Testing Service Management...")
	testServiceManagement(handler)
	
	// Test network diagnostics
	fmt.Println("\n[+] Testing Network Diagnostics...")
	testNetworkDiag(handler)
	
	// Test hardware information
	fmt.Println("\n[+] Testing Hardware Information...")
	testHardwareInfo(handler)
	
	fmt.Println("\n[+] All tests completed!")
}

func testSysInfo(handler *commands.Handler) {
	// Test system overview
	data := createCommand(1) // SYSINFO_OVERVIEW
	result, err := handler.Execute(3000, data) // COMMAND_SYSINFO
	if err != nil {
		log.Printf("[-] System overview test failed: %v", err)
	} else {
		fmt.Printf("[+] System overview: %d bytes returned\n", len(result))
	}
	
	// Test uptime
	data = createCommand(4) // SYSINFO_UPTIME
	result, err = handler.Execute(3000, data)
	if err != nil {
		log.Printf("[-] Uptime test failed: %v", err)
	} else {
		fmt.Printf("[+] Uptime info: %d bytes returned\n", len(result))
	}
}

func testServiceManagement(handler *commands.Handler) {
	// Test service list
	data := createCommand(1) // SERVICE_LIST
	result, err := handler.Execute(3010, data) // COMMAND_SERVICE
	if err != nil {
		log.Printf("[-] Service list test failed: %v", err)
	} else {
		fmt.Printf("[+] Service list: %d bytes returned\n", len(result))
	}
}

func testNetworkDiag(handler *commands.Handler) {
	// Test network interfaces
	data := createCommand(5) // NETDIAG_INTERFACES
	result, err := handler.Execute(3040, data) // COMMAND_NETDIAG
	if err != nil {
		log.Printf("[-] Network interfaces test failed: %v", err)
	} else {
		fmt.Printf("[+] Network interfaces: %d bytes returned\n", len(result))
	}
}

func testHardwareInfo(handler *commands.Handler) {
	// Test CPU info
	data := createCommand(1) // HWINFO_CPU
	result, err := handler.Execute(3050, data) // COMMAND_HWINFO
	if err != nil {
		log.Printf("[-] CPU info test failed: %v", err)
	} else {
		fmt.Printf("[+] CPU info: %d bytes returned\n", len(result))
	}
	
	// Test memory info
	data = createCommand(2) // HWINFO_MEMORY
	result, err = handler.Execute(3050, data)
	if err != nil {
		log.Printf("[-] Memory info test failed: %v", err)
	} else {
		fmt.Printf("[+] Memory info: %d bytes returned\n", len(result))
	}
}

// Helper function to create command data
func createCommand(subCmd uint32) []byte {
	data := make([]byte, 4)
	data[0] = byte(subCmd)
	data[1] = byte(subCmd >> 8)
	data[2] = byte(subCmd >> 16)
	data[3] = byte(subCmd >> 24)
	return data
}