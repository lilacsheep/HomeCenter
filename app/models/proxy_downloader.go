package models

import (
	"homeproxy/library/filedb2"
	"time"
)

type DownloadSettings struct {
	ID                  int    `json:"id" storm:"id,increment"`
	Aria2Url            string `json:"aria2_url"`              // aria2地址
	Aria2Token          string `json:"aria2_token"`            // aria2的Token
	AutoClean           int    `json:"auto_clean"`             // 自动清理Bt下载后文件夹内内容，根据文件大小判断
	AutoUpdateBTTracker string `json:"auto_update_bt_tracker"` // 自动更新bt-tracker, "" 为不更新,
}

func GetSettings() (*DownloadSettings, error) {
	settings := &DownloadSettings{}
	err := filedb2.DB.Get("settings", "download", settings)
	return settings, err
}

type DownloadFileList struct {
	ID       int       `json:"id" storm:"id,increment"`
	Vkey     string    `json:"vkey"`
	Path     string    `json:"path"`
	CreateAt time.Time `json:"create_at"`
}
