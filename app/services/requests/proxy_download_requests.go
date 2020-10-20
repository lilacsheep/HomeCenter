package requests

import (
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"homeproxy/app/models"
	"homeproxy/app/server/download"
	"net/http"
	"time"
)

type CreateDownloadTaskRequest struct {
	Url       string `json:"url"`
	ChunkSize int    `json:"chunk_size"`
}

func (self *CreateDownloadTaskRequest) sync(task *download.Task) {
	go func() {
		c, _ := models.DB.Collection(models.DownloadListTable)
		convertTask := func(task *download.Task) *gjson.Json {
			j := gjson.New(task, true)
			j.Remove("url")
			j.Remove("id")
			return j
		}
		for task.Status < 3 {
			j := convertTask(task)
			err := c.UpdateById(task.ID, j.ToMap())
			if err != nil {
				glog.Error(err.Error())
			}
			time.Sleep(time.Second)
		}
		j := convertTask(task)
		err := c.UpdateById(task.ID, j.ToMap())
		if err != nil {
			glog.Error(err.Error())
		}
		glog.Debugf("task: %s sync done", task.ID)
	}()
}
func (self *CreateDownloadTaskRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	if self.ChunkSize == 0 {
		self.ChunkSize = 32 * 1024 * 1024
	}
	task, err := download.NewDownLoadTask(self.Url, self.ChunkSize)
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		go task.Do()
		time.Sleep(time.Second)

		c, _ := models.DB.Collection(models.DownloadListTable)
		id, err := c.Insert(task)
		if err != nil {
			response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
		} else {
			task.ID = id
			self.sync(task)
			response.SuccessWithDetail(id)
		}
	}
	return
}

func NewCreateDownloadTaskRequest() *CreateDownloadTaskRequest {
	return &CreateDownloadTaskRequest{}
}

type QueryDownloadTaskRequest struct{}

func (self *QueryDownloadTaskRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	tasks, err := models.GetALLDownloadTasks()
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		response.DataTable(tasks, len(tasks))
	}
	return
}

func NewQueryDownloadTaskRequest() *QueryDownloadTaskRequest {
	return &QueryDownloadTaskRequest{}
}
