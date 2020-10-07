package models

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"homeproxy/library/mallory"
)

const ProxyRoleTable = "proxy_role_list"

func AllRoles() (proxies []mallory.ProxyRole) {
	if err := g.DB().Table(ProxyRoleTable).Structs(&proxies); err != nil {
		glog.Errorf("get all proxies error: %s", err.Error())
	}
	return proxies
}
