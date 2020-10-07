package requests

import (
	"github.com/gogf/gf/net/ghttp"
	"net/http"
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

func (resp *MessageResponse) DataTable(data interface{}, total int) {
	resp.ErrorCode = http.StatusOK
	resp.Message = "成功"
	resp.Detail = data
	resp.Count = total
}

func (resp *MessageResponse) Success() {
	resp.ErrorCode = http.StatusOK
	resp.Message = "成功"
}

func (resp *MessageResponse) SuccessWithDetail(data interface{}) {
	resp.ErrorCode = http.StatusOK
	resp.Message = "成功"
	resp.Detail = data
}

func (resp *MessageResponse) ErrorWithMessage(code int, message interface{}) {
	resp.ErrorCode = code
	resp.Message = "失败"
	resp.Detail = message
}
