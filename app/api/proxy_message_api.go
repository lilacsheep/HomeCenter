package api

import (
	"github.com/gogf/gf/net/ghttp"
	"homeproxy/app/services/requests"
)

type ProxyMessageApi struct {
	BaseControllers
}

func (self *ProxyMessageApi) All(r *ghttp.Request) {
	request := requests.NewGetAllMessagesRequest()
	self.DoRequest(request, r)
}
