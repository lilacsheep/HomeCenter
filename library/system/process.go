package system

import "github.com/shirou/gopsutil/process"

type Process struct {
	Pid        int32  `json:"pid"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	CpuPercent float64 `json:"cpu_percent"`
	MemPercent float32 `json:"mem_percent"`
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
			Pid:    process.Pid,
			Name:   name,
			Status: status,
			CpuPercent: cp,
			MemPercent: mp,
		})
	}
	return all
}

func GetProcess(pid int32) (*process.Process, error) {
	return process.NewProcess(pid)
}