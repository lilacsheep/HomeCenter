package models

import (
	"fmt"
	"homeproxy/library/filedb2"

	"github.com/asdine/storm/v3/q"
)

type ProxyInstance struct {
	ID         int    `json:"id" storm:"id,increment"`
	Protocol   int    `json:"protocol"`               // default 0(SSH)
	Address    string `json:"address" storm:"unique"` //
	Username   string `json:"username"`
	Password   string `json:"password"`
	PrivateKey string `json:"private_key"`
	Status     bool   `json:"status"`
	Delay      int    `json:"delay"`
}

func (self *ProxyInstance) Url() string {
	return fmt.Sprintf("ssh://%s@%s", self.Username, self.Address)
}

func UpdateProxyInstanceDelay(id int, delay int) error {
	return filedb2.DB.UpdateField(&ProxyInstance{ID: id}, "Delay", delay)
}

func GetEnableProxyInstances() (instances []ProxyInstance, err error) {
	err = filedb2.DB.Select(q.Eq("Status", true)).Find(&instances)
	return
}
