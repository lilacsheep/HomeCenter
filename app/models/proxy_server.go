package models

import (
	"homeproxy/library/filedb2"
)

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
	err := filedb2.DB.Get("settings", "server", &server)
	return &server, err
}
