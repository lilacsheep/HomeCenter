package tasks

import (
	"github.com/gogf/gf/os/gcron"
	"github.com/gogf/gf/util/gconv"
	"github.com/shirou/gopsutil/process"
	"homeproxy/app/models"
	"os"
	"runtime"
	"time"
)

var (
	windows = runtime.GOOS == "windows"
	History *HistoryInfo
)

func init() {
	gcron.AddSingleton("* * * * * *", QueryProxyMonitorInfoTask)
}

type HistoryInfo struct {
	ReadBytes  uint64 `json:"read_bytes"`
	WriteBytes uint64 `json:"write_bytes"`
	BytesSent  uint64 `json:"bytes_sent"`
	BytesRecv  uint64 `json:"bytes_recv"`
}

func QueryProxyMonitorInfoTask() {
	var (
		mainProxy *process.Process
		NowInfo   = &HistoryInfo{}
	)
	mainProxy, _ = process.NewProcess(gconv.Int32(os.Getpid()))

	if mainProxy != nil {
		c, _ := models.DB.Collection(models.ProxyMonitorTable)
		data := models.ProxyMonitorInfo{CreateAt: time.Now().Format("2006-01-02 15:04:05")}
		data.CpuPercent, _ = mainProxy.CPUPercent()
		if v, err := mainProxy.MemoryInfo(); err == nil {
			data.MemorySize = v.RSS
		}
		if v, err := mainProxy.Connections(); err == nil {
			data.Connections = len(v)
		}
		if v, err := mainProxy.IOCounters(); err == nil {
			NowInfo.WriteBytes = v.WriteBytes
			NowInfo.ReadBytes = v.ReadBytes
			if History != nil {
				data.WriteBytes = v.WriteBytes - History.WriteBytes
				data.ReadBytes = v.ReadBytes - History.ReadBytes
			}
		}
		if v, err := mainProxy.NetIOCounters(false); err == nil {
			NowInfo.BytesRecv = v[0].BytesRecv
			NowInfo.BytesSent = v[0].BytesSent
			if History != nil {

			}
		}
		if History != nil {
			data.BytesRecv = NowInfo.BytesRecv - History.BytesRecv
			data.BytesSent = NowInfo.BytesSent - History.BytesSent
			data.WriteBytes = NowInfo.WriteBytes - History.WriteBytes
			data.ReadBytes = NowInfo.ReadBytes - History.ReadBytes
			c.Insert(data)
			History = NowInfo
		} else {
			History = NowInfo
		}
	}
}
