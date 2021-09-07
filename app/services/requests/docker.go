package requests

import (
	"context"
	"homeproxy/library/docker"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/gogf/gf/net/ghttp"
)

// CreateContainerRequest 容器创建请求
type CreateContainerRequest struct {
	Name          string
	Config        container.Config         `json:"config"`
	HostConfig    container.HostConfig     `json:"host_config"`
	NetworkConfig network.NetworkingConfig `json:"network_config"`
}

func (req *CreateContainerRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	body, err := docker.Docker.ContainerCreate(context.Background(), &req.Config, &req.HostConfig, &req.NetworkConfig, nil, req.Name)
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
