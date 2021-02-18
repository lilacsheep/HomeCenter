package models

import (
	"homeproxy/app/server/download"
	"homeproxy/library/filedb2"
	"time"

	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gutil"
)

var (
	DownloadManager *downloadTaskManager
)

func init() {
	if found, _ := filedb2.DB.KeyExists("settings", "download"); !found {
		if !found {
			settings := DownloadSettings{
				Path:          gfile.Abs("download/"),
				ThreadNum:     32,
				NotifyOpen:    false,
				NotifyMessage: "",
			}
			filedb2.DB.Set("settings", "download", &settings)
			glog.Info("init download settings")
		}
	}

	if err := InitDownloadManager(); err != nil {
		panic(err)
	}
}

func InitDownloadManager() error {
	DownloadManager = NewDownloadTaskManager()
	DownloadManager.Settings = &DownloadSettings{}
	return DownloadManager.Init()
}

type DownloadSettings struct {
	ID            int    `json:"id" storm:"id,increment"`
	Path          string `json:"path"`           // 下载路径
	ThreadNum     int64  `json:"thread_num"`     // 默认的线程大小
	NotifyOpen    bool   `json:"notify_open"`    // 是否开启通知
	NotifyMessage string `json:"notify_message"` // 通知消息
	Aria2Enable   int64  `json:"aria2_enable"`   // 是否使用aria2 -- 1使用 --2不使用
	Aria2Url      string `json:"aria2_url"`      // aria2地址
	Aria2Token    string `json:"aria2_token"`    // aria2的Token
}

type downloadTaskManager struct {
	Settings *DownloadSettings
	tasks    *gmap.TreeMap
}

func (self *downloadTaskManager) Init() error {
	glog.Info("init download settings")
	err := filedb2.DB.Get("settings", "download", self.Settings)
	if err != nil {
		return err
	}
	glog.Info("init download settings done")

	if self.tasks.Size() == 0 {
		tasks, err := self.Tasks()
		if err != nil {
			return err
		} else {
			for _, task := range tasks {
				if task.Status == 3 {
					task.Status = 1
				}
				task.Init()
				self.tasks.Set(task.ID, task)
			}
		}
	}
	self.sync()
	return nil
}

func (self *downloadTaskManager) GetSettings() *DownloadSettings {
	return self.Settings
}

func (self *downloadTaskManager) UpdateSettings(data interface{}) error {
	new_ := gjson.New(data)
	self.Settings.Path = new_.GetString("path", self.Settings.Path)
	self.Settings.ThreadNum = new_.GetInt64("thread_num", self.Settings.ThreadNum)
	self.Settings.NotifyOpen = new_.GetBool("notify_open", self.Settings.NotifyOpen)
	self.Settings.NotifyMessage = new_.GetString("notify_message", self.Settings.NotifyMessage)
	return filedb2.DB.Set("settings", "download", self.Settings)
}

func (self *downloadTaskManager) NewTask(url string, threadNum int64, path string) (err error) {
	if threadNum == 0 {
		threadNum = self.Settings.ThreadNum
	}
	var (
		task  *download.Task
		path_ = self.Settings.Path
	)
	if path == "" {
		path_ = path
	}
	// 新建任务
	task, err = download.NewDownLoadTask(url, threadNum, path_)
	if err != nil {
		return
	}
	// 创建记录
	err = filedb2.DB.Save(task)
	if err != nil {
		return
	}
	task.Init()
	self.tasks.Set(task.ID, task)
	return
}

func (self *downloadTaskManager) Tasks() (tasks []*download.Task, err error) {
	err = filedb2.DB.All(&tasks)
	return
}

func (self *downloadTaskManager) StartTask(taskID int) {
	if value, found := self.tasks.Search(taskID); found {
		if task, ok := value.(*download.Task); ok {
			go task.Start()
		}
	}
}

func (self *downloadTaskManager) RemoveTask(taskID int) {
	if value := self.tasks.Remove(taskID); value != nil {
		task := value.(*download.Task)
		task.Cancel()
		filedb2.DB.DeleteStruct(task)
	} else {
		glog.Warning("task not exist")
	}
}

func (self *downloadTaskManager) CancelTask(taskID int) {
	if value, found := self.tasks.Search(taskID); found {
		value.(*download.Task).Cancel()
	}
}

func (self *downloadTaskManager) sync() {
	go func() {
		for {
			self.tasks.IteratorAsc(func(key, value interface{}) bool {
				var (
					new = value.(*download.Task)
					err error
				)
				err = filedb2.DB.Update(new)
				if err != nil {
					glog.Errorf("sync task: %d error: %s", key.(int), err.Error())
				}
				return true
			})
			time.Sleep(time.Millisecond * 500)
		}
	}()
}

func NewDownloadTaskManager() *downloadTaskManager {
	return &downloadTaskManager{
		tasks: gmap.NewTreeMap(gutil.ComparatorInt, true),
	}
}
