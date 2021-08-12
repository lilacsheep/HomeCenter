package api

import (
	"homeproxy/app/services/requests"

	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
)

type BaseControllers struct{}

// DoRequest 处理请求
func (c *BaseControllers) DoRequestValid(request requests.RequestInterface, r *ghttp.Request) {
	if err := r.Parse(request); err != nil {
		resp := requests.MessageResponse{}
		switch err.(type) {
			case gvalid.Error:
				resp.ErrorWithMessage(requests.ParamsErrorCode, err.(gvalid.Error).Maps())
				r.Response.WriteJsonExit(resp)
			default:
				resp.ErrorWithMessage(requests.ParamsErrorCode, err.Error())
				r.Response.WriteJsonExit(resp)
		}
		
	}
	resp := request.Exec(r)
	r.Response.WriteJsonExit(resp)
}

func (c *BaseControllers) DoRequest(request requests.RequestInterface, r *ghttp.Request) {
	resp := request.Exec(r)
	r.Response.WriteJsonExit(resp)
}
