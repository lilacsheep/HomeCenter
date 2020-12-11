package models

import (
	"github.com/gogf/gf/os/glog"
	"homeproxy/library/filedb2"
	"homeproxy/library/mallory"
)

func AllRoles() (proxies []mallory.ProxyRole) {
	if err := filedb2.DB.All(&proxies); err != nil {
		glog.Errorf("get all proxies error: %s", err.Error())
	}
	return proxies
}

func AllVisitLogs() (data []mallory.ProxyRoleAnalysis) {
	err := filedb2.DB.All(&data)
	if err != nil {
		glog.Errorf("get all proxies error: %s", err.Error())
	}
	return
}
