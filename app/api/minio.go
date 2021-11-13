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

func (a *MinioApi) BucketCreate(r *ghttp.Request) {
	request := requests.CreateMinioBucketRequest{}
	a.DoRequestValid(&request, r)
}

func (a *MinioApi) ObjectList(r *ghttp.Request) {
	request := requests.GetMinioObjectsRequest{}
	a.DoRequestValid(&request, r)
}

func (a *MinioApi) ObjectUpload(r *ghttp.Request) {
	request := requests.UploadObjectRequest{}
	a.DoRequestValid(&request, r)
}

func (a *MinioApi) ObjectInfo(r *ghttp.Request) {
	request := requests.GetObjectInfoRequest{}
	a.DoRequestValid(&request, r)
}

func (a *MinioApi) ObjectDownload(r *ghttp.Request) {
	request := requests.DownloadObjectRequest{}
	a.DoRequestValid(&request, r)
}

func (a *MinioApi) ObjectShare(r *ghttp.Request) {
	request := requests.ShareObjectRequest{}
	a.DoRequestValid(&request, r)
}