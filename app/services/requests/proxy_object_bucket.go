package requests

import (
	"homeproxy/app/models"
	"net/http"

	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

type CreateObjectBucketRequest struct {
	Name       string `json:"name"`
	Public     bool   `json:"public"`
	Referer    bool   `json:"referer"`
	RefererUrl string `json:"referer_url"`
}

func (self *CreateObjectBucketRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	data := gjson.New(self).Map()
	query := g.DB().Model(&models.Bucket{})
	if c, err := query.Clone().Where("`name` = ?", self.Name).Count(); err != nil {
		return *response.SystemError(err)
	} else {
		if c != 0 {
			return *response.ErrorWithMessage(http.StatusInternalServerError, "数据已经存在")
		}
	}
	_, err := query.Clone().Data(data).Insert()
	if err != nil {
		return *response.SystemError(err)
	}
	return *response.Success()
}

type ListBucketRequest struct {
	Pagination
}

func (self *ListBucketRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	query := g.DB().Model(&models.Bucket{})
	c, _ := query.Clone().Count()
	var data []models.Bucket
	err := query.Clone().Limit(self.OffsetLimit()...).Structs(&data)
	if err != nil {
		return *response.SystemError(err)
	}
	return *response.DataTable(data, c)
}

type DeleteBucketRequest struct {
	Id int `json:"id"`
}

func (self *DeleteBucketRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	query := g.DB().Model(&models.Bucket{})
	if _, err := query.Where("`id` = ?", self.Id).Delete(); err != nil {
		return *response.SystemError(err)
	}
	return *response.Success()
}

type UpdateBucketRequest struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Public     bool   `json:"public"`
	Referer    bool   `json:"referer"`
	RefererUrl string `json:"referer_url"`
}

func (self *UpdateBucketRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	data := gjson.New(self).Map()
	query := g.DB().Model(&models.Bucket{})
	if c, err := query.Clone().Where("`name` = ?", self.Name).Count(); err != nil {
		return *response.SystemError(err)
	} else {
		if c == 0 {
			return *response.ErrorWithMessage(http.StatusInternalServerError, "数据不存在")
		}
	}
	if _, err := query.Clone().Data(data).Where("`id` = ?", self.Id).Update(); err != nil {
		return *response.SystemError(err)
	}
	return *response.Success()
}
