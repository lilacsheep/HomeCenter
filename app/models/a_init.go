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
	settings := filedb.DefaultCollectionSettings()
	settings.Unique = "address"
	if err := DB.NewCollections(ProxyInstanceTable, settings); err != nil {
		if err != filedb.ErrCollectionExist {
			glog.Error("init collection error: %s", err.Error())
		}
	}
	if err := DB.NewCollections(ProxyRoleTable, nil); err != nil {
		if err != filedb.ErrCollectionExist {
			glog.Error("init collection error: %s", err.Error())
		}
	}
	settings = filedb.DefaultCollectionSettings()
	settings.Unique = "url"
	if err := DB.NewCollections(DownloadListTable, settings); err != nil {
		if err != filedb.ErrCollectionExist {
			glog.Error("init collection error: %s", err.Error())
		}
	}
	settings = filedb.DefaultCollectionSettings()
	settings.AutoDump = false
	settings.MaxRecord = 10
	if err := DB.NewCollections(ProxyMonitorTable, settings); err != nil {
		if err != filedb.ErrCollectionExist {
			glog.Error("init collection error: %s", err.Error())
		}
	}
	if err := DB.NewCollections(ProxyServerTable, nil); err != nil {
		if err != filedb.ErrCollectionExist {
			glog.Error("init collection error: %s", err.Error())
		}
	} else {
		server := ProxyServer{
			Name:      "default",
			Port:      1316,
			Status:    true,
			AutoProxy: false,
			AutoStart: true,
		}
		c, _ := DB.Collection(ProxyServerTable)
		_, err := c.Insert(server)
		if err != nil {
			glog.Errorf("init server info error: %s", err)
		}
	}
	go func() {
		for {
			DB.Dump()
			time.Sleep(2 * time.Second)
		}
	}()
}
