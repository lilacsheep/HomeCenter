package models

import "github.com/gogf/gf/frame/g"

const (
	LocalMinioServer = 0
	S3Server         = 1
)

type LocalMinioOption struct {
	Public        bool   `json:"public"`
	AutoStart     bool   `json:"auto_start"`
	Port          int    `json:"port"`
	WebUi         bool   `json:"webui"`
	WebUiPort     int    `json:"webui_port"`
	SavePath      string `json:"save_path"`
	Region        string `json:"region"`
	RegionComment string `json:"region_comment"`
	ConfigDir     string `json:"config_dir"`
	MinioDomain   string `json:"minio_domain"`
}

type ObjectServer struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Type      int    `json:"type"`
	Region    string `json:"region"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	EndPoint  string `json:"end_point"`
	Extend    string `json:"extend"`
}

func GetLocalMinioServerConfig() (*ObjectServer, error) {
	server := new(ObjectServer)
	err := g.DB().Model(new(ObjectServer)).Where("`type` = ?", LocalMinioServer).Scan(server)
	if err != nil {
		return nil, err
	}
	return server, nil
}
