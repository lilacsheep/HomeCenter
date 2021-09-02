package system

import "github.com/shirou/gopsutil/mem"

func MemInfo() *mem.VirtualMemoryStat {
	info, _ := mem.VirtualMemory()
	return info
}
