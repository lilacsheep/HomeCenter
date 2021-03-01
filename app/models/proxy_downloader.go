package models

import (
	"github.com/gogf/gf/os/glog"
	"homeproxy/library/filedb2"
)

func init() {
	if found, _ := filedb2.DB.KeyExists("settings", "download"); !found {
		if !found {
			settings := DownloadSettings{}
			filedb2.DB.Set("settings", "download", &settings)
			glog.Info("init download settings")
		}
	}
}

type DownloadSettings struct {
	ID              int    `json:"id" storm:"id,increment"`
	Aria2Url        string `json:"aria2_url"`   // aria2地址
	Aria2Token      string `json:"aria2_token"` // aria2的Token
}

func GetSettings() (*DownloadSettings, error) {
	settings := &DownloadSettings{}
	err := filedb2.DB.Get("settings", "download", settings)
	return settings, err
}
