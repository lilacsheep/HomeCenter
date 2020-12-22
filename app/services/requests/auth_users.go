package requests

import (
	"homeproxy/app/models"
	"homeproxy/library/filedb2"
	"net/http"

	"github.com/gogf/gf/net/ghttp"
)

type GetSelfRequest struct{}

func (self *GetSelfRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	if user := r.Session.Get("user"); user == nil {
		response.ErrorWithMessage(http.StatusUnauthorized, "need login")
	} else {
		response.SuccessWithDetail(user)
	}
	return
}

func NewGetSelfRequest() *GetSelfRequest {
	return &GetSelfRequest{}
}

type LoginRequest struct {
	Username string `v:"username@required|length:5,30#请输入用户名称|用户名称长度非法"`
	Password string `v:"password@required|password2"`
}

func (self LoginRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	user, err := models.UserLogin(self.Username, self.Password)
	if err != nil {
		response.ErrorWithMessage(500, err.Error())
	} else {
		user.Password = ""
		r.Session.Set("user", user)
		response.SuccessWithDetail(user)
	}
	return
}

func NewLoginRequest() *LoginRequest {
	return &LoginRequest{}
}

type ChangeSelfPasswordRequest struct {
	Current   string `json:"current" v:"current@required"`
	Password1 string `json:"password1" v:"password1@required|password2"`
	Password2 string `json:"password2" v:"password2@required|same:password1"`
}

func (self *ChangeSelfPasswordRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	user := r.Session.Get("user")
	if user == nil {
		response.ErrorWithMessage(http.StatusUnauthorized, "need login")
	} else {
		u, err := models.UserLogin(user.(*models.User).Username, self.Current)
		if err != nil {
			response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
		} else {
			err = filedb2.DB.UpdateField(u, "Password", self.Password1)
			if err != nil {
				response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
			} else {
				response.Success()
			}
		}
	}
	return
}

func NewChangeSelfPasswordRequest() *ChangeSelfPasswordRequest {
	return &ChangeSelfPasswordRequest{}
}

type LogoutRequest struct {}

func (self *LogoutRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	r.Session.RemoveAll()
	response.Success()
	return
}

func NewLogoutRequest() *LogoutRequest {
	return &LogoutRequest{}
}


type CreateUserRequest struct {
	Username string `v:"username@required|length:5,30#请输入用户名称|用户名称长度非法"`
	Password1 string `v:"password1@required|password2"`
	Password2 string `json:"password2" v:"password2@required|same:password1"`
}

func (self CreateUserRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	user, err := models.CreateUser(self.Username, self.Password1)
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		user.Password = ""
		response.SuccessWithDetail(user)
	}
	return
}

func NewCreateUserRequest() *CreateUserRequest {
	return &CreateUserRequest{}
}
