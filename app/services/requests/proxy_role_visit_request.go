package requests

import (
	"github.com/gogf/gf/net/ghttp"
	"homeproxy/app/models"
	"homeproxy/library/filedb"
	"homeproxy/library/mallory"
	"net/http"
)

type GetRoleAllVisitRequest struct{}

func (self *GetRoleAllVisitRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	data := models.AllVisitLogs()
	response.DataTable(data, len(data))
	return
}

func NewGetRoleAllVisitRequest() *GetRoleAllVisitRequest {
	return &GetRoleAllVisitRequest{}
}

type AddErrorSiteToProxyRequest struct {
	ID         string `json:"id"`
	InstanceID string `json:"instance_id"`
}

func (self *AddErrorSiteToProxyRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	data := mallory.ProxyRoleAnalysis{}
	err := filedb.DB.GetById(mallory.ProxyRoleAnalysisTable, self.ID, &data)
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		request := AddUrlRoleRequest{Url: data.Domain, InstanceID: self.InstanceID, Status: true}
		_, domain := request.UrlSplit(data.Domain)
		request.Url = domain
		response = request.Exec(r)
		if response.ErrorCode == 200 {
			filedb.DB.RemoveByID(mallory.ProxyRoleAnalysisTable, data.ID)
		}
	}
	return
}

func NewAddErrorSiteToProxyRequest() *AddErrorSiteToProxyRequest {
	return &AddErrorSiteToProxyRequest{}
}
