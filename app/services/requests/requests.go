package requests

import (
	"net/http"

	"github.com/gogf/gf/net/ghttp"
)

const (
	ParamsErrorCode = iota
	UserExistCode
	UsernameNotExistCode
	PasswordWrongCode
	UnknownCode int = 999
)

type RequestInterface interface {
	Exec(r *ghttp.Request) MessageResponse
}

type MessageResponse struct {
	ErrorCode int         `json:"code"`
	Message   string      `json:"message"`
	Detail    interface{} `json:"detail,omitempty"`
	Count     int         `json:"count,omitempty"`
	Next      int         `json:"next,omitempty"`
	Previous  int         `json:"previous,omitempty"`
}

func (resp *MessageResponse) DataTable(data interface{}, total int) *MessageResponse {
	resp.ErrorCode = http.StatusOK
	resp.Message = "成功"
	resp.Detail = data
	resp.Count = total
	return resp
}

func (resp *MessageResponse) Success() *MessageResponse {
	resp.ErrorCode = http.StatusOK
	resp.Message = "成功"
	return resp
}

func (resp *MessageResponse) SuccessWithDetail(data interface{}) *MessageResponse {
	resp.ErrorCode = http.StatusOK
	resp.Message = "成功"
	resp.Detail = data
	return resp
}

func (resp *MessageResponse) ErrorWithMessage(code int, message interface{}) *MessageResponse {
	resp.ErrorCode = code
	resp.Message = "失败"
	resp.Detail = message
	return resp
}

func (resp *MessageResponse) SystemError(err error) *MessageResponse {
	resp.ErrorCode = http.StatusInternalServerError
	resp.Message = "失败"
	resp.Detail = err.Error()
	return resp
}

type Pagination struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

func (self *Pagination) OffsetLimit() []int {
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
	return []int{offset, limit}
}

func (self *Pagination) Next() (int, int) {
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
