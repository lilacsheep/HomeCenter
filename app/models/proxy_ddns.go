package models

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gmeta"
)

type OperationRecord struct {
	DefaultModel
	SettingId  int    `json:"setting_id"`
	Date       string `json:"date"`
	Status     int    `json:"status"`
	Error      string `json:"error"`
	Value      string `json:"value"`
	gmeta.Meta `orm:"table:ddns_operation_record"`
}

type DDnsProviderSettings struct {
	DefaultModel
	Provider     string            `json:"provider"`
	Domain       string            `json:"domain"`
	SubDomain    string            `json:"sub_domain"`
	TimeInterval string            `json:"time_interval"`
	UsePublicIP  bool              `json:"use_public_ip"`
	NetCard      int               `json:"net_card"`
	RecordID     string            `json:"record_id"`
	DNSPodID     string            `json:"dnspod_id"`
	DNSPodToken  string            `json:"dnspod_token"`
	Status       bool              `json:"status"`
	UpdatedOn    string            `json:"updated_on"`
	History      []OperationRecord `json:"history" orm:"-"`
	gmeta.Meta   `orm:"table:ddns_operation_settings"`
}

func AllDDnsSettings() (settings []DDnsProviderSettings, err error) {
	err = g.DB().Model(&DDnsProviderSettings{}).Structs(&settings)
	return
}
