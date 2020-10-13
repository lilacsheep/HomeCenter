package requests

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
	"github.com/shirou/gopsutil/process"
	"net/http"
	"os"
)

type QueryProxyMonitorInfoRequest struct{}

func (self *QueryProxyMonitorInfoRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	p, err := process.NewProcess(gconv.Int32(os.Getpid()))
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		data := g.Map{}
		data["pid"] = os.Getpid()
		data["cpu_percent"], _ = p.CPUPercent()
		data["memory"], _ = p.MemoryInfo()
		data["fd"], _ = p.NumFDs()
		data["connections"], _ = p.Connections()
		data["counters"], _ = p.IOCounters()
		data["NetIOCounters"], _ = p.NetIOCounters(false)
		response.SuccessWithDetail(data)
	}
	return
}

func NewQueryProxyMonitorInfoRequest() *QueryProxyMonitorInfoRequest {
	return &QueryProxyMonitorInfoRequest{}
}
