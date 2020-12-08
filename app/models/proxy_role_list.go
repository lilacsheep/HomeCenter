package models

import (
	"github.com/gogf/gf/os/glog"
	"homeproxy/library/filedb"
	"homeproxy/library/mallory"
)

func init() {
	if err := filedb.DB.NewCollections(mallory.ProxyRoleTable, nil); err != nil {
		if err != filedb.ErrCollectionExist {
			glog.Error("init collection error: %s", err.Error())
		}
	}

	setting := filedb.DefaultCollectionSettings()
	setting.AutoDump = false
	setting.MaxRecord = 1000
	if err := filedb.DB.NewCollections(mallory.ProxyRoleAnalysisTable, setting); err != nil {
		if err != filedb.ErrCollectionExist {
			glog.Error("init collection error: %s", err.Error())
		}
	} else {
		// 每次启动清空访问失败的网站
		_ = filedb.DB.Truncate(mallory.ProxyRoleAnalysisTable)
	}

	setting = filedb.DefaultCollectionSettings()
	setting.AutoDump = false
	setting.MaxRecord = 1000
	if err := filedb.DB.NewCollections(mallory.ProxyVisitLogTable, setting); err != nil {
		if err != filedb.ErrCollectionExist {
			glog.Error("init collection error: %s", err.Error())
		}
	}
}

func AllRoles() (proxies []mallory.ProxyRole) {
	c, _ := filedb.DB.Collection(mallory.ProxyRoleTable)
	if err := c.All(&proxies); err != nil {
		glog.Errorf("get all proxies error: %s", err.Error())
	}
	return proxies
}

func AllVisitLogs() (data []mallory.ProxyRoleAnalysis) {
	err := filedb.DB.QueryAll(mallory.ProxyRoleAnalysisTable, &data)
	if err != nil {
		glog.Errorf("get all proxies error: %s", err.Error())
	}
	return
}
