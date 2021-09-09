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
