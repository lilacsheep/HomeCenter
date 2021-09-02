package api

import (
	"homeproxy/app/services/requests"

	"github.com/gogf/gf/net/ghttp"
)

type SystemApi struct {
	BaseControllers
}

func (a *SystemApi) Info(r *ghttp.Request) {
	request := &requests.SystemInfoRequest{}
	a.DoRequest(request, r)
}

func (a *SystemApi) Processes(r *ghttp.Request) {
	request := &requests.ProcessesRquest{}
	a.DoRequest(request, r)
}