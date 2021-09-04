package tasks

import (
	"runtime"
)

var (
	windows = runtime.GOOS == "windows"
	History *HistoryInfo
)

func SetupMonitor() {
	// gcron.AddSingleton("* * * * * *", QueryProxyMonitorInfoTask)
}

type HistoryInfo struct {
	ReadBytes  uint64 `json:"read_bytes"`
	WriteBytes uint64 `json:"write_bytes"`
	BytesSent  uint64 `json:"bytes_sent"`
	BytesRecv  uint64 `json:"bytes_recv"`
}

func QueryProxyMonitorInfoTask() {

}
