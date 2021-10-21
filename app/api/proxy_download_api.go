package api

import (
	"homeproxy/app/server/aria2"
	"homeproxy/app/services/requests"
	"net/http"

	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
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
	request := requests.NewQueryWebSocketRequest()
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

func (self *ProxyDownloadApi) TaskStatus(r *ghttp.Request) {
	request := requests.NewTaskStatusRequest()
	self.DoRequestValid(request, r)
}

func (self *ProxyDownloadApi) Options(r *ghttp.Request) {
	request := requests.NewGetAria2GlobalOptionsRequest()
	self.DoRequestValid(request, r)
}

func (self *ProxyDownloadApi) MakeDownloadUrl(r *ghttp.Request) {
	request := requests.NewMakeAria2FileDownloadUrlRequest()
	self.DoRequestValid(request, r)
}

func (self *ProxyDownloadApi) Download(r *ghttp.Request) {
	type Aria2FileDownloadRequest struct {
		GID       string `json:"gid" v:"required"`
		FileIndex string `json:"file_index" v:"required"`
	}
	request := &Aria2FileDownloadRequest{}
	response := requests.MessageResponse{}

	if err := r.Parse(request); err != nil {
		response.ErrorWithMessage(requests.ParamsErrorCode, err.(gvalid.Error).Maps())
		r.Response.WriteJsonExit(response)
	}
	path, err := aria2.Manager.GetReadPath(request.GID, request.FileIndex)
	if err != nil {
		response.SystemError(err)
		r.Response.WriteJsonExit(response)
	}

	if path == "" {
		response.ErrorWithMessage(http.StatusNotFound, "文件不存在")
		r.Response.WriteJsonExit(request)
	}

	r.Response.ServeFileDownload(path)
}
