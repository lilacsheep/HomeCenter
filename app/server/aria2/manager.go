package aria2

import (
	"strings"
	"github.com/zyxar/argo/rpc"
)


type manager struct{}

func (self *manager) ActiveTasks() (infos []rpc.StatusInfo, err error) {
	return Server.TellActive()
}

func (self *manager) StoppedTasks(offset, limit int) (infos []rpc.StatusInfo, err error) {
	return Server.TellStopped(offset, limit)
}

func (self *manager) NewTask(s string) error {
	urls := strings.Split(s, "\n")
	_, err := Server.AddURI(urls)
	return err
}