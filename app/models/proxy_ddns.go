package models

const (
	DDnsProviderSettingsTable string = "ddns_provider_settings"
)

type OperationRecord struct {
	Date   string `json:"date"`
	Status int    `json:"status"`
	Error  string `json:"error"`
	Value  string `json:"value"`
}

type DDnsProviderSettings struct {
	ID           string            `json:"id"`
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
	History      []OperationRecord `json:"history"`
}
