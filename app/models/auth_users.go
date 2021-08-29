package models

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gmeta"
)

type User struct {
	DefaultModel
	Username   string `json:"username" storm:"unique"`
	Password   string `json:"password"`
	Status     bool   `json:"status"`
	gmeta.Meta `orm:"table:auth_users"`
}

func UserLogin(username, password string) (*User, error) {
	var (
		user = User{}
	)
	err := g.DB().Model(&User{}).Where("`username` = ? AND `password` = ?", username, password).Struct(&user)
	return &user, err
}

func CreateUser(username, password string) (user *User, err error) {
	user = &User{}
	user.Username = username
	user.Password = password
	user.Status = true
	_, err = g.DB().Model(&User{}).Save(user)
	return
}
