package models

import (
	"homeproxy/library/filedb"
)

const ProxyServerTable = "proxy_server"

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
	if c, err = DB.Collection(ProxyServerTable); err != nil {
		return
	} else {
		err = c.GetFirst(&server)
		return server, err
	}
}
