package api

import (
	"github.com/gogf/gf/net/ghttp"
	"homeproxy/app/services/requests"
)

type ProxyDDnsApi struct {
	BaseControllers
}

func (self *ProxyDDnsApi) ChooseCards(r *ghttp.Request) {
	request := requests.NewChooseNetCardRequest()
	self.DoRequest(request, r)
}

func (self *ProxyDDnsApi) GetSettings(r *ghttp.Request) {
	request := requests.NewGetDDnsSettingsRequest()
	self.DoRequest(request, r)
}

func (self *ProxyDDnsApi) CreateSetting(r *ghttp.Request) {
	request := requests.NewCreateDDnsSettingRequest()
	self.DoRequestValid(request, r)
}

func (self *ProxyDDnsApi) GetRecords(r *ghttp.Request) {
	request := requests.NewGetProviderRecordsRequest()
	self.DoRequestValid(request, r)
}

func (self *ProxyDDnsApi) GetSettingInfo(r *ghttp.Request) {
	request := requests.NewGetSettingsInfoRequest()
	self.DoRequestValid(request, r)
}

func (self *ProxyDDnsApi) StartSetting(r *ghttp.Request) {
	request := requests.NewStartRoleTaskRequest()
	self.DoRequestValid(request, r)
}

func (self *ProxyDDnsApi) StopSetting(r *ghttp.Request) {
	request := requests.NewStopRoleTaskRequest()
	self.DoRequestValid(request, r)
}

func (self *ProxyDDnsApi) DeleteSetting(r *ghttp.Request) {
	request := requests.NewDeleteSettingRequest()
	self.DoRequestValid(request, r)
}

func (self *ProxyDDnsApi) RefreshSetting(r *ghttp.Request) {
	request := requests.NewRefreshSettingRequest()
	self.DoRequestValid(request, r)
}
