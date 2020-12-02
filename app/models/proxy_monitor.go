package models

import (
	"github.com/gogf/gf/os/glog"
	"homeproxy/library/filedb"
)

const (
	ProxyMonitorTable = "proxy_monitor_info"
)

func init() {
	settings := filedb.DefaultCollectionSettings()
	settings.AutoDump = false
	settings.MaxRecord = 10
	if err := filedb.DB.NewCollections(ProxyMonitorTable, settings); err != nil {
		if err != filedb.ErrCollectionExist {
			glog.Error("init collection error: %s", err.Error())
		}
	}
}

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
	c, _ := filedb.DB.Collection(ProxyMonitorTable)
	err = c.All(&info)
	return
}
