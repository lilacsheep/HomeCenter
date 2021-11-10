package api

import (
	"homeproxy/app/services/requests"

	"github.com/gogf/gf/net/ghttp"
)

type MinioApi struct {
	BaseControllers
}

func (a *MinioApi) QuerySettings(r *ghttp.Request) {
	request := requests.GetMinioSettingsRequest{}
	a.DoRequestValid(&request, r)
}

func (a *MinioApi) UpdateSettings(r *ghttp.Request) {
	request := requests.UpdateMinioSettingsRequest{}
	a.DoRequestValid(&request, r)
}

func (a *MinioApi) BucketsList(r *ghttp.Request) {
	request := requests.GetMinioBucketsRequest{}
	a.DoRequest(&request, r)
}

func (a *MinioApi) ObjectList(r *ghttp.Request) {
	request := requests.GetMinioObjectsRequest{}
	a.DoRequestValid(&request, r)
}