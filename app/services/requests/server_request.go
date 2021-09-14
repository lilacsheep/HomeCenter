package requests

import (
	"homeproxy/app/models"
	"net/http"

	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

type CreateServerRequest struct {
	Name       string `json:"name"`
	Address    string `json:"address" v:"address@required"`
	Port       int    `json:"port" v:"port@required"`
	Group      int    `json:"group"`
	Username   string `json:"username" v:"username@required"`
	Password   string `json:"password"`
	PrivateKey string `json:"private_key"`
	Remark     string `json:"remark"`
}

func (self *CreateServerRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	if self.Password == "" && self.PrivateKey == "" {
		return *response.ErrorWithMessage(http.StatusInternalServerError, "password or private_key need")
	}
	data := gjson.New(self)
	if self.Name == "" {
		data.Set("name", self.Address)
	}
	_, err := g.DB().Model(&models.Server{}).Data(data.Map()).Insert()
	if err != nil {
		return *response.SystemError(err)
	}
	return *response.Success()
}

type ServerListRequest struct {
	Pagination
	Group int `json:"group"`
}

func (self *ServerListRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	query := g.DB().Model(&models.Server{})
	if self.Group != 0 {
		query = query.Where("`group` = ?", self.Group)
	}
	var hosts []models.Server
	c, _ := query.Count()
	err := query.Limit(self.OffsetLimit()...).Structs(&hosts)
	if err != nil {
		return *response.SystemError(err)
	}
	return *response.DataTable(hosts, c)
}

type ServerUpdateRequest struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Address    string `json:"address"`
	Port       int    `json:"port"`
	Group      int    `json:"group"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	PrivateKey string `json:"private_key"`
	Remark     string `json:"remark"`
}

func (self *ServerUpdateRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	query := g.DB().Model(&models.Server{}).Where("`id` = ?", self.Id)
	var server models.Server
	err := query.Clone().Struct(&server)
	if err != nil {
		return *response.SystemError(err)
	}
	data := gjson.New(self)
	err = data.Remove("id")
	if err != nil {
		return *response.SystemError(err)
	}
	_, err = query.Clone().Data(data.Map()).Update()
	if err != nil {
		return *response.SystemError(err)
	}
	return *response.Success()
}

type ServerDeleteRequest struct {
	Id int `json:"id"`
}

func (self *ServerDeleteRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	query := g.DB().Model(&models.Server{}).Where("`id` = ?", self.Id)
	if c, _ := query.Clone().Count(); c != 0 {
		return *response.ErrorWithMessage(http.StatusInternalServerError, "数据不存在")
	}
	_, err := query.Delete()
	if err != nil {
		return *response.SystemError(err)
	}
	return *response.Success()
}

type CreateServerGroupRequest struct {
	Name   string `json:"name"`
	Remark string `json:"remark"`
}

func (self *CreateServerGroupRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	query := g.DB().Model(&models.ServerGroup{})
	c, _ := query.Clone().Where("`name` = ?", self.Name).Count()
	if c != 0 {
		return *response.ErrorWithMessage(http.StatusInternalServerError, "数据已存在")
	}
	data := gjson.New(self)
	_, err := query.Data(data.Map()).Insert()
	if err != nil {
		return *response.SystemError(err)
	}
	return *response.Success()
}

type RemoveServerGroupRequest struct {
	Id string `json:"id"`
}

func (self *RemoveServerGroupRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	query := g.DB().Model(&models.ServerGroup{}).Where("`id` = ?", self.Id)
	if c, _ := query.Clone().Count(); c == 0 {
		return *response.ErrorWithMessage(http.StatusInternalServerError, "数据不存在")
	}
	_, err := query.Delete()
	if err != nil {
		return *response.SystemError(err)
	}
	return *response.Success()
}

type ServerGroupListRequest struct {
	Pagination
}

func (self *ServerGroupListRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	query := g.DB().Model(&models.ServerGroup{})
	var groups []models.ServerGroup
	c, _ := query.Clone().Count()
	err := query.Clone().Limit(self.OffsetLimit()...).Structs(&groups)
	if err != nil {
		return *response.SystemError(err)
	}
	return *response.DataTable(groups, c)
}
