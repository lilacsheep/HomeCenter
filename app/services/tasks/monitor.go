package tasks

import (
	"homeproxy/app/models"
	"homeproxy/library/filedb2"
	"os"
	"runtime"
	"time"

	"github.com/gogf/gf/os/gcron"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	"github.com/shirou/gopsutil/process"
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
		}
		if v, err := mainProxy.NetIOCounters(false); err == nil {
			NowInfo.BytesRecv = v[0].BytesRecv
			NowInfo.BytesSent = v[0].BytesSent
		}
		if History != nil {
			data.BytesRecv = NowInfo.BytesRecv - History.BytesRecv
			data.BytesSent = NowInfo.BytesSent - History.BytesSent
			data.WriteBytes = NowInfo.WriteBytes - History.WriteBytes
			data.ReadBytes = NowInfo.ReadBytes - History.ReadBytes
			err := filedb2.DB.Save(&data)
			if err != nil {
				glog.Errorf("save monitor info error: %s", err.Error())
			}
			History = NowInfo
		} else {
			History = NowInfo
		}
	}
}
