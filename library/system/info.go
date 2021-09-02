package system

import "github.com/shirou/gopsutil/host"

func GetInfo() *host.InfoStat {
	info, _ := host.Info()
	return info
}
