package requests

import (
	"homeproxy/library/system"

	"github.com/gogf/gf/net/ghttp"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

type systemInfo struct {
	Host     *host.InfoStat         `json:"host"`
	Mem      *mem.VirtualMemoryStat `json:"memory"`
	CpuTimes []cpu.TimesStat        `json:"cpu_times"`
	CpuInfo  []cpu.InfoStat         `json:"cpu_info"`
	NetCards []net.IOCountersStat   `json:"net_cards"`
	Disk     []*disk.UsageStat      `json:"disk"`
}

type SystemInfoRequest struct{}

func (req *SystemInfoRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	info := systemInfo{
		Host:     system.GetInfo(),
		Mem:      system.MemInfo(),
		CpuTimes: system.CpuTimes(false),
		CpuInfo:  system.CpuInfo(),
		NetCards: system.IOCounters(true),
		Disk:     system.DiskAllUsage(),
	}
	return *response.SuccessWithDetail(info)
}

type ProcessesRquest struct {}

func (req *ProcessesRquest) Exec(r *ghttp.Request) (response MessageResponse) {
	processes := system.Processes()
	return *response.SuccessWithDetail(processes)
}