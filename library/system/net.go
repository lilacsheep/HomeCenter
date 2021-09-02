package system

import "github.com/shirou/gopsutil/net"

func AllConnetions() ([]net.ConnectionStat, error) {
	return net.Connections("all")
}

func IOCounters(pernic bool) []net.IOCountersStat {
	info, _ := net.IOCounters(pernic)
	return info
}

// 获取网卡流量
func CompareNetFlow(t1 net.IOCountersStat, t2 net.IOCountersStat) (uint64, uint64) {
	return (t2.BytesRecv - t1.BytesRecv), (t2.BytesSent - t1.BytesSent)
}

// 当前信息与上次网卡信息做对比获取流量
func GetLatestNewFlow(t1 net.IOCountersStat) (uint64, uint64) {
	infos := IOCounters(true)
	for _, i := range infos {
		if i.Name == t1.Name {
			return (i.BytesRecv - t1.BytesRecv), (i.BytesSent - t1.BytesSent)
		}
	}
	return 0, 0
}
