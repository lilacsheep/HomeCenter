package requests

import (
	"homeproxy/library/filedb2"
	"homeproxy/library/mallory"
	"net/http"

	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
	"github.com/gogf/gf/net/ghttp"
)

type GetRoleAllVisitRequest struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Filter string `json:"filter"`
}

func (self *GetRoleAllVisitRequest) pagination() (int, int) {
	var (
		limit  = 10
		page   = 1
		offset = 0
	)
	if self.Limit != 0 {
		limit = self.Limit
	}
	if self.Page != 0 {
		page = self.Page
	}
	offset = (page - 1) * limit
	return offset, limit
}

func (self *GetRoleAllVisitRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	var data []mallory.ProxyRoleAnalysis
	offset, limit := self.pagination()
	var query storm.Query
	if self.Filter != "" {
		query = filedb2.DB.Select(q.Re("Domain", self.Filter))
	} else {
		query = filedb2.DB.Select()
	}
	c, _ := query.Count(&mallory.ProxyRoleAnalysis{})

	err := query.Skip(offset).Limit(limit).Find(&data)
	if err != nil {
		if err == storm.ErrNotFound {
			response.DataTable([]mallory.ProxyVisitLog{}, 0)
		} else {
			response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
		}
	} else {
		response.DataTable(data, c)
	}
	return
}

func NewGetRoleAllVisitRequest() *GetRoleAllVisitRequest {
	return &GetRoleAllVisitRequest{}
}

type AddErrorSiteToProxyRequest struct {
	ID         int `json:"id"`
	InstanceID int `json:"instance_id"`
}

func (self *AddErrorSiteToProxyRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	data := mallory.ProxyRoleAnalysis{}
	err := filedb2.DB.One("ID", self.ID, &data)

	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		request := AddUrlRoleRequest{Url: data.Domain, InstanceID: self.InstanceID, Status: true}
		_, domain := request.UrlSplit(data.Domain)
		request.Url = domain
		response = request.Exec(r)
		if response.ErrorCode == 200 {
			filedb2.DB.DeleteStruct(&data)
		}
	}
	return
}

func NewAddErrorSiteToProxyRequest() *AddErrorSiteToProxyRequest {
	return &AddErrorSiteToProxyRequest{}
}
