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
)

func init() {
	gcron.AddSingleton("* * * * * *", QueryProxyMonitorInfoTask)
}

func QueryProxyMonitorInfoTask() {
	var (
		mainProxy *process.Process
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
			data.WriteBytes = v.WriteBytes
			data.ReadBytes = v.ReadBytes
		}
		if v, err := mainProxy.NetIOCounters(false); err == nil {
			data.BytesRecv = v[0].BytesRecv
			data.BytesSent = v[0].BytesSent
		}

		c.Insert(data)
	}
}
