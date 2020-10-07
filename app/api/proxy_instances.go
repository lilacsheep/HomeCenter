package api

import (
	"github.com/gogf/gf/net/ghttp"
	"homeproxy/app/services/requests"
)

type ProxyInstanceApi struct {
	BaseControllers
}

func (self *ProxyInstanceApi) Create(r *ghttp.Request) {
	request := requests.NewCreateProxyInstanceRequest()
	self.DoRequestValid(request, r)
}

func (self *ProxyInstanceApi) Query(r *ghttp.Request) {
	request := requests.NewQueryAllInstanceRequest()
	self.DoRequest(request, r)
}
