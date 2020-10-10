package models

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"homeproxy/library/filedb"
)

const ProxyInstanceTable = "proxy_instances"

type ProxyInstance struct {
	ID         string `json:"id"`
	Protocol   int    `json:"protocol"` // default 0(SSH)
	Address    string `json:"address"`  //
	Username   string `json:"username"`
	Password   string `json:"password"`
	PrivateKey string `json:"private_key"`
	Status     bool   `json:"status"`
	Delay      int    `json:"delay"`
}

func (self *ProxyInstance) Url() string {
	return fmt.Sprintf("ssh://%s@%s", self.Username, self.Address)
}

func GetEnableProxyInstances() (instances []ProxyInstance, err error) {
	var c *filedb.Collection
	if c, err = DB.Collection(ProxyInstanceTable); err != nil {
		return
	} else {
		if err = c.Search(g.Map{"status": true}, &instances); err != nil {
			return
		}
	}
	return
}
