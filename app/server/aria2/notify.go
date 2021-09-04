package aria2

import (
	"homeproxy/app/models"
	"io/ioutil"
	"os"
	"path/filepath"

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
	settings, err := models.GetSettings()
	if err != nil {
		glog.Error("get settings error: %s", err.Error())
		return
	}
	if settings.AutoClean == 0 {
		return
	}
	cleanSize = settings.AutoClean * 1048576
	
	for _, event := range events {
		taskInfo, err := server.TellStatus(event.Gid)
		if err != nil {
			glog.Errorf("get task %s file error: %s", event.Gid, err.Error())
		} else {
			if taskInfo.BitTorrent.Info.Name != "" {
				if len(taskInfo.Files) > 1 {
					for _, info := range taskInfo.Files {
						if gconv.Int(info.Length) < cleanSize {
							glog.Infof("clean file %s size is %s b need %d b", info.Path, gconv.Int(info.Length), cleanSize)
							os.RemoveAll(info.Path)
						}
					}
				}
				dirPath := filepath.Join(taskInfo.Dir, taskInfo.BitTorrent.Info.Name)
				dirs, err := getDirList(dirPath)
				if err != nil {
					glog.Errorf("get dir path error: %s", err.Error())
				} else {
					for _, d := range dirs {
						t, _ := ioutil.ReadDir(d)
						if len(t) == 0 {
							os.RemoveAll(d)
						}
					}
				}
			}
		}
	}
}
func (CustomNotify) OnDownloadError(events []rpc.Event)      {}
func (CustomNotify) OnBtDownloadComplete(events []rpc.Event) {}


func getDirList(dirpath string) ([]string, error) {
	var dirs []string
	err := filepath.Walk(dirpath,
		func(path string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() {
				dirs = append(dirs, path)
				return nil
			}

			return nil
		})
	return dirs, err
}