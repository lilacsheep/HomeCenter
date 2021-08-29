package boot

import (
	"homeproxy/app/models"
	"homeproxy/app/server"
	"homeproxy/app/server/aria2"
	"homeproxy/app/services/tasks"
	"homeproxy/library/filedb2"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcron"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/glog"
)

func Setup() error {
	filedb2.Init()

	sqlFile := []string{"dbsql/object_table.sql", "dbsql/global_config.sql", "dbsql/instances.sql", "dbsql/auth_users.sql"}
	for _, f := range sqlFile {
		s := gfile.GetContents(f)
		_, err := g.DB().Exec(s)
		if err != nil {
			return err
		}
	}

	// 初始化用户
	if c, _ := g.DB().Model(&models.User{}).Count(); c == 0 {
		models.CreateUser("admin", "!QAZ2wsx")
	}
	// 初始化下载配置

	downloadSettings, err := models.GetSettings()
	if err != nil {
		return err
	}
	if downloadSettings.Aria2Url != "" {
		gcron.AddSingleton("*/2 * * * * *", tasks.ReloadAira2Manager)
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
