package requests

import (
	"homeproxy/app/models"
	"homeproxy/library/filedb2"
	"net/http"

	"github.com/gogf/gf/net/ghttp"
)

type QueryProxyMonitorInfoRequest struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

func (self *QueryProxyMonitorInfoRequest) Pagination() (int, int) {
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

func (self *QueryProxyMonitorInfoRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	var info []models.ProxyMonitorInfo
	_, limit := self.Pagination()
	err := filedb2.DB.Select().Limit(limit).Reverse().Find(&info)
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		response.DataTable(info, len(info))
	}
	return
}

func NewQueryProxyMonitorInfoRequest() *QueryProxyMonitorInfoRequest {
	return &QueryProxyMonitorInfoRequest{}
}
