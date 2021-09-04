package models

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gmeta"
)

type ProxyRole struct {
	DefaultModel
	InstanceID int    `json:"instance_id"`
	Status     bool   `json:"status"`
	Sub        string `json:"sub"`
	Domain     string `json:"domain"`
	gmeta.Meta `orm:"table:proxy_role"`
}

func AllProxyRole() ([]ProxyRole, error) {
	var all []ProxyRole
	err := g.DB().Model(&ProxyRole{}).Structs(&all)
	return all, err
}
