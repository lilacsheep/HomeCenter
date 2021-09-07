package api

import (
	"homeproxy/app/services/requests"

	"github.com/gogf/gf/net/ghttp"
)

type VolumeApi struct {
	BaseControllers
}

func (c *VolumeApi) Create(r *ghttp.Request) {
	request := requests.VolumeCreateRequest{}
	c.DoRequestValid(&request, r)
}

func (c *VolumeApi) List(r *ghttp.Request) {
	request := requests.VolumeListRequest{}
	c.DoRequestValid(&request, r)
}

func (c *VolumeApi) Remove(r *ghttp.Request) {
	request := requests.VolumeRemoveRequest{}
	c.DoRequestValid(&request, r)
}
