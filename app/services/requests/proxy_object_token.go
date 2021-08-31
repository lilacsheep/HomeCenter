package requests

import (
	"homeproxy/app/models"
	"net/http"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/guid"
)

type CreateObjectTokenRequest struct {
	Name      string `json:"name"`
	Effective int    `json:"effective"`
	Upload    bool   `json:"upload"`
	Download  bool   `json:"download"`
	Delete    bool   `json:"delete"`
	List      bool   `json:"list"`
}

func (self *CreateObjectTokenRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	query := g.DB().Model(&models.ObjectToken{})
	if c, err := query.Clone().Where("`name` = ?", self.Name).Count(); err != nil {
		return *response.SystemError(err)
	} else {
		if c == 0 {
			data := gconv.Map(self)
			data["secret_key"] = guid.S()
			if _, err := query.Data(data).Insert(); err != nil {
				return *response.SystemError(err)
			}
		} else {
			return *response.ErrorWithMessage(http.StatusInternalServerError, "token名已存在")
		}
	}
	return *response.Success()
}

type ListObjectTokenRequest struct {
	Pagination
}

func (self *ListObjectTokenRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	query := g.DB().Model(&models.ObjectToken{})
	var tokens []models.ObjectToken
	var c int
	c, _ = query.Count()
	err := query.Limit(self.OffsetLimit()...).Structs(&tokens)
	if err != nil {
		return *response.SystemError(err)
	}
	return *response.DataTable(tokens, c)
}

type DeleteObjectTokenRequest struct {
	Id int
}

func (self *DeleteObjectTokenRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	_, err := g.DB().Model(&models.ObjectToken{}).Where("`id` = ?", self.Id).Delete()
	if err != nil {
		return *response.SystemError(err)
	}
	return *response.Success()
}
