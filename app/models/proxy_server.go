package models

import (
	"github.com/gogf/gf/os/glog"
	"homeproxy/library/filedb"
)

const ProxyServerTable = "proxy_server"

func init() {
	settings := filedb.DefaultCollectionSettings()
	settings.Unique = "address"
	if err := filedb.DB.NewCollections(ProxyInstanceTable, settings); err != nil {
		if err != filedb.ErrCollectionExist {
			glog.Error("init collection error: %s", err.Error())
		}
	}

	if err := filedb.DB.NewCollections(ProxyServerTable, nil); err != nil {
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
		c, _ := filedb.DB.Collection(ProxyServerTable)
		_, err := c.Insert(server)
		if err != nil {
			glog.Errorf("init server info error: %s", err)
		}
	}
}

type ProxyServer struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Port      int    `json:"port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Status    bool   `json:"status"`
	AutoProxy bool   `json:"auto_proxy"`
	AllProxy  bool   `json:"all_proxy"`
	AutoStart bool   `json:"auto_start"`
}

func GetProxyServer() (server ProxyServer, err error) {
	var c *filedb.Collection
	if c, err = filedb.DB.Collection(ProxyServerTable); err != nil {
		return
	} else {
		err = c.GetFirst(&server)
		return server, err
	}
}
