package models

import (
	"time"

	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
)

type DownloadSettings struct {
	Aria2Url            string `json:"aria2_url"`              // aria2地址
	Aria2Token          string `json:"aria2_token"`            // aria2的Token
	AutoClean           int    `json:"auto_clean"`             // 自动清理Bt下载后文件夹内内容，根据文件大小判断
	AutoUpdateBTTracker string `json:"auto_update_bt_tracker"` // 自动更新bt-tracker, "" 为不更新,
}

func GetSettings() (*DownloadSettings, error) {
	query := g.DB().Model(&GlobalConfig{})
	c, err := query.Where("`group` = ?", "download").Count()
	if err != nil {
		return nil, err
	}
	if c == 0 {
		query.Data(g.List{
			{"group": "download", "key": "aria2_url", "type": "string", "value": "", "desc": "aria2地址"},
			{"group": "download", "key": "aria2_token", "type": "string", "value": "", "desc": "aria2的Token"},
			{"group": "download", "key": "auto_clean", "type": "int", "value": "0", "desc": "自动清理Bt下载后文件夹内内容，根据文件大小判断"},
			{"group": "download", "key": "auto_update_bt_tracker", "type": "string", "value": "", "desc": "自动更新bt-tracker,空为不更新"},
		}).Save()
	}
	configs, err := GetConfigsMap("download")
	if err != nil {
		return nil, err
	}
	settings := &DownloadSettings{}
	err = gjson.New(configs).Struct(settings)
	return settings, err
}

type DownloadFileList struct {
	ID       int       `json:"id" storm:"id,increment"`
	Vkey     string    `json:"vkey"`
	Path     string    `json:"path"`
	CreateAt time.Time `json:"create_at"`
}
