package models

import (
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gmeta"
)

type GlobalConfig struct {
	DefaultModel
	Group      string `json:"group"`
	Key        string `json:"key"`
	Type       string `json:"type"`
	Value      string `json:"value"`
	Desc       string // 描述
	gmeta.Meta `orm:"table:global_configs"`
}

func (o GlobalConfig) Var() *g.Var {
	return g.NewVar(o.Value, true)
}

func UpdateConfig(group, key string, value interface{}) error {
	_, err := g.DB().Model(&GlobalConfig{}).Data(
		g.Map{"value": gconv.String(value)}).Where("`group` = ? AND `key` = ?", group, key).Update()
	return err
}

func GetGroupConfigs(group string) (configs []GlobalConfig, err error) {
	err = g.DB().Model(&GlobalConfig{}).Where("`group` = ?", group).Structs(&configs)
	return
}

func GetConfigsMap(group string) (map[string]string, error) {
	configs, err := GetGroupConfigs(group)
	if err != nil {
		return nil, err
	}
	r := make(map[string]string)

	for _, c := range configs {
		r[c.Key] = c.Value
	}
	return r, nil
}

func ConfigToStruct(group string, to interface{}) error {
	configs, err := GetConfigsMap(group)
	if err != nil {
		return err
	}
	return gjson.New(configs).Struct(to)
}
