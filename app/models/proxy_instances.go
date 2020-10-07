package models

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
)

const ProxyInstanceTable = "proxy_instances"

type ProxyInstance struct {
	ID         string `orm:"id"`
	Protocol   int    `orm:"protocol"` // default 0(SSH)
	Address    string `orm:"address"`  //
	Username   string `orm:"username"`
	Password   string `orm:"password"`
	PrivateKey string `orm:"private_key"`
	Status     bool   `orm:"status"`
}

func (self *ProxyInstance) Url() string {
	// ssh://user:passwd@192.168.1.1:1122
	if self.Password == "" {
		return fmt.Sprintf("ssh://%s@%s", self.Username, self.Address)
	} else {
		return fmt.Sprintf("ssh://%s:%s@%s", self.Username, self.Password, self.Address)
	}
}

func GetEnableProxyInstances() (instances []ProxyInstance, err error) {
	err = g.DB().Table(ProxyInstanceTable).Structs(&instances, "status = true")
	return
}

func GetAllProxyInstances() (instances []ProxyInstance, err error) {
	err = g.DB().Table(ProxyInstanceTable).Structs(&instances)
	return
}
