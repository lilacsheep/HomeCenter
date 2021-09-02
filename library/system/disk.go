package system

import "github.com/shirou/gopsutil/disk"

func DiskPartitions() []disk.PartitionStat {
	info, _ := disk.Partitions(true)
	return info
}

func DiskUsage(path string) *disk.UsageStat {
	info, _ := disk.Usage(path)
	return info
}

func DiskIOCounters(names ...string) (map[string]disk.IOCountersStat, error){
	info, err := disk.IOCounters(names...)
	return info, err
}

func DiskAllUsage() []*disk.UsageStat {
	var info  []*disk.UsageStat
	for _, i := range DiskPartitions() {
		info = append(info, DiskUsage(i.Mountpoint))
	}
	return info
}