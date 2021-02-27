package aria2

import (
	"strings"
	"github.com/zyxar/argo/rpc"
)

var Manager *manager

type manager struct{}

func (self *manager) GetGlobalStat() (info rpc.GlobalStatInfo, err error) {
	return server.GetGlobalStat()
}

func (self *manager) ActiveTasks() (infos []rpc.StatusInfo, err error) {
	return server.TellActive()
}

func (self *manager) TellStopped(offset, limit int) (infos []rpc.StatusInfo, err error) {
	return server.TellStopped(offset, limit)
}

func (self *manager) TellWaiting(offset, limit int) (infos []rpc.StatusInfo, err error) {
	return server.TellWaiting(offset, limit)
}

func (self *manager) UnpauseTask(gid string) error {
	_, err := server.Unpause(gid)
	return err
}

func (self *manager) PauseTask(gid string, force bool) (err error) {
	if force {
		_, err = server.Pause(gid)
	} else {
		_, err = server.ForcePause(gid)
	}
	return err
}

func (self *manager) RemoveTask(gid string, force bool) (err error) {
	if force {
		_, err = server.ForceRemove(gid)
	} else {
		_, err = server.Remove(gid)
	}
	server.RemoveDownloadResult(gid)
	return err
}

func (self *manager) NewTask(s string) error {
	urls := strings.Split(s, "\n")
	_, err := server.AddURI(urls)
	return err
}

func (self *manager) AddTorrent(filename string) error {
	_, err := server.AddTorrent(filename)
	return err
}

func (self *manager) TaskStatus(gid string) (info rpc.StatusInfo, err error) {
	return server.TellStatus(gid)
}