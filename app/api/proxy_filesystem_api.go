package api

import (
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
	"homeproxy/app/services/requests"
	"net/http"
)

type ProxyFilesystemApi struct {
	BaseControllers
}

func (self *ProxyFilesystemApi) Nodes(r *ghttp.Request) {
	request := requests.NewAllFileSystemNodesRequest()
	self.DoRequest(request, r)
}

func (self *ProxyFilesystemApi) Files(r *ghttp.Request) {
	request := requests.NewAllFilesystemFilesRequest()
	self.DoRequestValid(request, r)
}

func (self *ProxyFilesystemApi) FileInfo(r *ghttp.Request) {
	request := requests.NewGetFilesystemFileInfoRequest()
	self.DoRequestValid(request, r)
}

func (self *ProxyFilesystemApi) DownloadFile(r *ghttp.Request) {
	request := &requests.DownloadFilesystemFileRequest{}
	if err := r.Parse(request); err != nil {
		resp := requests.MessageResponse{}
		switch err.(type) {
			case gvalid.Error:
				resp.ErrorWithMessage(requests.ParamsErrorCode, err.(gvalid.Error).Maps())
				r.Response.WriteJsonExit(resp)
			default:
				resp.ErrorWithMessage(requests.ParamsErrorCode, err.Error())
				r.Response.WriteJsonExit(resp)
		}
	} else {
		response := request.Exec(r)
		if response.ErrorCode == http.StatusNotFound {
			r.Response.WriteJsonExit(response)
		}
	}
}

func (self *ProxyFilesystemApi) RemoveFile(r *ghttp.Request) {
	request := requests.NewRemoveFilesystemFileRequest()
	self.DoRequestValid(request, r)
}

func (self *ProxyFilesystemApi) UploadFile(r *ghttp.Request) {
	request := requests.NewUploadFilesystemFileRequest()
	self.DoRequestValid(request, r)
}

func (self *ProxyFilesystemApi) CreateNode(r *ghttp.Request) {
	request := requests.NewCreateFilesystemNodeRequest()
	self.DoRequestValid(request, r)
}

func (self *ProxyFilesystemApi) RemoveNode(r *ghttp.Request) {
	request := requests.NewRemoveFilesystemNodeRequest()
	self.DoRequestValid(request, r)
}

func (self *ProxyFilesystemApi) CreateDir(r *ghttp.Request) {
	request := requests.NewCreateFilesystemDirRequest()
	self.DoRequestValid(request, r)
}

func (self *ProxyFilesystemApi) RemoveDir(r *ghttp.Request) {
	request := requests.NewRemoveFilesystemDirRequest()
	self.DoRequestValid(request, r)
}
