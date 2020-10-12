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

func (self *ProxyInstanceApi) Update(r *ghttp.Request) {
	request := requests.NewUpdateInstanceRequest()
	self.DoRequestValid(request, r)
}

func (self *ProxyInstanceApi) Remove(r *ghttp.Request) {
	request := requests.NewRemoveInstanceRequest()
	self.DoRequestValid(request, r)
}

func (self *ProxyInstanceApi) RemoveFromPool(r *ghttp.Request) {
	request := requests.NewRemoveInstanceFromPoolRequest()
	self.DoRequestValid(request, r)
}

func (self *ProxyInstanceApi) IntoPool(r *ghttp.Request) {
	request := requests.NewAddInstanceIntoPoolRequest()
	self.DoRequestValid(request, r)
}
