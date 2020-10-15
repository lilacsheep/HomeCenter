package requests

import (
	"github.com/gogf/gf/net/ghttp"
	"homeproxy/app/models"
	"net/http"
)

type QueryProxyMonitorInfoRequest struct{}

func (self *QueryProxyMonitorInfoRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	info, err := models.GetAllProxyMonitorInfo()
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err)
	} else {
		response.DataTable(info, len(info))
	}
	return
}

func NewQueryProxyMonitorInfoRequest() *QueryProxyMonitorInfoRequest {
	return &QueryProxyMonitorInfoRequest{}
}
