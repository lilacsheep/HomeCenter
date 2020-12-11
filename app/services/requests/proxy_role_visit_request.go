package requests

import (
	"github.com/gogf/gf/net/ghttp"
	"homeproxy/app/models"
	"homeproxy/library/filedb2"
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
