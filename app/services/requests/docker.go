package requests

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"homeproxy/library/docker"
	"io/ioutil"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/go-connections/nat"
	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gfile"
)

// CreateContainerRequest 容器创建请求
type CreateContainerRequest struct {
	Name   string
	Config struct {
		Env   []string
		Cmd   []string
		Image string `v:"required"`
	} `json:"config"`
	HostConfig struct {
		RestartPolicy   container.RestartPolicy
		Mounts          []mount.Mount
		PublishAllPorts bool
		PortBindings    map[string]interface{}
		Resource        container.Resources
	} `json:"host_config"`
}

func (req *CreateContainerRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	var (
		config     = container.Config{}
		hostConfig = container.HostConfig{PublishAllPorts: req.HostConfig.PublishAllPorts}
	)
	config.Image = req.Config.Image
	if len(req.Config.Env) != 0 {
		config.Env = req.Config.Env
	}
	hostConfig.RestartPolicy = req.HostConfig.RestartPolicy
	if len(req.HostConfig.Mounts) != 0 {
		for _, m := range req.HostConfig.Mounts {
			if !gfile.Exists(m.Source) {
				gfile.Mkdir(m.Source)
			}
		}
		hostConfig.Mounts = req.HostConfig.Mounts
	}
	if len(req.Config.Cmd) != 0 {
		config.Cmd = req.Config.Cmd
	}
	if len(req.HostConfig.PortBindings) != 0 {
		hostConfig.PortBindings = make(nat.PortMap, len(req.HostConfig.PortBindings))
		for k, v := range req.HostConfig.PortBindings {
			bind := nat.PortBinding{}
			err := gvar.New(v).Scan(&bind)
			if err != nil {
				return *response.SystemError(err)
			}
			p := nat.Port(k)
			hostConfig.PortBindings[p] = append(hostConfig.PortBindings[p], bind)
		}
	}
	body, err := docker.Docker.ContainerCreate(context.Background(), &config, &hostConfig, nil, nil, req.Name)
	if err != nil {
		return *response.SystemError(err)
	}
	return *response.SuccessWithDetail(body.ID)
}

// ContainersListRequest 容器清单
type ContainersListRequest struct {
	types.ContainerListOptions
}

// Exec 执行
func (req *ContainersListRequest) Exec(r *ghttp.Request) MessageResponse {
	response := MessageResponse{}
	containers, err := docker.Docker.ContainerList(context.Background(), req.ContainerListOptions)
	if err != nil {
		return *response.SystemError(err)
	}
	return *response.SuccessWithDetail(containers)
}

type ContainerStartRequest struct {
	ContainerID string `json:"id"`
	types.ContainerStartOptions
}

func (req ContainerStartRequest) Exec(r *ghttp.Request) MessageResponse {
	response := MessageResponse{}
	err := docker.Docker.ContainerStart(context.Background(), req.ContainerID, req.ContainerStartOptions)
	if err != nil {
		return *response.SystemError(err)
	}
	return *response.Success()
}

type ContainerStopRequest struct {
	ContainerID string `json:"id"`
}

func (req ContainerStopRequest) Exec(r *ghttp.Request) MessageResponse {
	response := MessageResponse{}
	err := docker.Docker.ContainerStop(context.Background(), req.ContainerID, nil)
	if err != nil {
		return *response.SystemError(err)
	}
	return *response.Success()
}

type ContainerPauseRequest struct {
	ContainerID string `json:"id"`
}

func (req ContainerPauseRequest) Exec(r *ghttp.Request) MessageResponse {
	response := MessageResponse{}
	err := docker.Docker.ContainerPause(context.Background(), req.ContainerID)
	if err != nil {
		return *response.SystemError(err)
	}
	return *response.Success()
}

type ContainerUnpauseRequest struct {
	ContainerID string `json:"id"`
}

func (req ContainerUnpauseRequest) Exec(r *ghttp.Request) MessageResponse {
	response := MessageResponse{}
	err := docker.Docker.ContainerUnpause(context.Background(), req.ContainerID)
	if err != nil {
		return *response.SystemError(err)
	}
	return *response.Success()
}

type ContainerUpdateRequest struct {
	ContainerID string `json:"id"`
	container.UpdateConfig
}

func (req ContainerUpdateRequest) Exec(r *ghttp.Request) MessageResponse {
	response := MessageResponse{}
	_, err := docker.Docker.ContainerUpdate(context.Background(), req.ContainerID, req.UpdateConfig)
	if err != nil {
		return *response.SystemError(err)
	}
	return *response.Success()
}

type ContainerRemoveRequest struct {
	ContainerID string `json:"id"`
	types.ContainerRemoveOptions
}

func (req ContainerRemoveRequest) Exec(r *ghttp.Request) MessageResponse {
	response := MessageResponse{}
	err := docker.Docker.ContainerRemove(context.Background(), req.ContainerID, req.ContainerRemoveOptions)
	if err != nil {
		return *response.SystemError(err)
	}
	return *response.Success()
}

type ContainerRestartRequest struct {
	ContainerID string `json:"id"`
}

func (req ContainerRestartRequest) Exec(r *ghttp.Request) MessageResponse {
	response := MessageResponse{}
	err := docker.Docker.ContainerRestart(context.Background(), req.ContainerID, nil)
	if err != nil {
		return *response.SystemError(err)
	}
	return *response.Success()
}

type ContainerInfoRequest struct {
	ContainerID string `json:"id"`
}

func (req ContainerInfoRequest) Exec(r *ghttp.Request) MessageResponse {
	response := MessageResponse{}
	info, err := docker.Docker.ContainerInspect(context.Background(), req.ContainerID)
	if err != nil {
		return *response.SystemError(err)
	}
	return *response.SuccessWithDetail(info)
}

type ContainerStatsRequest struct {
	ContainerID string `json:"id"`
}

func (req ContainerStatsRequest) Exec(r *ghttp.Request) MessageResponse {
	response := MessageResponse{}
	info, err := docker.Docker.ContainerStats(context.Background(), req.ContainerID, false)
	if err != nil {
		return *response.SystemError(err)
	} else {
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(info.Body)
		var s map[string]interface{}
		_ = json.Unmarshal(buf.Bytes(), &s)
		return *response.SuccessWithDetail(s)
	}
}

type ContainerRenameRequest struct {
	ContainerId string `json:"id"`
	Name        string `json:"name"`
}

func (req ContainerRenameRequest) Exec(r *ghttp.Request) MessageResponse {
	response := MessageResponse{}
	err := docker.Docker.ContainerRename(context.Background(), req.ContainerId, req.Name)
	if err != nil {
		return *response.SystemError(err)
	}
	return *response.Success()
}

type ListImagesRequest struct {
	types.ImageListOptions
}

func (req ListImagesRequest) Exec(r *ghttp.Request) MessageResponse {
	response := MessageResponse{}
	images, err := docker.Docker.ImageList(context.Background(), req.ImageListOptions)
	if err != nil {
		return *response.SystemError(err)
	}
	return *response.SuccessWithDetail(images)
}

type ImageHistoryRequest struct {
	ImageId string `json:"id"`
}

func (req ImageHistoryRequest) Exec(r *ghttp.Request) MessageResponse {
	response := MessageResponse{}
	body, err := docker.Docker.ImageHistory(context.Background(), req.ImageId)
	if err != nil {
		return *response.SystemError(err)
	}
	return *response.SuccessWithDetail(body)
}

type ImageInspectRequest struct {
	ImageId string `json:"id"`
}

func (req ImageInspectRequest) Exec(r *ghttp.Request) MessageResponse {
	response := MessageResponse{}
	body, _, err := docker.Docker.ImageInspectWithRaw(context.Background(), req.ImageId)
	if err != nil {
		return *response.SystemError(err)
	}
	return *response.SuccessWithDetail(body)
}

type DeleteImageRequest struct {
	ImageId string `json:"id"`
	types.ImageRemoveOptions
}

func (req DeleteImageRequest) Exec(r *ghttp.Request) MessageResponse {
	response := MessageResponse{}
	_, err := docker.Docker.ImageRemove(context.Background(), req.ImageId, req.ImageRemoveOptions)
	if err != nil {
		return *response.SystemError(err)
	}
	return *response.Success()
}

type PullImageRequest struct {
	Ref      string `json:"ref"`
	Username string `json:"username"`
	Password string `json:"password"`
	//types.ImagePullOptions
}

func (req PullImageRequest) Exec(r *ghttp.Request) MessageResponse {
	response := MessageResponse{}
	option := types.ImagePullOptions{}
	if req.Username != "" && req.Password != "" {
		auth := types.AuthConfig{Username: req.Username, Password: req.Password}
		authB, _ := json.Marshal(auth)
		option.RegistryAuth = base64.URLEncoding.EncodeToString(authB)
	}
	resp, err := docker.Docker.ImagePull(context.Background(), req.Ref, option)
	if err != nil {
		return *response.SystemError(err)
	}
	m, _ := ioutil.ReadAll(resp)
	return *response.SuccessWithDetail(m)
}

type VolumeListRequest struct {
	filters.Args
}

func (req VolumeListRequest) Exec(r *ghttp.Request) MessageResponse {
	response := MessageResponse{}

	body, err := docker.Docker.VolumeList(context.Background(), req.Args)
	if err != nil {
		return *response.SystemError(err)
	}
	return *response.SuccessWithDetail(body.Volumes)
}

type VolumeCreateRequest struct {
	volume.VolumeCreateBody
}

func (req VolumeCreateRequest) Exec(r *ghttp.Request) MessageResponse {
	response := MessageResponse{}
	_, err := docker.Docker.VolumeCreate(context.Background(), req.VolumeCreateBody)
	if err != nil {
		return *response.SystemError(err)
	}
	return *response.Success()
}

type VolumeRemoveRequest struct {
	ID    string `json:"volume_id"`
	Force bool   `json:"force"`
}

func (req VolumeRemoveRequest) Exec(r *ghttp.Request) MessageResponse {
	response := MessageResponse{}
	err := docker.Docker.VolumeRemove(context.Background(), req.ID, req.Force)
	if err != nil {
		return *response.SystemError(err)
	}
	return *response.Success()
}
