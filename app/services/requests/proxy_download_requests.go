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
	Url string `json:"url"`
}

func (self *CreateDownloadTaskRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	gid, err := aria2.Manager.NewTask(self.Url)
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		response.SuccessWithDetail(gid)
	}
	return
}

type CreateTorrentDownloadRequest struct{}

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
		infos []rpc.StatusInfo
		err   error
	)

	actives, err := aria2.Manager.ActiveTasks()
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
		return
	}
	stopped, err := aria2.Manager.TellStopped(0, 999)
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
		return
	}
	warnings, err := aria2.Manager.TellWaiting(0, 999)
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
		return
	}
	for _, task := range actives {
		if len(task.FollowedBy) > 0 {
			infos = append(infos, task)
		}
	}
	for _, task := range warnings {
		if len(task.FollowedBy) > 0 {
			infos = append(infos, task)
		}
	}
	for _, task := range stopped {
		if len(task.FollowedBy) > 0 {
			infos = append(infos, task)
		}
	}
	response.DataTable(infos, len(infos))
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

type PauseDownloadTaskRequest struct {
	TaskId string `json:"id"`
}

func (self *PauseDownloadTaskRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	err := aria2.Manager.PauseTask(self.TaskId, false)
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		response.Success()
	}
	return
}

func NewPauseDownloadTaskRequest() *PauseDownloadTaskRequest {
	return &PauseDownloadTaskRequest{}
}

type UnpauseDownloadTaskRequest struct {
	TaskId string `json:"id"`
}

func (self *UnpauseDownloadTaskRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	err := aria2.Manager.UnpauseTask(self.TaskId)
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		response.Success()
	}
	return
}

func NewUnpauseDownloadTaskRequest() *UnpauseDownloadTaskRequest {
	return &UnpauseDownloadTaskRequest{}
}

type GetDownloadSettingsRequest struct{}

func (self *GetDownloadSettingsRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	settings, err := models.GetSettings()
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		response.SuccessWithDetail(settings)
	}
	return
}

func NewGetDownloadSettingsRequest() *GetDownloadSettingsRequest {
	return &GetDownloadSettingsRequest{}
}

type UpdateDownloadSettingsRequest struct {
	Aria2Url            string `json:"aria2_url"`
	Aria2Token          string `json:"aria2_token"`
	AutoClean           int    `json:"auto_clean"`
	AutoUpdateBTTracker string `json:"auto_update_bt_tracker"`
}

func (self *UpdateDownloadSettingsRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	err := aria2.UpdateSettings(self)
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

type GlobalStatInfoRequest struct{}

func (self *GlobalStatInfoRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	info, err := aria2.Manager.GetGlobalStat()
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		response.SuccessWithDetail(info)
	}
	return
}

func NewGlobalStatInfoRequest() *GlobalStatInfoRequest {
	return &GlobalStatInfoRequest{}
}

type TaskStatusRequest struct {
	Gid string `json:"id"`
}

func (self *TaskStatusRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	info, err := aria2.Manager.TaskStatus(self.Gid)
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		response.SuccessWithDetail(info)
	}
	return
}

func NewTaskStatusRequest() *TaskStatusRequest {
	return &TaskStatusRequest{}
}

type GetAria2GlobalOptionsRequest struct{}

func (self *GetAria2GlobalOptionsRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	opts := []map[string]interface{}{}
	options, err := aria2.Manager.GetGlobalOption()
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		for k, v := range options {
			opts = append(opts, map[string]interface{}{"key": k, "value": v})
		}
		response.SuccessWithDetail(opts)
	}
	return
}

func NewGetAria2GlobalOptionsRequest() *GetAria2GlobalOptionsRequest {
	return &GetAria2GlobalOptionsRequest{}
}
