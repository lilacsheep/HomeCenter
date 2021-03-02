package aria2

import (
	"homeproxy/app/models"
	"homeproxy/library/filedb2"
	"os"

	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	"github.com/zyxar/argo/rpc"
)


type CustomNotify struct{}

func (CustomNotify) OnDownloadStart(events []rpc.Event)  {}
func (CustomNotify) OnDownloadPause(events []rpc.Event)      {}
func (CustomNotify) OnDownloadStop(events []rpc.Event)       {}
func (CustomNotify) OnDownloadComplete(events []rpc.Event)   {
	var cleanSize int
	settings := models.DownloadSettings{}
	err := filedb2.DB.Get("settings", "download", &settings)
	if err != nil {
		glog.Error("get settings error: %s", err.Error())
		return
	}
	if settings.AutoClean == 0 {
		return
	}
	cleanSize = settings.AutoClean * 1048576
	for _, event := range events {
		infos, err := server.GetFiles(event.Gid)
		if err != nil {
			glog.Error("get task %s file error: %s", event.Gid, err.Error())
		} else {
			if len(infos) > 1 {
				for _, info := range infos {
					if gconv.Int(info.Length) < cleanSize {
						glog.Info("clean file %s size is %s b need %d b", info.Path, gconv.Int(info.Length), cleanSize)
						os.RemoveAll(info.Path)
					}
				}
			}
		}
	}

}
func (CustomNotify) OnDownloadError(events []rpc.Event)      {}
func (CustomNotify) OnBtDownloadComplete(events []rpc.Event) {}