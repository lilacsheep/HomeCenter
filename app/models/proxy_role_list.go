package models

import (
	"github.com/gogf/gf/os/glog"
	"homeproxy/library/mallory"
)

const ProxyRoleTable = "proxy_role_list"

func AllRoles() (proxies []mallory.ProxyRole) {
	c, _ := DB.Collection(ProxyRoleTable)
	if err := c.All(&proxies); err != nil {
		glog.Errorf("get all proxies error: %s", err.Error())
	}
	return proxies
}
