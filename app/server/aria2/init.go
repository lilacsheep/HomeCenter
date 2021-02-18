package aria2

import (
	"homeproxy/app/models"
	"homeproxy/library/filedb2"

	"github.com/gogf/gf/os/glog"
	"github.com/zyxar/argo/rpc"
)


var Server *rpc.Client

func init() {
	if err := InitClient(); err != nil {
		glog.Errorf("init aria2 client error: %s", err.Error())
	}
}


func InitClient() error {
	settings := models.DownloadSettings{}
	err := filedb2.DB.Get("settings", "download", &settings)
	return err
}