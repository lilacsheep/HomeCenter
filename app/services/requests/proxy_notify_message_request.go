package requests

import (
	"github.com/gogf/gf/net/ghttp"
	"homeproxy/app/models"
	"homeproxy/library/filedb"
	"net/http"
)

type GetAllMessagesRequest struct{}

func (self *GetAllMessagesRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	var (
		data []models.NotifyMessage
	)
	err := filedb.DB.QueryAll(models.ProxyNotifyMessageTable, &data)
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		response.DataTable(data, len(data))
	}
	return
}

func NewGetAllMessagesRequest() *GetAllMessagesRequest {
	return &GetAllMessagesRequest{}
}
