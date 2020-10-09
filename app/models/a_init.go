package models

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"homeproxy/library/filedb"
	"time"
)

var DB *filedb.Database

func init() {
	path := g.Cfg().GetString("db.path", "db")
	name := g.Cfg().GetString("db.name", "default")
	DB = filedb.NewDatabase(name, path)
	if err := DB.NewCollections(ProxyInstanceTable); err != nil {
		if err != filedb.ErrCollectionExist {
			glog.Error("init collection error: %s", err.Error())
		}
	}
	if err := DB.NewCollections(ProxyRoleTable); err != nil {
		if err != filedb.ErrCollectionExist {
			glog.Error("init collection error: %s", err.Error())
		}
	}
	if err := DB.NewCollections(ProxyServerTable); err != nil {
		if err != filedb.ErrCollectionExist {
			glog.Error("init collection error: %s", err.Error())
		}
	} else {
		server := ProxyServer{
			Name:      "default",
			Port:      1316,
			Status:    true,
			AutoProxy: false,
		}
		c, _ := DB.Collection(ProxyServerTable)
		c.Insert(&server)
	}
	go func() {
		for {
			DB.Dump()
			time.Sleep(2 * time.Second)
		}
	}()
}
