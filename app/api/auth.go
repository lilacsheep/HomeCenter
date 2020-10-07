package api

import (
	"github.com/gogf/gf/net/ghttp"
	"homeproxy/app/services/requests"
)

type AuthController struct {
	BaseControllers
}

func (self *AuthController) MySelf(r *ghttp.Request) {
	request := requests.NewGetSelfRequest()
	self.DoRequest(request, r)
}

func (self *AuthController) CreateUser(r *ghttp.Request) {
	request := requests.NewCreateUserRequest()
	self.DoRequestValid(request, r)
}

func (self *AuthController) LoginUser(r *ghttp.Request) {
	request := requests.NewLoginRequest()
	self.DoRequestValid(request, r)
}
