package mallory

import (
	"time"
)

type ProxyRole struct {
	ID         int    `json:"id" storm:"id,increment"`
	InstanceID int    `json:"instance_id"`
	Status     bool   `json:"status"`
	Sub        string `json:"sub"`
	Domain     string `json:"domain"`
}

type ProxyRoleAnalysis struct {
	ID     int    `json:"id" storm:"id,increment"`
	Domain string `json:"domain" storm:"unique"`
	Times  int    `json:"times"`
	Error  string `json:"error"`
}

type ProxyVisitLog struct {
	ID       int       `json:"id" storm:"id,increment"`
	Address  string    `json:"address" storm:"index"`
	Host     string    `json:"host" storm:"index"`
	Datetime time.Time `json:"datetime"`
}
