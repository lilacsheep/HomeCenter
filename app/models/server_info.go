package models

import "github.com/gogf/gf/util/gmeta"

type ServerGroup struct {
	DefaultModel
	Name       string `json:"name"`
	Remark     string `json:"remark"`
	gmeta.Meta `orm:"table:server_group"`
}

type Server struct {
	DefaultModel
	Name       string `json:"name"`
	Address    string `json:"address"`
	Port       int    `json:"port"`
	Group      int    `json:"group"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	PrivateKey string `json:"private_key"`
	Remark     string `json:"remark"`
	Status     bool   `json:"status"`
	gmeta.Meta `orm:"table:server_host"`
}
