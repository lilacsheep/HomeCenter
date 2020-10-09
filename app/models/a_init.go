package models

import (
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/glog"
	"homeproxy/library/filedb"
	"time"
)

var DB *filedb.Database

var (
	Dbname string
	Dbpath string
)

func InitDB() {
	if !gfile.Exists(Dbpath) {
		_ = gfile.Mkdir(Dbpath)
	}
	DB = filedb.NewDatabase(Dbname, Dbpath)
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
