package models

import (
	"homeproxy/library/filedb"
	"homeproxy/library/filedb2"

	"github.com/gogf/gf/os/glog"
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
}

type ProxyServer struct {
	ID        int    `json:"id" storm:"id,increment"`
	Name      string `json:"name"`
	Port      int    `json:"port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Status    bool   `json:"status"`
	AutoProxy bool   `json:"auto_proxy"`
	AllProxy  bool   `json:"all_proxy"`
	AutoStart bool   `json:"auto_start"`
}

func GetProxyServer() (*ProxyServer, error) {
	server := ProxyServer{}
	err := filedb2.DB.One("Name", "default", &server)
	return &server, err
}
