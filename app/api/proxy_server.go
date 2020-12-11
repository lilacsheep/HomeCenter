package api

import (
	"github.com/gogf/gf/net/ghttp"
	"homeproxy/app/services/requests"
)

type ProxyServerApi struct {
	*BaseControllers
}

func (self *ProxyServerApi) Start(r *ghttp.Request) {
	request := requests.NewStartProxyServerRequest()
	self.DoRequest(request, r)
}

func (self *ProxyServerApi) Stop(r *ghttp.Request) {
	request := requests.NewStopProxyServerRequest()
	self.DoRequestValid(request, r)
}

func (self *ProxyServerApi) Update(r *ghttp.Request) {
	request := requests.NewUpdateProxyServerRequest()
	self.DoRequestValid(request, r)
}

func (self *ProxyServerApi) Info(r *ghttp.Request) {
	request := requests.NewInfoProxyServerRequest()
	self.DoRequestValid(request, r)
}

func (self *ProxyServerApi) Monitor(r *ghttp.Request) {
	request := requests.NewQueryProxyMonitorInfoRequest()
	self.DoRequest(request, r)
}

func (self *ProxyServerApi) Logs(r *ghttp.Request) {
	request := requests.NewQueryLogRequest()
	self.DoRequestValid(request, r)
}