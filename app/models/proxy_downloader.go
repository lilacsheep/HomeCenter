package models

import (
	"homeproxy/app/server/download"
	"homeproxy/library/filedb"
	"homeproxy/library/filedb2"
	"time"

	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gutil"
)

const (
	DownloadListTable     = "proxy_download_list"
	DownloadSettingsTable = "proxy_download_settings"
)

var (
	DownloadManager    *downloadTaskManager
	SettingsCollection *filedb.Collection
	TasksCollection    *filedb.Collection
)

func init() {
	settings := filedb.DefaultCollectionSettings()
	settings.Unique = "url"
	if err := filedb.DB.NewCollections(DownloadListTable, settings); err != nil {
		if err != filedb.ErrCollectionExist {
			glog.Error("init collection error: %s", err.Error())
		}
	}
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
	var err error
	TasksCollection, err = filedb.DB.Collection(DownloadListTable)
	if err != nil {
		return err
	}
	DownloadManager = NewDownloadTaskManager()
	DownloadManager.Settings = &DownloadSettings{}
	return DownloadManager.Init()
}

type DownloadSettings struct {
	ID            string `json:"id"`
	Path          string `json:"path"`           // 下载路径
	ThreadNum     int64  `json:"thread_num"`     // 默认的线程大小
	NotifyOpen    bool   `json:"notify_open"`    // 是否开启通知
	NotifyMessage string `json:"notify_message"` // 通知消息
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
				if task.Status == 2 {
					task.Status = 0
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
	local := gjson.New(self.Settings)
	new_ := gjson.New(data)
	_ = local.Remove("id")
	local.Set("path", new_.Get("path", local.Get("path")))
	local.Set("thread_num", new_.Get("thread_num", local.Get("thread_num")))
	local.Set("notify_open", new_.Get("notify_open", local.Get("notify_open")))
	local.Set("notify_message", new_.Get("notify_message", local.Get("notify_message")))
	if err := SettingsCollection.UpdateById(self.Settings.ID, local.ToMap()); err != nil {
		return err
	}
	return SettingsCollection.GetFirst(self.Settings)
}

func (self *downloadTaskManager) NewTask(url string, threadNum int64, path string) (err error) {
	if threadNum == 0 {
		threadNum = self.Settings.ThreadNum
	}
	var (
		task  *download.Task
		id    string
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
	id, err = TasksCollection.Insert(task)
	if err != nil {
		return
	}
	task.ID = id
	task.Init()
	self.tasks.Set(id, task)
	return
}

func (self *downloadTaskManager) Tasks() (tasks []*download.Task, err error) {
	err = TasksCollection.All(&tasks)
	return
}

func (self *downloadTaskManager) StartTask(taskID string) {
	if value, found := self.tasks.Search(taskID); found {
		if task, ok := value.(*download.Task); ok {
			go task.Start()
		}
	}
}

func (self *downloadTaskManager) RemoveTask(taskID string) {
	if value := self.tasks.Remove(taskID); value != nil {
		value.(*download.Task).Cancel()
		TasksCollection.RemoveById(taskID)
	} else {
		glog.Warning("task not exist")
	}
}

func (self *downloadTaskManager) CancelTask(taskID string) {
	if value, found := self.tasks.Search(taskID); found {
		value.(*download.Task).Cancel()
	}
}

func (self *downloadTaskManager) sync() {
	convertTask := func(task *download.Task) *gjson.Json {
		j := gjson.New(task, true)
		_ = j.Remove("url")
		_ = j.Remove("id")
		return j
	}

	go func() {
		for {
			self.tasks.IteratorAsc(func(key, value interface{}) bool {
				task := value.(*download.Task)
				j := convertTask(task)
				err := TasksCollection.UpdateById(task.ID, j.ToMap())
				if err != nil {
					glog.Errorf("sync task: %s error: %s", task.ID, err.Error())
				}
				return true
			})
			time.Sleep(time.Second)
		}
	}()
}

func NewDownloadTaskManager() *downloadTaskManager {
	return &downloadTaskManager{
		tasks: gmap.NewTreeMap(gutil.ComparatorString, true),
	}
}
