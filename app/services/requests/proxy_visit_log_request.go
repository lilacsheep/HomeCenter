package requests

import (
	"homeproxy/library/filedb2"
	"homeproxy/library/mallory"
	"net/http"

	"github.com/gogf/gf/net/ghttp"
)

type QueryLogRequest struct {
	*Pagination
	Reverse bool `json:"reverse"`
}

func (self *QueryLogRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	var logs []mallory.ProxyVisitLog
	offset, limit := self.Next()
	query := filedb2.DB.Select().Limit(limit).Skip(offset)
	if self.Reverse {
		query = query.Reverse()
	}
	err := query.Find(&logs)
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		response.SuccessWithDetail(logs)
	}
	return
}

func NewQueryLogRequest() *QueryLogRequest {
	return &QueryLogRequest{}
}
