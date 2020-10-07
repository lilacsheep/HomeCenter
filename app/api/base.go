package api

import (
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
	"homeproxy/app/services/requests"
)

type BaseControllers struct{}

// DoRequest 处理请求
func (c *BaseControllers) DoRequestValid(request requests.RequestInterface, r *ghttp.Request) {
	if err := r.Parse(request); err != nil {
		resp := requests.MessageResponse{}
		resp.ErrorWithMessage(requests.ParamsErrorCode, err.(*gvalid.Error).Maps())
		r.Response.WriteJsonExit(resp)
	}
	resp := request.Exec(r)
	r.Response.WriteJsonExit(resp)
}

func (c *BaseControllers) DoRequest(request requests.RequestInterface, r *ghttp.Request) {
	resp := request.Exec(r)
	r.Response.WriteJsonExit(resp)
}
