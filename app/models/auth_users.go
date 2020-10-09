package models

import (
	"time"
)

const AuthUserTable = "auth_users"

type AuthUser struct {
	ID       int
	Username string
	Password string
	Status   bool
	CreateAt time.Time
}

//func UserLogin(username, password string) (user *AuthUser, err error) {
//	err = g.DB().Table(AuthUserTable).Where("username = ? and password = ?", username, password).Struct(user)
//	return
//}
//
//func CreateUser(username, password string) (user *AuthUser, err error) {
//	user = &AuthUser{Username: username, Password: password}
//	_, err = g.DB().Table(AuthUserTable).Insert(&user)
//	return
//}
