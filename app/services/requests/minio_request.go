package requests

import (
	"context"
	"fmt"
	"homeproxy/app/models"
	"io"
	"net/http"
	"time"

	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/util/gconv"
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
	Public         bool   `json:"public"`
	AutoStart      bool   `json:"auto_start"`
	Port           int    `json:"port"`
	AccessKey      string `json:"access_key"`
	SecretKey      string `json:"secret_key"`
	WebUi          bool   `json:"webui"`
	WebUiPort      int    `json:"webui_port"`
	SavePath       string `json:"save_path"`
	ConfigDir      string `json:"config_dir"`
	MinioDomain    string `json:"minio_domain"`
	MinioServerUrl string `json:"server_url"`
	Region         string `json:"region"`
	RegionComment  string `json:"region_comment"`
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

type CreateMinioBucketRequest struct {
	Name string `json:"name"`
}

func (self *CreateMinioBucketRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	cli, err := models.GetMinioClient()
	if err != nil {
		return *response.SystemError(err)
	}
	err = cli.MakeBucket(context.Background(), self.Name, minio.MakeBucketOptions{})
	if err != nil {
		return *response.SystemError(err)
	}
	return *response.Success()
}

type UploadObjectRequest struct {
	Name   string `json:"name"`
	Prefix string `json:"prefix"`
}

func (self *UploadObjectRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	file := r.GetUploadFile("file")
	if file == nil {
		return *response.ErrorWithMessage(http.StatusBadRequest, "file is none")
	}
	cli, err := models.GetMinioClient()
	if err != nil {
		return *response.SystemError(err)
	}
	fd, err := file.Open()
	if err != nil {
		return *response.SystemError(err)
	}
	defer fd.Close()
	name := gfile.Basename(file.Filename)
	info, err := cli.PutObject(context.Background(), self.Name, name, fd, file.Size, minio.PutObjectOptions{})
	if err != nil {
		return *response.SystemError(err)
	}
	return *response.SuccessWithDetail(info)
}

type DownloadObjectRequest struct {
	BucketName string `json:"bucket_name"`
	ObjectName string `json:"object_name"`
}

func (self *DownloadObjectRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	cli, err := models.GetMinioClient()
	if err != nil {
		return *response.SystemError(err)
	}

	obj, err := cli.GetObject(context.Background(), self.BucketName, self.ObjectName, minio.GetObjectOptions{})
	if err != nil {
		return *response.SystemError(err)
	}
	defer obj.Close()
	info, err := obj.Stat()
	if err != nil {
		return *response.SystemError(err)
	}
	name := gfile.Basename(info.Key)
	r.Response.Header().Set("Content-Type", "application/force-download")
	r.Response.Header().Set("Accept-Ranges", "bytes")
	r.Response.Header().Set("Content-Length", gconv.String(info.Size))
	r.Response.Header().Set("Content-Disposition", fmt.Sprintf(`attachment;filename="%s"`, name))

	_, err = io.Copy(r.Response.Writer, obj)
	if err != nil {
		return *response.SystemError(err)
	}
	r.Exit()
	return *response.Success()
}

type GetObjectInfoRequest struct {
	BucketName string `json:"bucket_name"`
	ObjectName string `json:"object_name"`
}

func (self *GetObjectInfoRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	cli, err := models.GetMinioClient()
	if err != nil {
		return *response.SystemError(err)
	}
	info, err := cli.StatObject(context.Background(), self.BucketName, self.ObjectName, minio.GetObjectOptions{})
	if err != nil {
		return *response.SystemError(err)
	}
	return *response.SuccessWithDetail(info)
}

type ShareObjectRequest struct {
	BucketName string `json:"bucket_name"`
	ObjectName string `json:"object_name"`
}

func (self *ShareObjectRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	var (
		cli      *minio.Client
		err      error
	)
	Options, err := models.GetMinioConfig()
	if err != nil {
		return *response.SystemError(err)
	}
	if Options.MinioServerUrl == "" {
		cli, err = models.GetMinioClient()
		if err != nil {
			return *response.SystemError(err)
		}
	} else {
		cli, err = models.GetPublicClient()
		if err != nil {
			return *response.SystemError(err)
		}
	}

	query, err := cli.PresignedGetObject(context.Background(), self.BucketName, self.ObjectName, time.Hour*168-1, nil)
	if err != nil {
		return *response.SystemError(err)
	}

	return *response.SuccessWithDetail(query.String())
}
