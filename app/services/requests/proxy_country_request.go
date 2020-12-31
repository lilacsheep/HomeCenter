package requests

import (
	"homeproxy/library/common"

	"github.com/gogf/gf/net/ghttp"
)


type GetAllCountryRequest struct{}

func (self *GetAllCountryRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	data := common.Countrys.Map()
	response.SuccessWithDetail(data)
	return
}

func NewGetAllCountryRequest() *GetAllCountryRequest {
	return &GetAllCountryRequest{}
}