package tasks

import (
	"homeproxy/library/system"
	"runtime"
	"sync"
	"time"

	"github.com/gogf/gf/os/gcron"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

var (
	windows = runtime.GOOS == "windows"
	History *SystemInfo
)

func init() {
	History = &SystemInfo{
		Host:     system.GetInfo(),
		Mem:      system.MemInfo(),
		CpuTimes: system.CpuTimes(false),
		CpuInfo:  system.CpuInfo(),
		NetCards: system.IOCounters(true),
		Disk:     system.DiskAllUsage(),
		Datetime: time.Now(),
		l:        sync.Mutex{},
	}
}

type SystemInfo struct {
	Host       *host.InfoStat         `json:"host"`
	Mem        *mem.VirtualMemoryStat `json:"memory"`
	CpuTimes   []cpu.TimesStat        `json:"cpu_times"`
	CpuInfo    []cpu.InfoStat         `json:"cpu_info"`
	NetCards   []net.IOCountersStat   `json:"net_cards"`
	Disk       []*disk.UsageStat      `json:"disk"`
	CpuPercent float64                `json:"cpu_percent"`
	Datetime   time.Time              `json:"datetime"`
	l          sync.Mutex
}

func (s *SystemInfo) Sync() {
	s.l.Lock()
	defer s.l.Unlock()
	t2 := system.CpuTimes(false)
	s.Host = system.GetInfo()
	s.Mem = system.MemInfo()
	s.CpuInfo = system.CpuInfo()
	s.NetCards = system.IOCounters(true)
	s.Disk = system.DiskAllUsage()
	s.Datetime = time.Now()
	s.CpuPercent = system.CpuPercent(t2[0], s.CpuTimes[0])
	s.CpuTimes = t2
}

func (s *SystemInfo) Clone() *SystemInfo {
	s.l.Lock()
	defer s.l.Unlock()
	return s
}


func SetupMonitor() {
	gcron.AddSingleton("* * * * * *", QueryProxyMonitorInfoTask)
}

func QueryProxyMonitorInfoTask() {
	History.Sync()
}
