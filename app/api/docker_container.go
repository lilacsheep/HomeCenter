package api

import (
	"homeproxy/app/services/requests"

	"github.com/gogf/gf/net/ghttp"
)

type ContainersController struct {
	BaseControllers
}

func (c *ContainersController) List(r *ghttp.Request) {
	request := requests.ContainersListRequest{}
	c.DoRequestValid(&request, r)
}

func (c *ContainersController) Start(r *ghttp.Request) {
	request := requests.ContainerStartRequest{}
	c.DoRequestValid(&request, r)
}

func (c *ContainersController) Stop(r *ghttp.Request) {
	request := requests.ContainerStopRequest{}
	c.DoRequestValid(&request, r)
}

func (c *ContainersController) Restart(r *ghttp.Request) {
	request := requests.ContainerStopRequest{}
	c.DoRequestValid(&request, r)
}

func (c *ContainersController) Update(r *ghttp.Request) {
	request := requests.ContainerUpdateRequest{}
	c.DoRequestValid(&request, r)
}

func (c *ContainersController) Pause(r *ghttp.Request) {
	request := requests.ContainerPauseRequest{}
	c.DoRequestValid(&request, r)
}

func (c *ContainersController) Unpause(r *ghttp.Request) {
	request := requests.ContainerUnpauseRequest{}
	c.DoRequestValid(&request, r)
}

func (c *ContainersController) Info(r *ghttp.Request) {
	request := requests.ContainerInfoRequest{}
	c.DoRequestValid(&request, r)
}

func (c *ContainersController) Stats(r *ghttp.Request) {
	request := requests.ContainerStatsRequest{}
	c.DoRequestValid(&request, r)
}

func (c *ContainersController) Rename(r *ghttp.Request) {
	request := requests.ContainerRenameRequest{}
	c.DoRequestValid(&request, r)
}
