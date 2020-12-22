package api

import (
	"homeproxy/app/services/requests"

	"github.com/gogf/gf/net/ghttp"
)

type AuthController struct {
	BaseControllers
}

//func (self *AuthController) MySelf(r *ghttp.Request) {
//	request := requests.NewGetSelfRequest()
//	self.DoRequest(request, r)
//}
//
// func (self *AuthController) CreateUser(r *ghttp.Request) {
// 	request := requests.NewCreateUserRequest()
// 	self.DoRequestValid(request, r)
// }

func (self *AuthController) LoginUser(r *ghttp.Request) {
	request := requests.NewLoginRequest()
	self.DoRequestValid(request, r)
}
