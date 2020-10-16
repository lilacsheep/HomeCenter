package models

const (
	ProxyMonitorTable = "proxy_monitor_info"
)

type ProxyMonitorInfo struct {
	ID          string  `json:"id"`
	CpuPercent  float64 `json:"cpu_percent"`
	MemorySize  uint64  `json:"memory_size"`
	ReadBytes   uint64  `json:"read_bytes"`
	WriteBytes  uint64  `json:"write_bytes"`
	Connections int     `json:"connections"`
	BytesSent   uint64  `json:"bytes_sent"`
	BytesRecv   uint64  `json:"bytes_recv"`
	CreateAt    string  `json:"create_at"`
}

func GetAllProxyMonitorInfo() (info []ProxyMonitorInfo, err error) {
	c, _ := DB.Collection(ProxyMonitorTable)
	err = c.All(&info)
	return
}
