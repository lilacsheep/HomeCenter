package aria2

import (
	"os/exec"
	"strings"

	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/util/gconv"
	"github.com/zyxar/argo/rpc"
)

var Manager *manager

type manager struct {
	Change bool
}

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
	status, err := server.TellStatus(gid)
	if err != nil {
		return err
	}
	if status.Status == "active" {
		server.Pause(gid)
	}

	for _, t := range status.Files {
		if gfile.Exists(t.Path) {
			gfile.Remove(t.Path)
		}
	}
	if force {
		_, err = server.ForceRemove(gid)
	} else {
		_, err = server.Remove(gid)
	}
	server.RemoveDownloadResult(gid)
	return err
}

func (self *manager) NewTask(s string) (string, error) {
	urls := strings.Split(s, "\n")
	return server.AddURI(urls)
}

func (self *manager) AddTorrent(filename string) error {
	_, err := server.AddTorrent(filename)
	return err
}

func (self *manager) TaskStatus(gid string) (info rpc.StatusInfo, err error) {
	return server.TellStatus(gid)
}

func (self *manager) GetGlobalOption() (options rpc.Option, err error) {
	return server.GetGlobalOption()
}

func (self *manager) Close() error {
	return server.Close()
}

func (self *manager) GetOption(key string) (string, error) {
	options, err := server.GetGlobalOption()
	if err != nil {
		return "", err
	}
	for k, v := range options {
		if k == key {
			return gconv.String(v), nil
		}
	}
	return "", err
}

func RestartAria2() {
	cmd := exec.Command("systemctl", "restart", "aria2")
	cmd.Run()
}
