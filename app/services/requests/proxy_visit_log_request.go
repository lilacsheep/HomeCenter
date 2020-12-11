package requests

import (
	"homeproxy/library/filedb2"
	"homeproxy/library/mallory"
	"net/http"

	"github.com/gogf/gf/net/ghttp"
)

type QueryLogRequest struct {
	Limit   int  `json:"limit"`
	Page    int  `json:"page"`
	Reverse bool `json:"reverse"`
}

func (self *QueryLogRequest) Pagination() (int, int) {
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

func (self *QueryLogRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	var logs []mallory.ProxyVisitLog
	offset, limit := self.Pagination()
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
