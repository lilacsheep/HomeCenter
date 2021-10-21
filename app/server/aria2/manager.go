package aria2

import (
	"context"
	"fmt"
	"homeproxy/app/models"
	"homeproxy/library/docker"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	"github.com/zyxar/argo/rpc"
)

/*
-v 本地文件夹1:/downloads Aria2下载位置
-v 本地文件夹2:/config Aria2配置文件位置
-e PUID=1026 Linux用户UID
-e PGID=100 Linux用户GID
-e SECRET=yourtoken Aria2 token
-e CACHE=1024M Aria2磁盘缓存配置
-e PORT=6800 RPC通讯端口
-e WEBUI=true 启用WEBUI
-e WEBUI_PORT=8080 WEBUI端口
-e BTPORT=32516 DHT和BT监听端口
-e UT=true 启动容器时更新trackers
-e CTU= 启动容器时更新自定义trackes地址 https://cdn.jsdelivr.net/gh/XIU2/TrackersListCollection@master/best_aria2.txt
-e RUT=true	每天凌晨3点更新trackers
-e SMD=true 保存磁力链接为种子文件
-e FA= 磁盘预分配模式none,falloc,trunc,prealloc
-p 6800:6800 Aria2 RPC连接端口
-p 6881:6881 Aria2 tcp下载端口
-p 6881:6881/udp Aria2 p2p udp下载端口
--restart unless-stopped
*/
var (
	Manager *manager
	server  rpc.Client
)

const ImageName = "superng6/aria2:webui-latest"
const ContainerName = "aria2_auto"

type manager struct {
	Change   bool
	Settings *models.DownloadSettings
	init     bool
}

func (self *manager) Init() error {
	self.Settings, _ = models.GetAria2Settings()
	if !self.Settings.AutoStart {
		return nil
	}
	if self.Settings.ContainerId == "" {
		glog.Info("未发现容器ID, 重新创建Aria2容器...")
		ContainerId, err := self.CreateContainer()
		if err != nil {
			return err
		}
		self.Settings.ContainerId = ContainerId
		models.UpdateConfig("aria2", "container_id", ContainerId)
		err = docker.Docker.ContainerStart(context.Background(), ContainerId, types.ContainerStartOptions{})
		if err != nil {
			return err
		}
		server, err = rpc.New(context.Background(), fmt.Sprintf("http://127.0.0.1:%d/jsonrpc", self.Settings.Port), self.Settings.SECRET, time.Second, CustomNotify{})
		if err != nil {
			return err
		}
		self.init = true
	} else {
		info, err := docker.Docker.ContainerInspect(context.Background(), self.Settings.ContainerId)
		if err != nil {
			if !client.IsErrNotFound(err) {
				return err
			}
			glog.Infof("未发现容器ID：%s, 重新创建Aria2容器...", self.Settings.ContainerId)
			models.UpdateConfig("aria2", "container_id", "")
			return self.Init()
		}
		if !info.State.Running {
			err := docker.Docker.ContainerStart(context.Background(), self.Settings.ContainerId, types.ContainerStartOptions{})
			if err != nil {
				return err
			}
		}
		server, err = rpc.New(context.Background(), fmt.Sprintf("http://127.0.0.1:%d/jsonrpc", self.Settings.Port), self.Settings.SECRET, time.Second, CustomNotify{})
		if err != nil {
			return err
		}
		self.init = true
	}
	return nil
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

func (self *manager) AllTasks() (infos []rpc.StatusInfo, err error) {
	infos, err = self.ActiveTasks()
	if err != nil {
		return
	}
	stopped, err := self.TellStopped(0, 999)
	if err != nil {
		return
	}
	warnings, err := self.TellWaiting(0, 999)
	if err != nil {
		return
	}
	infos = append(infos, warnings...)
	infos = append(infos, stopped...)
	return
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

func (self *manager) StopContainer() error {
	return docker.Docker.ContainerStop(context.Background(), self.Settings.ContainerId, nil)
}

func (self *manager) StartContainer() error {
	return docker.Docker.ContainerStart(context.Background(), self.Settings.ContainerId, types.ContainerStartOptions{})
}

func (self *manager) CheckImage() (bool, error) {
	body, err := docker.Docker.ImageList(context.Background(), types.ImageListOptions{All: true})
	if err != nil {
		return false, err
	}
	for _, image := range body {
		if image.RepoTags[0] == ImageName {
			return true, nil
		}
	}
	return false, nil
}

func (self *manager) DownloadImage() error {
	out, err := docker.Docker.ImagePull(context.Background(), ImageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()
	var buff = make([]byte, 1024)
	for {
		_, err = out.Read(buff)
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		glog.Info(buff)
	}
	return nil
}

func (self *manager) CreateContainer() (string, error) {
	if v, err := self.CheckImage(); err != nil {
		return "", err
	} else {
		if !v {
			glog.Infof("未发现镜像：%s, 下载中...", ImageName)
			err = self.DownloadImage()
			if err != nil {
				return "", err
			}
		}
	}
	Config := container.Config{
		Image: ImageName,
		Env:   self.Settings.Env(),
	}
	HostConfig := self.Settings.ContainerHostConfig()
	NetworkConfig := network.NetworkingConfig{}

	body, err := docker.Docker.ContainerCreate(context.Background(), &Config, HostConfig, &NetworkConfig, nil, ContainerName)
	if err != nil {
		return "", err
	}
	return body.ID, nil
}

func (self *manager) Size() string {
	p := gfile.Abs(self.Settings.DownloadPath)
	return gfile.ReadableSize(p)
}

func (self *manager) GetReadPath(gid, index string) (string, error) {
	info, err := self.TaskStatus(gid)
	if err != nil {
		return "", err
	}

	var path string
	for _, i := range info.Files {
		if i.Index == index {
			path = i.Path
		}
	}

	p := strings.Replace(path, "/downloads", gfile.Abs(self.Settings.DownloadPath), -1)

	return p, nil
}
