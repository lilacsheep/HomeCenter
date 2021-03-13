package api

import (
	"homeproxy/app/models"
	"homeproxy/app/services/requests"
	"homeproxy/library/filedb2"
	"net/http"

	"github.com/asdine/storm/v3"
	"github.com/gogf/gf/net/ghttp"
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

func (self *ProxyDownloadApi) Options(r *ghttp.Request) {
	request := requests.NewGetAria2GlobalOptionsRequest()
	self.DoRequestValid(request, r)
}


func (self *ProxyDownloadApi) MakeDownloadUrl(r *ghttp.Request) {
	request := requests.NewMakeAria2FileDownloadUrlRequest()
	self.DoRequestValid(request, r)
}

func (self *ProxyDownloadApi) Download(r *ghttp.Request) {
	vkey := r.GetString("vkey", "")
	if vkey == "" {
		r.Response.WriteStatus(http.StatusNotFound)
	} else {

		info := models.DownloadFileList{}
		err := filedb2.DB.One("Vkey", vkey, &info)
		if err != nil {
			if err == storm.ErrNotFound {
				r.Response.WriteStatus(http.StatusNotFound)
			} else {
				r.Response.WriteStatus(http.StatusInternalServerError)
			}
		} else {
			r.Response.ServeFileDownload(info.Path)
		}
	}
}