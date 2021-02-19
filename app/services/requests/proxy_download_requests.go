package requests

import (
	"homeproxy/app/models"
	"homeproxy/app/server/aria2"
	"net/http"
	"path"

	"github.com/gogf/gf/net/ghttp"
	"github.com/zyxar/argo/rpc"
)

type CreateDownloadTaskRequest struct {
	Url       string `json:"url"`
	Path      string `json:"path"`
	ThreadNum int64  `json:"thread_num"`
}

func (self *CreateDownloadTaskRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	err := aria2.Manager.NewTask(self.Url)
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		response.Success()
	}
	return
}

type CreateTorrentDownloadRequest struct {}

func (self *CreateTorrentDownloadRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	file := r.GetUploadFile("file")
	filename, err := file.Save("/tmp")
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		err := aria2.Manager.AddTorrent(path.Join("/tmp", filename))
		if err != nil {
			response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
		} else {
			response.Success()
		}
	}
	return
}

func NewCreateTorrentDownloadRequest() *CreateTorrentDownloadRequest {
	return &CreateTorrentDownloadRequest{}
}

func NewCreateDownloadTaskRequest() *CreateDownloadTaskRequest {
	return &CreateDownloadTaskRequest{}
}

type QueryDownloadTaskRequest struct {
	Status string `json:"status"`
}

func (self *QueryDownloadTaskRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	var (
		tasks map[string][]rpc.StatusInfo
		err   error
	)
	tasks = make(map[string][]rpc.StatusInfo)

	tasks["active"], err = aria2.Manager.ActiveTasks()
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		tasks["stopped"], err = aria2.Manager.TellStopped(0, 999)
		if err != nil {
			response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
		} else {
			response.DataTable(tasks, len(tasks))
		}
	}
	return
}

func NewQueryDownloadTaskRequest() *QueryDownloadTaskRequest {
	return &QueryDownloadTaskRequest{}
}

type RemoveDownloadTaskRequest struct {
	Gid string `json:"id"`
}

func (self *RemoveDownloadTaskRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	err := aria2.Manager.RemoveTask(self.Gid, false)
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		response.Success()
	}
	return
}

func NewRemoveDownloadTaskRequest() *RemoveDownloadTaskRequest {
	return &RemoveDownloadTaskRequest{}
}

type CancelDownloadTaskRequest struct {
	TaskId int `json:"id"`
}

func (self *CancelDownloadTaskRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	models.DownloadManager.CancelTask(self.TaskId)
	response.Success()
	return
}

func NewCancelDownloadTaskRequest() *CancelDownloadTaskRequest {
	return &CancelDownloadTaskRequest{}
}

type StartDownloadTaskRequest struct {
	TaskId int `json:"id"`
}

func (self *StartDownloadTaskRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	models.DownloadManager.StartTask(self.TaskId)
	response.Success()
	return
}

func NewStartDownloadTaskRequest() *StartDownloadTaskRequest {
	return &StartDownloadTaskRequest{}
}

type GetDownloadSettingsRequest struct{}

func (self *GetDownloadSettingsRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	settings := models.DownloadManager.GetSettings()
	response.SuccessWithDetail(settings)
	return
}

func NewGetDownloadSettingsRequest() *GetDownloadSettingsRequest {
	return &GetDownloadSettingsRequest{}
}

type UpdateDownloadSettingsRequest struct {
	Path          string `json:"path"`       // 下载路径
	ThreadNum     int64  `json:"thread_num"` // 默认的线程大小
	NotifyOpen    bool   `json:"notify_open"`
	NotifyMessage string `json:"notify_message"`
	Aria2Enable   bool   `json:"aria2_enable"`
	Aria2Url      string `json:"aria2_url"`
	Aria2Token    string `json:"aria2_token"`
}

func (self *UpdateDownloadSettingsRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	err := models.DownloadManager.UpdateSettings(self)
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		response.Success()
	}
	return
}

func NewUpdateDownloadSettingsRequest() *UpdateDownloadSettingsRequest {
	return &UpdateDownloadSettingsRequest{}
}
