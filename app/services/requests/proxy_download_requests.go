package requests

import (
	"homeproxy/app/models"
	"homeproxy/app/server/aria2"
	"net/http"
	"path"

	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
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

type QueryWebSocketRequest struct{}

func (self *QueryWebSocketRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	var (
		infos []rpc.StatusInfo
		err   error
	)
	ws, err := r.WebSocket()
	if err != nil {
		return *response.SystemError(err)
	}
	defer ws.Close()

	for {
		_, d, err := ws.ReadMessage()
		if err != nil {
			return *response.SystemError(err)
		}
		j := gjson.New(g.Map{})
		switch gvar.New(d).String() {
		case "tasks":
			infos, _ = aria2.Manager.AllTasks()
			j.Set("type", "tasks")
			j.Set("data", infos)
			ws.WriteJSON(j.Map())
		case "stats":
			stats, _ := aria2.Manager.GetGlobalStat()
			j.Set("type", "stats")
			j.Set("data", stats)
			ws.WriteJSON(j.Map())
		}
	}
}

func NewQueryWebSocketRequest() *QueryWebSocketRequest {
	return &QueryWebSocketRequest{}
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
	settings, err := models.GetAria2Settings()
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
	Port         int    `json:"port"`          // RPC连接端口
	TcpPort      int    `json:"tcp_port"`      // tcp下载端口
	UdpPort      int    `json:"udp_port"`      // Aria2 p2p udp下载端口
	DownloadPath string `json:"download_path"` // Aria2下载位置
	ConfigPath   string `json:"config_path"`   // Aria2配置文件位置
	SECRET       string `json:"token"`         // Aria2 token
	WebUi        bool   `json:"webui"`         // 启用WEBUI
	WebUiPort    int    `json:"webui_port"`    // WEBUI端口
	BTPort       int    `json:"bt_port"`       // DHT和BT监听端口
	UT           bool   `json:"ut"`            // 启动容器时更新trackers
	CTU          string `json:"ctu"`           // 启动容器时更新自定义trackes地址
	RUT          bool   `json:"rut"`           // 每天凌晨3点更新trackers
	SMD          bool   `json:"smd"`           // 保存磁力链接为种子文件
	FA           string `json:"fa"`            // 磁盘预分配模式none,falloc,trunc,prealloc
	AutoClean    int    `json:"auto_clean"`    // 自动清理Bt下载后文件夹内内容，根据文件大小判断
	AutoStart    bool   `json:"auto_start"`    // 启动Aria2任务
	PublicVisit  bool   `json:"public_visit"`  // 允许公共访问
}

func (self *UpdateDownloadSettingsRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	settings := models.DownloadSettings{}
	var needstart bool
	err := models.ConfigToStruct("aria2", &settings)
	if err != nil {
		return *response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		if settings.NeedReCreate(gvar.New(self)) {
			aria2.Manager.StopContainer()
			aria2.Manager.RemoveContianer()
			needstart = true
		}

		for k, v := range gvar.New(self).MapStrStr() {
			models.UpdateConfig("aria2", k, v)
		}

		if self.AutoStart != settings.AutoStart {
			if self.AutoStart {
				err = aria2.Manager.Init()
				if err != nil {
					return *response.SystemError(err)
				}
				needstart = false
			} else {
				aria2.Manager.StopContainer()
			}
		}
		if needstart && self.AutoStart {
			err = aria2.Manager.Init()
			if err != nil {
				return *response.SystemError(err)
			}
		}

	}
	return *response.Success()
}

func NewUpdateDownloadSettingsRequest() *UpdateDownloadSettingsRequest {
	return &UpdateDownloadSettingsRequest{}
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

type MakeAria2FileDownloadUrlRequest struct {
	GID       string `json:"gid"`
	FileIndex string `json:"file_index"`
}

func (self *MakeAria2FileDownloadUrlRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	// task, err := aria2.Manager.TaskStatus(self.GID)
	// if err != nil {
	// 	response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	// 	return
	// }
	// var download = &models.DownloadFileList{}

	// for _, i := range task.Files {
	// 	if i.Index == self.FileIndex {
	// 		download.Path = i.Path
	// 		download.CreateAt = time.Now()
	// 		download.Vkey = guid.S()
	// 		filedb2.DB.Save(download)
	// 	}
	// }
	// if download.Vkey == "" {
	// 	response.ErrorWithMessage(http.StatusInternalServerError, "资源不存在")
	// } else {
	// 	response.SuccessWithDetail(download.Vkey)
	// }
	return
}

func NewMakeAria2FileDownloadUrlRequest() *MakeAria2FileDownloadUrlRequest {
	return &MakeAria2FileDownloadUrlRequest{}
}
