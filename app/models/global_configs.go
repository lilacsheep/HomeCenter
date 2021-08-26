package models

import (
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
