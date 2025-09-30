package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

// HwInfoHandler handles hardware information commands
type HwInfoHandler struct{}

// NewHwInfoHandler creates a new hardware info handler
func NewHwInfoHandler() *HwInfoHandler {
	return &HwInfoHandler{}
}

// Handle processes hardware info commands
func (h *HwInfoHandler) Handle(data []byte) ([]byte, error) {
	if len(data) < 4 {
		return writeError(fmt.Errorf("invalid hwinfo command data")), nil
	}

	subCommand := binary.LittleEndian.Uint32(data[:4])

	switch subCommand {
	case HWINFO_CPU:
		return h.getCPUInfo()
	case HWINFO_MEMORY:
		return h.getMemoryInfo()
	case HWINFO_DISK:
		return h.getDiskInfo()
	case HWINFO_ALL:
		return h.getAllHardwareInfo()
	default:
		buf := new(bytes.Buffer)
		writeString(buf, fmt.Sprintf("Hardware info command %d not yet implemented", subCommand))
		return writeSuccess(buf.Bytes()), nil
	}
}

// getCPUInfo returns detailed CPU information
func (h *HwInfoHandler) getCPUInfo() ([]byte, error) {
	buf := new(bytes.Buffer)
	
	cpuInfo, err := cpu.Info()
	if err != nil {
		return writeError(err), nil
	}
	
	binary.Write(buf, binary.LittleEndian, uint32(len(cpuInfo)))
	for _, cpu := range cpuInfo {
		writeString(buf, cpu.ModelName)
		writeString(buf, cpu.VendorID)
		binary.Write(buf, binary.LittleEndian, uint32(cpu.Cores))
		binary.Write(buf, binary.LittleEndian, cpu.Mhz)
		binary.Write(buf, binary.LittleEndian, uint32(cpu.CacheSize))
	}
	
	return writeSuccess(buf.Bytes()), nil
}

// getMemoryInfo returns memory information
func (h *HwInfoHandler) getMemoryInfo() ([]byte, error) {
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
	
	return writeSuccess(buf.Bytes()), nil
}

// getDiskInfo returns disk information
func (h *HwInfoHandler) getDiskInfo() ([]byte, error) {
	buf := new(bytes.Buffer)
	
	partitions, err := disk.Partitions(false)
	if err != nil {
		return writeError(err), nil
	}
	
	binary.Write(buf, binary.LittleEndian, uint32(len(partitions)))
	for _, partition := range partitions {
		writeString(buf, partition.Device)
		writeString(buf, partition.Mountpoint)
		writeString(buf, partition.Fstype)
		
		usage, err := disk.Usage(partition.Mountpoint)
		if err == nil {
			binary.Write(buf, binary.LittleEndian, usage.Total)
			binary.Write(buf, binary.LittleEndian, usage.Used)
			binary.Write(buf, binary.LittleEndian, usage.Free)
			binary.Write(buf, binary.LittleEndian, usage.UsedPercent)
		} else {
			binary.Write(buf, binary.LittleEndian, uint64(0))
			binary.Write(buf, binary.LittleEndian, uint64(0))
			binary.Write(buf, binary.LittleEndian, uint64(0))
			binary.Write(buf, binary.LittleEndian, float64(0))
		}
	}
	
	return writeSuccess(buf.Bytes()), nil
}

// getAllHardwareInfo returns comprehensive hardware information
func (h *HwInfoHandler) getAllHardwareInfo() ([]byte, error) {
	buf := new(bytes.Buffer)
	
	// Get CPU info
	cpuData, err := h.getCPUInfo()
	if err != nil {
		cpuData = writeError(err)
	}
	binary.Write(buf, binary.LittleEndian, uint32(len(cpuData)))
	buf.Write(cpuData)
	
	// Get Memory info
	memData, err := h.getMemoryInfo()
	if err != nil {
		memData = writeError(err)
	}
	binary.Write(buf, binary.LittleEndian, uint32(len(memData)))
	buf.Write(memData)
	
	// Get Disk info
	diskData, err := h.getDiskInfo()
	if err != nil {
		diskData = writeError(err)
	}
	binary.Write(buf, binary.LittleEndian, uint32(len(diskData)))
	buf.Write(diskData)
	
	return writeSuccess(buf.Bytes()), nil
}