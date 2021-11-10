package boot

import (
	"homeproxy/app/models"
	"homeproxy/app/server/aria2"
	"homeproxy/app/services/tasks"
	"homeproxy/library/docker"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gfile"
)

func Setup() error {
	sqlFile := []string{
		"dbsql/objects.sql", "dbsql/global_config.sql", "dbsql/instances.sql", "dbsql/auth_users.sql",
		"dbsql/object_bucket.sql", "dbsql/object_token.sql", "dbsql/proxy_role.sql",
		"dbsql/ddns_operation_settings.sql", "dbsql/host.sql", "dbsql/host_group.sql",
		"dbsql/container_template.sql",
	}
	if gfile.Exists("dbsql") {
		for _, f := range sqlFile {
			s := gfile.GetContents(f)
			_, err := g.DB().Exec(s)
			if err != nil {
				return err
			}
		}
	}

	// 初始化用户
	if c, _ := g.DB().Model(&models.User{}).Count(); c == 0 {
		models.CreateUser("admin", "!QAZ2wsx")
	}

	// 初始化下载配置

	// downloadSettings, err := models.GetAria2Settings()
	// if err != nil {
	// 	return err
	// }
	// if downloadSettings.GetAria2Url() != "" {
	// 	gcron.AddSingleton("*/2 * * * * *", tasks.ReloadAira2Manager)
	// }
	docker.InitDockerClient()
	err := aria2.InitClient()
	if err != nil {
		return err
	}
	if minioOption, err := models.GetMinioConfig(); err != nil {
		return err
	} else {
		err = minioOption.Auto()
		if err != nil {
			return err
		}
	}
	// server.Setup()
	tasks.InitDDnsTask()
	return nil
}
