package requests

import (
	"homeproxy/app/models"
	"net/http"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

type ListObjectRequest struct {
	Pagination
	Prefix string `json:"prefix"`
}

func (self *ListObjectRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	bucket := r.GetString("bucket")
	if bucket == "" {
		return *response.ErrorWithMessage(http.StatusBadRequest, "No bucket specified")
	}
	query := g.DB().Model(&models.ObjectInfo{})

	query = query.Where("`bucket` = ?", bucket)

	if self.Prefix != "" {
		query = query.Where("`name` like %%s%%", self.Prefix)
	}
	var objects []models.ObjectInfo

	if err := query.Structs(&objects); err != nil {
		return *response.SystemError(err)
	}
	return *response.SuccessWithDetail(objects)
}

type UploadObjectRequest struct{}

func (self *UploadObjectRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	return
}

type DeleteObjectRequest struct {
	Key  string `json:"key" v:"required"`
	Name string `json:"name" v:"required"`
}

func (self *DeleteObjectRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	bucket := r.GetString("bucket")
	if bucket == "" {
		return *response.ErrorWithMessage(http.StatusBadRequest, "No bucket specified")
	}
	query := g.DB().Model(&models.ObjectInfo{})

	if _, err := query.Where(g.Map{"bucket": bucket, "key": self.Key, "name": self.Name}).Delete(); err != nil {
		return *response.SystemError(err)
	}
	return *response.Success()
}
