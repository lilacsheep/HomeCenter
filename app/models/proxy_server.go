package models

import (
	"homeproxy/library/filedb2"
)

type ProxyServer struct {
	ID        int    `json:"id" storm:"id,increment"`
	Name      string `json:"name"`
	Port      int    `json:"port"`
	Username  string `json:"username"`
	DNSAddr   string `json:"dns_addr"`
	Password  string `json:"password"`
	Status    bool   `json:"status"`
	ProxyMode int    `json:"proxy_mode"`
	AutoStart bool   `json:"auto_start"`
}

func GetProxyServer() (*ProxyServer, error) {
	server := ProxyServer{}
	err := filedb2.DB.Get("settings", "server", &server)
	return &server, err
}
