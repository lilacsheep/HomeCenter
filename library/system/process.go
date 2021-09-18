package system

import (
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
)

type Process struct {
	Pid         int32                   `json:"pid"`
	Name        string                  `json:"name"`
	Status      string                  `json:"status"`
	CpuPercent  float64                 `json:"cpu_percent"`
	MemPercent  float32                 `json:"mem_percent"`
	Username    string                  `json:"username"`
	Connections []net.ConnectionStat    `json:"connections"`
	Cmdline     string                  `json:"cmdline"`
	IOCounters  *process.IOCountersStat `json:"io"`
	MemInfo     *process.MemoryInfoStat `json:"mem_info"`
	Cwd         string                  `json:"cwd"`
	Env         []string                `json:"env"`
	CreateTime  int64                   `json:"create_time"`
	Fds         []process.OpenFilesStat `json:"fds"`
}

func Processes() []Process {
	processes, _ := process.Processes()
	var all []Process
	for _, process := range processes {
		status, _ := process.Status()
		cp, _ := process.CPUPercent()
		mp, _ := process.MemoryPercent()
		name, _ := process.Name()
		all = append(all, Process{
			Pid:        process.Pid,
			Name:       name,
			Status:     status,
			CpuPercent: cp,
			MemPercent: mp,
		})
	}
	return all
}

func GetProcess(pid int32) (*process.Process, error) {
	return process.NewProcess(pid)
}

func ProcessInfo(pid int32) (*Process, error) {
	p, err := GetProcess(pid)
	if err != nil {
		return nil, err
	}
	info := Process{}
	info.Name, _ = p.Name()
	info.Pid = p.Pid
	info.Cwd, _ = p.Cwd()
	info.CpuPercent, _ = p.CPUPercent()
	info.Username, _ = p.Username()
	info.Status, _ = p.Status()
	info.Connections, _ = p.Connections()
	info.Cmdline, _ = p.Cmdline()
	info.CreateTime, _ = p.CreateTime()
	info.Fds, _ = p.OpenFiles()
	info.IOCounters, _ = p.IOCounters()
	info.MemInfo, _ = p.MemoryInfo()
	info.MemPercent, _ = p.MemoryPercent()
	info.Env, _ = p.Environ()
	return &info, nil
}
