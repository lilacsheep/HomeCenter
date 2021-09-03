package aria2

import (
	"context"
	"homeproxy/app/models"
	"time"

	"github.com/gogf/gf/os/glog"
	"github.com/zyxar/argo/rpc"
)

var server rpc.Client

const (
	ConfigPathKey      = "conf-path"
	ConfigBTTrackerKey = "bt-tracker"
)

func InitClient() error {
	settings, err := models.GetSettings()
	if err != nil {
		return err
	}
	if settings.Aria2Url != "" {
		server, err = rpc.New(context.Background(), settings.Aria2Url, settings.Aria2Token, time.Second, CustomNotify{})
		if err != nil {
			return err
		}
		Manager = &manager{}
		NewAutoUpdateBTTracker(settings.AutoUpdateBTTracker)
		glog.Info("aria2 connection success")
	} else {
		glog.Info("Aria2 not enabled")
	}
	return nil
}
