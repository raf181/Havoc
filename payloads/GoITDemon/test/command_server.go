package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net/http"
)

var commandQueue = []Command{
	{ID: 1, CommandID: 3000, Data: []byte{1, 0, 0, 0}}, // SYSINFO_OVERVIEW
	{ID: 2, CommandID: 3010, Data: []byte{1, 0, 0, 0}}, // SERVICE_LIST
	{ID: 3, CommandID: 3040, Data: []byte{5, 0, 0, 0}}, // NETDIAG_INTERFACES
	{ID: 4, CommandID: 3050, Data: []byte{6, 0, 0, 0}}, // HWINFO_ALL
}

var currentCommand = 0

type Command struct {
	ID        uint32
	CommandID uint32
	Data      []byte
}

func main() {
	fmt.Println("[+] Starting GoITDemon test server with IT commands on :40056")

	http.HandleFunc("/", handleRequest)
	http.HandleFunc("/api/login", handleRequest)
	http.HandleFunc("/api/session", handleRequest)
	http.HandleFunc("/api/check", handleRequest)

	fmt.Println("[+] Server will send IT administration commands to test GoITDemon")
	fmt.Printf("[+] Commands queued: %d\n", len(commandQueue))
	fmt.Println("[+] Server listening on http://127.0.0.1:40056")

	log.Fatal(http.ListenAndServe(":40056", nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fmt.Printf("[*] %s %s from %s\n", r.Method, r.URL.Path, r.RemoteAddr)

	w.Header().Set("Content-Type", "application/octet-stream")

	if len(body) == 166 {
		// Initial metadata - parse some info
		fmt.Printf("    Received initial metadata (%d bytes)\n", len(body))
		if len(body) >= 8 {
			magicValue := binary.LittleEndian.Uint32(body[0:4])
			agentID := binary.LittleEndian.Uint32(body[4:8])
			fmt.Printf("    Magic: 0x%08x, Agent ID: 0x%08x\n", magicValue, agentID)
		}
		w.Write([]byte("OK"))
		fmt.Printf("    Response: Initial checkin acknowledged\n")
	} else if len(body) == 4 {
		// Task request
		agentID := binary.LittleEndian.Uint32(body)
		fmt.Printf("    Task request from agent 0x%08x\n", agentID)
		
		// Send next command if available
		if currentCommand < len(commandQueue) {
			cmd := commandQueue[currentCommand]
			response := createCommandPacket(cmd)
			w.Write(response)
			fmt.Printf("    Sent command %d: ID=%d, Type=%d, Data=%d bytes\n", 
				currentCommand+1, cmd.ID, cmd.CommandID, len(cmd.Data))
			currentCommand++
		} else {
			// No more commands
			w.Write([]byte{})
			fmt.Printf("    No more commands to send\n")
		}
	} else if len(body) > 4 {
		// Command response
		fmt.Printf("    Received command response (%d bytes)\n", len(body))
		if len(body) >= 8 {
			taskID := binary.LittleEndian.Uint32(body[0:4])
			dataLen := binary.LittleEndian.Uint32(body[4:8])
			fmt.Printf("    Task ID: %d, Data length: %d\n", taskID, dataLen)
			
			// Parse response data
			if len(body) >= 12 {
				successFlag := binary.LittleEndian.Uint32(body[8:12])
				fmt.Printf("    Success: %s\n", map[uint32]string{0: "Failed", 1: "Success"}[successFlag])
				
				if dataLen > 4 && len(body) >= int(12+dataLen) {
					responseData := body[12:12+dataLen]
					fmt.Printf("    Response preview: %s...\n", string(responseData[:min(int(dataLen), 50)]))
				}
			}
		}
		w.Write([]byte("ACK"))
	} else {
		// Empty request
		w.WriteHeader(http.StatusOK)
	}
	
	fmt.Println()
}

func createCommandPacket(cmd Command) []byte {
	packet := make([]byte, 16+len(cmd.Data))
	
	// Task ID
	binary.LittleEndian.PutUint32(packet[0:4], cmd.ID)
	
	// Command ID
	binary.LittleEndian.PutUint32(packet[4:8], cmd.CommandID)
	
	// Data length
	binary.LittleEndian.PutUint32(packet[8:12], uint32(len(cmd.Data)))
	
	// Reserved
	binary.LittleEndian.PutUint32(packet[12:16], 0)
	
	// Data
	copy(packet[16:], cmd.Data)
	
	return packet
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}