package models

import (
	"homeproxy/library/filedb2"
)


type ProxyMonitorInfo struct {
	ID          int     `json:"id" storm:"id,increment"`
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
	filedb2.DB.All(&info)
	return
}
