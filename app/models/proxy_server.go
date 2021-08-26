package models

import (
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
)

type ProxyServer struct {
	Port       int    `json:"port"`
	Username   string `json:"username"`
	DNSAddr    string `json:"dns_addr"`
	Password   string `json:"password"`
	Balance    bool   `json:"balance"`
	ProxyMode  int    `json:"proxy_mode"`
	AutoStart  bool   `json:"auto_start"`
	EnableAuth bool   `json:"enable_auth"`
}

func getProxyServer() (*ProxyServer, error) {
	var configs []GlobalConfig
	err := g.DB().Model(&GlobalConfig{}).Where("group", "mallory").Structs(&configs)
	if err != nil {
		return nil, err
	}
	configMap := make(map[string]string)
	server := &ProxyServer{}

	for _, c := range configs {
		configMap[c.Key] = c.Value
	}
	err = gjson.New(configMap).Struct(server)
	return server, err
}

func DefaultProxyServer() (*ProxyServer, error) {
	c, err := g.DB().Model(&GlobalConfig{}).Where("group", "mallory").Count()
	if err != nil {
		return nil, err
	}
	if c == 0 {
		g.DB().Model(&GlobalConfig{}).Data(g.List{
			{"group": "mallory", "key": "port", "type": "int", "value": "1316", "desc": "代理启动端口"},
			{"group": "mallory", "key": "username", "type": "string", "value": "", "desc": "HttpBasic用户名"},
			{"group": "mallory", "key": "password", "type": "string", "value": "", "desc": "HttpBasic密码"},
			{"group": "mallory", "key": "proxy_mode", "type": "int", "value": "1", "desc": "代理模式"},
			{"group": "mallory", "key": "auto_start", "type": "bool", "value": "true", "desc": "代理自动启动"},
			{"group": "mallory", "key": "enable_auth", "type": "bool", "value": "false", "desc": "代理开启HttpBasic认证"},
			{"group": "mallory", "key": "dns_addr", "type": "string", "value": "114.114.114.114", "desc": "DNS地址"},
			{"group": "mallory", "key": "balance", "type": "bool", "value": "false", "desc": "负载均衡模式"},
		}).Save()
	}

	return getProxyServer()
}
