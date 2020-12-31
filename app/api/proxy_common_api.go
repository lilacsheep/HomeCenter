package api

import (
	"homeproxy/app/services/requests"

	"github.com/gogf/gf/net/ghttp")


type CommonApi struct {
	BaseControllers
}

func (self *CommonApi) Countrys(r *ghttp.Request) {
	request := requests.NewGetAllCountryRequest()
	self.DoRequest(request, r)
}