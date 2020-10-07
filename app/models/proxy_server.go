package models

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
)

const ProxyServerTable = "proxy_server"

func init() {
	server := ProxyServer{
		Name:      "default",
		Port:      1316,
		Status:    true,
		AutoProxy: false,
	}
	querySet := g.DB().Table(ProxyServerTable)

	if c, err := querySet.Count(); err != nil {
		glog.Errorf("init proxy server error: %s", err)
	} else {
		if c == 0 {
			_, err := querySet.Insert(&server)
			if err != nil {
				glog.Errorf("init proxy server error: %s", err)
			}
		}
	}
}

type ProxyServer struct {
	Name      string `orm:"name"`
	Port      int    `orm:"port"`
	Username  string `orm:"username"`
	Password  string `orm:"password"`
	Status    bool   `orm:"status"`
	AutoProxy bool   `orm:"auto_proxy"`
}

func GetProxyServer() (server ProxyServer, err error) {
	if err = g.DB().Table(ProxyServerTable).Struct(&server); err != nil {
		return
	}
	return
}
