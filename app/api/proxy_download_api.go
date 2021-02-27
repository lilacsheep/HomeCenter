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

func (self *ProxyDownloadApi) AddTorrent(r *ghttp.Request) {
	request := requests.NewCreateTorrentDownloadRequest()
	self.DoRequest(request, r)
}

func (self *ProxyDownloadApi) Query(r *ghttp.Request) {
	request := requests.NewQueryDownloadTaskRequest()
	self.DoRequestValid(request, r)
}

func (self *ProxyDownloadApi) Remove(r *ghttp.Request) {
	request := requests.NewRemoveDownloadTaskRequest()
	self.DoRequestValid(request, r)
}

func (self *ProxyDownloadApi) Pause(r *ghttp.Request) {
	request := requests.NewPauseDownloadTaskRequest()
	self.DoRequestValid(request, r)
}

func (self *ProxyDownloadApi) UnPause(r *ghttp.Request) {
	request := requests.NewUnpauseDownloadTaskRequest()
	self.DoRequestValid(request, r)
}

func (self *ProxyDownloadApi) Settings(r *ghttp.Request) {
	request := requests.NewGetDownloadSettingsRequest()
	self.DoRequest(request, r)
}

func (self *ProxyDownloadApi) UpdateSettings(r *ghttp.Request) {
	request := requests.NewUpdateDownloadSettingsRequest()
	self.DoRequestValid(request, r)
}


func (self *ProxyDownloadApi) GlobalStatInfo(r *ghttp.Request) {
	request := requests.NewGlobalStatInfoRequest()
	self.DoRequestValid(request, r)
}

func (self *ProxyDownloadApi) TaskStatus(r *ghttp.Request) {
	request := requests.NewTaskStatusRequest()
	self.DoRequestValid(request, r)
}