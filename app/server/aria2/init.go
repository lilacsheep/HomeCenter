package aria2

import (
	"context"
	"homeproxy/app/models"
	"homeproxy/library/filedb2"
	"time"

	"github.com/gogf/gf/os/glog"
	"github.com/zyxar/argo/rpc"
)


var Server rpc.Client

func init() {
	if err := InitClient(); err != nil {
		glog.Errorf("init aria2 client error: %s", err.Error())
	} else {
		options, _ := Server.GetGlobalOption()
		glog.Info(options)
	}
}


func InitClient() error {
	settings := models.DownloadSettings{}
	err := filedb2.DB.Get("settings", "download", &settings)
	if err != nil {
		return err
	}
	if settings.Aria2Enable {
		Server, err = rpc.New(context.Background(), settings.Aria2Url, settings.Aria2Token, time.Second, rpc.DummyNotifier{})
		if err != nil {
			return err
		}
		glog.Info("aria2 connection success")
	} else {
		glog.Info("Aria2 not enabled")
	}
	return nil
}