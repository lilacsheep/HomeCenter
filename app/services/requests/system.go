package requests

import (
	"homeproxy/app/services/tasks"
	"homeproxy/library/system"

	"github.com/gogf/gf/net/ghttp"
)

type SystemInfoRequest struct{}

func (req *SystemInfoRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	return *response.SuccessWithDetail(tasks.History.Clone())
}

type ProcessesRquest struct{}

func (req *ProcessesRquest) Exec(r *ghttp.Request) (response MessageResponse) {
	processes := system.Processes()
	return *response.SuccessWithDetail(processes)
}

type ProcessInfoRequest struct {
	Pid int32
}

func (req *ProcessInfoRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	process, err := system.ProcessInfo(req.Pid)
	if err != nil {
		return *response.SystemError(err)
	}
	return *response.SuccessWithDetail(process)
}
