package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

// SysHealthHandler handles system health monitoring commands
type SysHealthHandler struct{}

// NewSysHealthHandler creates a new system health handler
func NewSysHealthHandler() *SysHealthHandler {
	return &SysHealthHandler{}
}

// Handle processes system health commands
func (s *SysHealthHandler) Handle(data []byte) ([]byte, error) {
	if len(data) < 4 {
		return writeError(fmt.Errorf("invalid syshealth command data")), nil
	}

	subCommand := binary.LittleEndian.Uint32(data[:4])

	switch subCommand {
	case SYSHEALTH_DISK:
		return s.getDiskUsage()
	case SYSHEALTH_MEMORY:
		return s.getMemoryUsage()
	case SYSHEALTH_CPU:
		return s.getCPUUsage()
	default:
		buf := new(bytes.Buffer)
		writeString(buf, fmt.Sprintf("System health command %d not yet implemented", subCommand))
		return writeSuccess(buf.Bytes()), nil
	}
}

// getDiskUsage returns disk usage information
func (s *SysHealthHandler) getDiskUsage() ([]byte, error) {
	buf := new(bytes.Buffer)
	
	partitions, err := disk.Partitions(false)
	if err != nil {
		return writeError(err), nil
	}
	
	binary.Write(buf, binary.LittleEndian, uint32(len(partitions)))
	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}
		
		writeString(buf, partition.Device)
		writeString(buf, partition.Mountpoint)
		binary.Write(buf, binary.LittleEndian, usage.Total)
		binary.Write(buf, binary.LittleEndian, usage.Used)
		binary.Write(buf, binary.LittleEndian, usage.Free)
		binary.Write(buf, binary.LittleEndian, usage.UsedPercent)
	}
	
	return writeSuccess(buf.Bytes()), nil
}

// getMemoryUsage returns memory usage information
func (s *SysHealthHandler) getMemoryUsage() ([]byte, error) {
	buf := new(bytes.Buffer)
	
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return writeError(err), nil
	}
	
	binary.Write(buf, binary.LittleEndian, memInfo.Total)
	binary.Write(buf, binary.LittleEndian, memInfo.Available)
	binary.Write(buf, binary.LittleEndian, memInfo.Used)
	binary.Write(buf, binary.LittleEndian, memInfo.Free)
	binary.Write(buf, binary.LittleEndian, memInfo.UsedPercent)
	
	// Add swap information
	swapInfo, err := mem.SwapMemory()
	if err == nil {
		binary.Write(buf, binary.LittleEndian, swapInfo.Total)
		binary.Write(buf, binary.LittleEndian, swapInfo.Used)
		binary.Write(buf, binary.LittleEndian, swapInfo.Free)
		binary.Write(buf, binary.LittleEndian, swapInfo.UsedPercent)
	} else {
		binary.Write(buf, binary.LittleEndian, uint64(0))
		binary.Write(buf, binary.LittleEndian, uint64(0))
		binary.Write(buf, binary.LittleEndian, uint64(0))
		binary.Write(buf, binary.LittleEndian, float64(0))
	}
	
	return writeSuccess(buf.Bytes()), nil
}

// getCPUUsage returns CPU usage information
func (s *SysHealthHandler) getCPUUsage() ([]byte, error) {
	buf := new(bytes.Buffer)
	
	// Get CPU usage percentage
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		return writeError(err), nil
	}
	
	if len(cpuPercent) > 0 {
		binary.Write(buf, binary.LittleEndian, cpuPercent[0])
	} else {
		binary.Write(buf, binary.LittleEndian, float64(0))
	}
	
	// Get per-CPU usage
	cpuPercentPerCPU, err := cpu.Percent(0, true)
	if err == nil {
		binary.Write(buf, binary.LittleEndian, uint32(len(cpuPercentPerCPU)))
		for _, usage := range cpuPercentPerCPU {
			binary.Write(buf, binary.LittleEndian, usage)
		}
	} else {
		binary.Write(buf, binary.LittleEndian, uint32(0))
	}
	
	return writeSuccess(buf.Bytes()), nil
}