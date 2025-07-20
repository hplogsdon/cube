package worker

import (
	"github.com/c9s/goprocinfo/linux"
	"log"
)

type Stats struct {
	MemStats  *linux.MemInfo
	DiskStats *linux.Disk
	CPUStats  *linux.CPUStat
	LoadStats *linux.LoadAvg
}

func (s *Stats) MemUsedKb() uint64 {
	return s.MemTotalKb() - s.MemAvailableKb()
}

func (s *Stats) MemUsedPercent() uint64 {
	return s.MemAvailableKb() / s.MemTotalKb()
}

func (s *Stats) MemTotalKb() uint64 {
	return s.MemStats.MemTotal
}

func (s *Stats) MemAvailableKb() uint64 {
	return s.MemStats.MemAvailable
}

func (s *Stats) DiskTotal() uint64 {
	return s.DiskStats.All
}

func (s *Stats) DiskFree() uint64 {
	return s.DiskStats.Free
}

func (s *Stats) DiskUsed() uint64 {
	return s.DiskStats.Used
}

func (s *Stats) CPUUsage() float64 {
	idle := s.CPUStats.Idle + s.CPUStats.IOWait
	nonIdle := s.CPUStats.User + s.CPUStats.Nice + s.CPUStats.System + s.CPUStats.IRQ + s.CPUStats.SoftIRQ + s.CPUStats.Steal
	total := idle + nonIdle

	if total == 0 && idle == 0 {
		return 0.00
	}
	return float64(total) - float64(idle)/float64(total)
}

func GetStats() *Stats {
	return &Stats{
		MemStats:  GetMemoryInfo(),
		DiskStats: GetDiskInfo(),
		CPUStats:  GetCPUInfo(),
		LoadStats: GetLoadInfo(),
	}
}

func GetMemoryInfo() *linux.MemInfo {
	memStats, err := linux.ReadMemInfo("/proc/meminfo")
	if err != nil {
		log.Printf("Error reading memory info: %s", err)
		return &linux.MemInfo{}
	}
	return memStats
}

func GetDiskInfo() *linux.Disk {
	diskStats, err := linux.ReadDisk("/")
	if err != nil {
		log.Printf("Error reading disk stats: %s", err)
		return &linux.Disk{}
	}
	return diskStats
}

func GetCPUInfo() *linux.CPUStat {
	cpuStats, err := linux.ReadStat("/proc/stat")
	if err != nil {
		log.Printf("Error reading cpu stats: %s", err)
		return &linux.CPUStat{}
	}
	return &cpuStats.CPUStatAll
}

func GetLoadInfo() *linux.LoadAvg {
	loadAvg, err := linux.ReadLoadAvg("/proc/loadavg")
	if err != nil {
		log.Printf("Error reading load avg: %s", err)
		return &linux.LoadAvg{}
	}
	return loadAvg
}
