package requests

import (
	"homeproxy/app/models"

	"github.com/gogf/gf/net/ghttp"
)

//
//type GetSelfRequest struct{}
//
//func (self *GetSelfRequest) Exec(r *ghttp.Request) (response MessageResponse) {
//	if user := r.Session.Get("user"); user == nil {
//		response.ErrorWithMessage(http.StatusUnauthorized, "need login")
//	} else {
//		response.SuccessWithDetail(user)
//	}
//	return
//}
//
//func NewGetSelfRequest() *GetSelfRequest {
//	return &GetSelfRequest{}
//}
//
type LoginRequest struct {
	Username string `v:"username@required|length:6,30#请输入用户名称|用户名称长度非法"`
	Password string `v:"password@required|password2"`
}

func (self LoginRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	user, err := models.UserLogin(self.Username, self.Password)
	if err != nil {
		response.ErrorWithMessage(500, err.Error())
	} else {
		user.Password = ""
		response.SuccessWithDetail(user)
	}
	return
}

func NewLoginRequest() *LoginRequest {
	return &LoginRequest{}
}
//
//type CreateUserRequest struct {
//	Username string `v:"username@required|length:6,30#请输入用户名称|用户名称长度非法"`
//	Password string `v:"password@required|password2"`
//}
//
//func (self CreateUserRequest) Exec(r *ghttp.Request) (response MessageResponse) {
//	user, err := models.CreateUser(self.Username, self.Password)
//	if err != nil {
//		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
//	} else {
//		user.Password = ""
//		response.SuccessWithDetail(user)
//	}
//	return
//}
//
//func NewCreateUserRequest() *CreateUserRequest {
//	return &CreateUserRequest{}
//}
