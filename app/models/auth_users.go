package models

import (
	"homeproxy/library/filedb2"
	"time"

	"github.com/asdine/storm/v3/q"
)

func init() {
	defaultUser := User{Username: "admin", Password: "!QAZ2wsx", Status: true, CreateAt: time.Now()}
	if c, _ := filedb2.DB.Count(&User{});c == 0 {
		filedb2.DB.Save(&defaultUser)
	}
}

type User struct {
	ID       int       `json:"id" storm:"id,increment"`
	Username string    `json:"username" storm:"unique"`
	Password string    `json:"password"`
	Status   bool      `json:"status"`
	CreateAt time.Time `json:"create_at"`
}

func UserLogin(username, password string) (user *User, err error) {
	err = filedb2.DB.Select(q.And(q.Eq("Username", username), q.Eq("Password", password))).Find(&user)
	return
}

func CreateUser(username, password string) (user *User, err error) {
	user.Username = username
	user.Password = password
	user.Status = true
	user.CreateAt = time.Now()
	err = filedb2.DB.Save(user)
	return
}
