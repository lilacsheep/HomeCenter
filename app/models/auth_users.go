package models

import (
	"homeproxy/library/filedb2"
	"time"

	"github.com/asdine/storm/v3/q"
)

type User struct {
	ID       int       `json:"id" storm:"id,increment"`
	Username string    `json:"username" storm:"unique"`
	Password string    `json:"password"`
	Status   bool      `json:"status"`
	CreateAt time.Time `json:"create_at"`
}

func UserLogin(username, password string) (*User, error) {
	var (
		user = User{}
	)
	query := filedb2.DB.Select(q.And(q.Eq("Username", username), q.Eq("Password", password)))
	err := query.First(&user)
	return &user, err
}

func CreateUser(username, password string) (user *User, err error) {
	user = &User{}
	user.Username = username
	user.Password = password
	user.Status = true
	user.CreateAt = time.Now()
	err = filedb2.DB.Save(user)
	return
}
