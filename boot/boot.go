package boot

import (
	"homeproxy/app/models"
	"homeproxy/app/server"
	"homeproxy/app/server/aria2"
	"homeproxy/app/services/tasks"
	"homeproxy/library/filedb2"
	"time"

	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/glog"
)

func Setup() error {
	filedb2.Init()
	// 初始化用户
	defaultUser := models.User{Username: "admin", Password: "!QAZ2wsx", Status: true, CreateAt: time.Now()}
	if c, _ := filedb2.DB.Count(&models.User{});c == 0 {
		filedb2.DB.Save(&defaultUser)
	}
	// 初始化下载配置
	if found, _ := filedb2.DB.KeyExists("settings", "download"); !found {
		if !found {
			settings := models.DownloadSettings{}
			filedb2.DB.Set("settings", "download", &settings)
			glog.Info("init download settings")
		}
	}

	// 初始化文件管理节点
	count, _ := filedb2.DB.Count(&models.ProxyFileSystemNode{})
	if count == 0 {
		node := models.ProxyFileSystemNode{
			Path:     gfile.Abs("download/"),
			Name:     "下载",
			CreateAt: time.Now(),
		}
		filedb2.DB.Save(&node)
	}

	glog.Debugf("clean all monitor info")
	filedb2.DB.Drop(&models.ProxyMonitorInfo{})
	aria2.InitClient()

	server.Setup()
	tasks.InitDDnsTask()
	tasks.SetupMonitor()
	
	return nil
}