package models

import (
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/net/ghttp"
	"homeproxy/library/filedb2"
)

const (
	ProxyNotifyMessageTable = "proxy_notify_message"
	BarkServerUrl           = "https://api.day.app"
)

type NotifyMessage struct {
	ID                string `json:"id"`
	Title             string `json:"title"`
	Message           string `json:"message"`
	AutomaticallyCopy int    `json:"automaticallyCopy"` //指定是否需要保存推送信息到历史记录，1 为保存，其他值为不保存。如果不指定这个参数，推送信息将按照APP内设置来决定是否保存。
	IsActive          int    `json:"isActive"`
	Copy              string `json:"copy"`
	Status            bool   `json:"status"`
	CreateAt          string `json:"create_at"`
}

func (self *NotifyMessage) Push(key string) error {
	client := ghttp.NewClient()
	body := gjson.New(self)
	body.Set("key", key)
	_ = body.Remove("create_at")
	resp, err := client.Post(BarkServerUrl, body.MustToJson())
	if err != nil {
		return err
	}
	defer resp.Close()
	self.Status = true
	return self.sync()
}

func (self *NotifyMessage) sync() error {
	return filedb2.DB.Update(self)
}
