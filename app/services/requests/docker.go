package requests

import (
	"context"
	"homeproxy/library/docker"

	"github.com/gogf/gf/net/ghttp"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
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
