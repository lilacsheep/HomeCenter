package models

import (
	"context"
	"errors"
	"fmt"
	"homeproxy/library/docker"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/grand"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const ImageName = "minio/minio:latest"
const ContainerName = "minio_auto"

type LocalMinioOption struct {
	Public        bool   `json:"public"`
	AutoStart     bool   `json:"auto_start"`
	Port          int    `json:"port"`
	AccessKey     string `json:"access_key"`
	SecretKey     string `json:"secret_key"`
	WebUi         bool   `json:"webui"`
	WebUiPort     int    `json:"webui_port"`
	SavePath      string `json:"save_path"`
	ConfigDir     string `json:"config_dir"`
	MinioDomain   string `json:"minio_domain"`
	Region        string `json:"region"`
	RegionComment string `json:"region_comment"`
	ContainerID   string `json:"container_id"`
}

func (o *LocalMinioOption) Auto() (err error) {
	err = ConfigToStruct("minio", o)
	if err != nil {
		return err
	}
	if !o.AutoStart {
		_, err = docker.Docker.ContainerInspect(context.Background(), o.ContainerID)
		if err != nil {
			if !client.IsErrNotFound(err) {
				return err
			} else {
				o.ContainerID = ""
				UpdateConfig("minio", "container_id", "")
				return nil
			}
		}
		o.StopContainer()
		o.RemoveContainer()
		return
	}
	if o.ContainerID != "" {
		_, err = docker.Docker.ContainerInspect(context.Background(), o.ContainerID)
		if err != nil {
			if !client.IsErrNotFound(err) {
				return err
			} else {
				o.ContainerID = ""
				UpdateConfig("minio", "container_id", "")
			}
		} else {
			o.StopContainer()
			o.RemoveContainer()
		}
	}
	option := filters.NewArgs(filters.KeyValuePair{Key: "name", Value: ContainerName})
	containers, err := docker.Docker.ContainerList(context.Background(), types.ContainerListOptions{
		All:     true,
		Filters: option,
	})
	if err != nil {
		return err
	}
	for _, c := range containers {
		glog.Infof("clean container: %s", c.ID)
		o.ContainerID = ""
		UpdateConfig("minio", "container_id", "")
		docker.Docker.ContainerRemove(context.Background(), c.ID, types.ContainerRemoveOptions{Force: true})
	}

	err = o.CreateContainer()
	if err != nil {
		return err
	}
	err = o.StartContainer()
	if err != nil {
		return err
	}
	return nil
}

func (o *LocalMinioOption) NeedReCreate(options *g.Var) bool {
	now := gvar.New(o).MapStrStr()
	new := options.MapStrStr()

	reCreateKeys := []string{"port", "access_key",
		"secret_key", "webui", "webui_port", "save_path", "config_dir"}

	for _, key := range reCreateKeys {
		newOption, _ := new[key]
		if now[key] != newOption {
			return true
		}
	}
	return false
}

func (o *LocalMinioOption) Env() []string {
	var envs []string
	if o.Region != "" {
		envs = append(envs, fmt.Sprintf("MINIO_REGION_NAME=%s", o.Region))
	}
	if o.RegionComment != "" {
		envs = append(envs, fmt.Sprintf("MINIO_REGION_COMMENT=%s", o.RegionComment))
	}
	envs = append(envs, fmt.Sprintf("MINIO_ROOT_USER=%s", o.AccessKey))
	envs = append(envs, fmt.Sprintf("MINIO_ROOT_PASSWORD=%s", o.SecretKey))
	if !o.WebUi {
		envs = append(envs, "MINIO_BROWSER=off")
	}
	if o.MinioDomain != "" {
		envs = append(envs, fmt.Sprintf("MINIO_ROOT_USER=%s", o.MinioDomain))
	}
	return envs
}

func (o *LocalMinioOption) ContainerHostConfig() *container.HostConfig {
	var host string

	if o.Public {
		host = "0.0.0.0"
	} else {
		host = "127.0.0.1"
	}

	portMap := nat.PortMap{
		nat.Port("9000/tcp"): []nat.PortBinding{{HostIP: host, HostPort: fmt.Sprintf("%d/tcp", o.Port)}},
	}
	if o.WebUi {
		portMap[nat.Port("9001/tcp")] = []nat.PortBinding{{HostIP: host, HostPort: fmt.Sprintf("%d/tcp", o.WebUiPort)}}
	}
	// downloadPath := gfile.Abs(o.SavePath)
	configPath := gfile.Abs(o.ConfigDir)
	// if !gfile.Exists(downloadPath) {
	// 	gfile.Mkdir(downloadPath)
	// }

	if !gfile.Exists(configPath) {
		gfile.Mkdir(configPath)
	}
	mt := []mount.Mount{
		{Type: "bind", Source: configPath, Target: "/root/.minio"},
	}
	for _, d := range strings.Split(o.SavePath, ",") {
		if d == "" {
			continue
		}
		dirPath := gfile.Abs(d)
		if !gfile.Exists(dirPath) {
			gfile.Mkdir(dirPath)
		}
		mt = append(mt, mount.Mount{
			Type:   "bind",
			Source: dirPath,
			Target: dirPath,
		})
	}
	HostConfig := container.HostConfig{
		PortBindings:  portMap,
		RestartPolicy: container.RestartPolicy{Name: "always"},
		Mounts:        mt,
	}
	return &HostConfig
}

func (o *LocalMinioOption) CreateContainer() error {
	o.reload()
	if o.ContainerID != "" {
		return errors.New(fmt.Sprintf("container id not none: %s", o.ContainerID))
	}
	if ok, err := checkMinioImages(); err != nil {
		return err
	} else {
		if !ok {
			err = downloadMinioImage()
			if err != nil {
				return err
			}
		}
	}
	var cmd = []string{"server"}
	for _, d := range strings.Split(o.SavePath, ",") {
		if d == "" {
			continue
		}
		cmd = append(cmd, d)
	}

	if o.WebUi {
		cmd = append(cmd, "--console-address")
		cmd = append(cmd, "0.0.0.0:9001")
	}

	HostConfig := o.ContainerHostConfig()
	ExposedPorts := make(nat.PortSet)

	for k := range HostConfig.PortBindings {
		ExposedPorts[k] = struct{}{}
	}
	Config := container.Config{
		Image:        ImageName,
		Env:          o.Env(),
		Cmd:          cmd,
		ExposedPorts: ExposedPorts,
	}

	NetworkConfig := network.NetworkingConfig{}

	body, err := docker.Docker.ContainerCreate(context.Background(), &Config, HostConfig, &NetworkConfig, nil, ContainerName)
	if err != nil {
		return err
	}
	o.ContainerID = body.ID
	UpdateConfig("minio", "container_id", body.ID)
	return nil
}

func (o *LocalMinioOption) StopContainer() error {
	o.reload()
	if o.ContainerID == "" {
		return nil
	}
	glog.Infof("stop container id: %s", o.ContainerID)
	return docker.Docker.ContainerStop(context.Background(), o.ContainerID, nil)
}

func (o *LocalMinioOption) RemoveContainer() error {
	o.reload()
	if o.ContainerID == "" {
		return nil
	}
	glog.Infof("remove container id: %s", o.ContainerID)
	err := docker.Docker.ContainerRemove(context.Background(), o.ContainerID, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   true,
		Force:         true,
	})
	if err != nil {
		return err
	}
	o.ContainerID = ""
	UpdateConfig("minio", "container_id", "")
	return nil
}

func (o *LocalMinioOption) StartContainer() error {
	o.reload()
	if o.ContainerID == "" {
		return nil
	}
	glog.Infof("start container id: %s", o.ContainerID)
	return docker.Docker.ContainerStart(context.Background(), o.ContainerID, types.ContainerStartOptions{})
}

func (o *LocalMinioOption) reload() {
	ConfigToStruct("minio", o)
}

func checkMinioImages() (bool, error) {
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

func downloadMinioImage() error {
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

func GetMinioConfig() (*LocalMinioOption, error) {
	options := &LocalMinioOption{}
	c, err := g.DB().Model(&GlobalConfig{}).Where("group", "minio").Count()
	if err != nil {
		return nil, err
	}
	if c == 0 {
		g.DB().Model(&GlobalConfig{}).Data(g.List{
			{"group": "minio", "key": "auto_start", "type": "bool", "value": "false", "desc": "开机自启动minio"},
			{"group": "minio", "key": "port", "type": "int", "value": "9000", "desc": "api端口"},
			{"group": "minio", "key": "access_key", "type": "string", "value": grand.S(16), "desc": "密钥账号"},
			{"group": "minio", "key": "secret_key", "type": "string", "value": grand.S(24), "desc": "密钥认证"},
			{"group": "minio", "key": "webui", "type": "bool", "value": "false", "desc": "web console 开启"},
			{"group": "minio", "key": "webui_port", "type": "int", "value": "9001", "desc": "web console端口"},
			{"group": "minio", "key": "public", "type": "bool", "value": "true", "desc": "允许远程访问"},
			{"group": "minio", "key": "save_path", "type": "string", "value": "/data/minio/data", "desc": "存储目录"},
			{"group": "minio", "key": "config_dir", "type": "string", "value": "/data/minio/config", "desc": "配置文件地址"},
			{"group": "minio", "key": "minio_domain", "type": "string", "value": "", "desc": "minio域名地址"},
			{"group": "minio", "key": "region", "type": "string", "value": "us-east-1", "desc": "区域"},
			{"group": "minio", "key": "region_comment", "type": "string", "value": "", "desc": "区域备注"},
			{"group": "minio", "key": "container_id", "type": "string", "value": "", "desc": "本机容器ID"},
		}).Save()
	}
	err = ConfigToStruct("minio", options)
	return options, err
}

func GetMinioClient() (*minio.Client, error) {
	options, err := GetMinioConfig()
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("127.0.0.1:%d", options.Port)

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(options.AccessKey, options.SecretKey, ""),
		Secure: false,
	})
	return minioClient, err
}
