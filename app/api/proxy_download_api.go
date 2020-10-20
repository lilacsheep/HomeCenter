package api

import (
	"github.com/gogf/gf/net/ghttp"
	"homeproxy/app/services/requests"
)

type ProxyDownloadApi struct {
	BaseControllers
}

func (self *ProxyDownloadApi) Create(r *ghttp.Request) {
	request := requests.NewCreateDownloadTaskRequest()
	self.DoRequestValid(request, r)
}

func (self *ProxyDownloadApi) Query(r *ghttp.Request) {
	request := requests.NewQueryDownloadTaskRequest()
	self.DoRequestValid(request, r)
}
