package requests

import (
	"context"
	"homeproxy/app/models"

	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/net/ghttp"
	"github.com/minio/minio-go/v7"
)

type GetMinioSettingsRequest struct{}

func (self *GetMinioSettingsRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	settings, err := models.GetMinioConfig()
	if err != nil {
		return *response.SystemError(err)
	} else {
		return *response.SuccessWithDetail(settings)
	}
}

type UpdateMinioSettingsRequest struct {
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
}

func (self *UpdateMinioSettingsRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	oldOptions, err := models.GetMinioConfig()
	if err != nil {
		return *response.SystemError(err)
	}
	newOptions := gvar.New(self)

	for k, v := range gvar.New(self).MapStrStr() {
		models.UpdateConfig("minio", k, v)
	}
	if self.AutoStart && !oldOptions.AutoStart {
		oldOptions.StopContainer()
		oldOptions.RemoveContainer()
		err = oldOptions.Auto()
		if err != nil {
			return *response.SystemError(err)
		}
	}
	if self.AutoStart && oldOptions.AutoStart {
		if oldOptions.NeedReCreate(newOptions) {
			oldOptions.StopContainer()
			oldOptions.RemoveContainer()
			err = oldOptions.Auto()
			if err != nil {
				return *response.SystemError(err)
			}
		}
	}
	if !self.AutoStart && oldOptions.AutoStart {
		oldOptions.StopContainer()
		oldOptions.RemoveContainer()
	}
	return *response.Success()
}

type GetMinioBucketsRequest struct{}

func (self *GetMinioBucketsRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	cli, err := models.GetMinioClient()
	if err != nil {
		return *response.SystemError(err)
	}
	buckets, err := cli.ListBuckets(context.Background())
	if err != nil {
		return *response.SystemError(err)
	}
	return *response.SuccessWithDetail(buckets)
}

type GetMinioObjectsRequest struct {
	BucketName string `json:"bucket_name" v:"required"`
	Prefix     string `json:"prefix"`
}

func (self *GetMinioObjectsRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	cli, err := models.GetMinioClient()
	if err != nil {
		return *response.SystemError(err)
	}
	objectCh := cli.ListObjects(context.Background(), self.BucketName, minio.ListObjectsOptions{
		Prefix:    self.Prefix,
		Recursive: false,
		UseV1:     false,
	})
	var objs []minio.ObjectInfo
	for object := range objectCh {
		if object.Err != nil {
			return *response.SystemError(object.Err)
		}
		objs = append(objs, object)
	}

	return *response.SuccessWithDetail(objs)
}
