package api

import (
	"homeproxy/app/services/requests"

	"github.com/gogf/gf/net/ghttp"
)

type ImageController struct {
	BaseControllers
}

func (c *ImageController) List(r *ghttp.Request) {
	request := requests.ListImagesRequest{}
	c.DoRequestValid(&request, r)
}

func (c *ImageController) Pull(r *ghttp.Request) {
	request := requests.PullImageRequest{}
	c.DoRequestValid(&request, r)
}

func (c *ImageController) Remove(r *ghttp.Request) {
	request := requests.DeleteImageRequest{}
	c.DoRequestValid(&request, r)
}

func (c *ImageController) Info(r *ghttp.Request) {
	request := requests.ImageInspectRequest{}
	c.DoRequestValid(&request, r)
}

func (c *ImageController) History(r *ghttp.Request) {
	request := requests.ImageHistoryRequest{}
	c.DoRequestValid(&request, r)
}
