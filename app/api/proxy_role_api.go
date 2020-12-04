package api

import (
	"github.com/gogf/gf/net/ghttp"
	"homeproxy/app/services/requests"
)

type ProxyRoleApi struct {
	*BaseControllers
}

func (self *ProxyRoleApi) AddRole(r *ghttp.Request) {
	request := requests.NewAddUrlRoleRequest()
	self.DoRequestValid(request, r)
}

func (self *ProxyRoleApi) Remove(r *ghttp.Request) {
	request := requests.NewRemoveUrlRoleRequest()
	self.DoRequestValid(request, r)
}

func (self *ProxyRoleApi) Change(r *ghttp.Request) {
	request := requests.NewChangeRoleInstanceRequest()
	self.DoRequestValid(request, r)
}

func (self *ProxyRoleApi) All(r *ghttp.Request) {
	request := requests.NewQueryAllRoleRequest()
	self.DoRequest(request, r)
}

func (self *ProxyRoleApi) RoleVisitLogs(r *ghttp.Request) {
	request := requests.NewGetRoleAllVisitRequest()
	self.DoRequest(request, r)
}

func (self *ProxyRoleApi) AddLogToRole(r *ghttp.Request) {
	request := requests.NewAddErrorSiteToProxyRequest()
	self.DoRequestValid(request, r)
}
