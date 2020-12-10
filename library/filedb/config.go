package filedb

import (
	"homeproxy/library/config"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/glog"
)

var (
	DB      *Database
)

func init() {
	if !gfile.Exists(config.Dbpath) {
		_ = gfile.Mkdir(config.Dbpath)
	}

	glog.Debugf("init db: %s path: %s", config.Dbname, config.Dbpath)
	DB = NewDatabase(config.Dbname, config.Dbpath)
}

type CollectionSettings struct {
	MaxRecord int    `json:"max_record"` // 集合记录最大记录数
	AutoDump  bool   `json:"auto_dump"`  // 需要自动文件保存
	Rule      g.Map  `json:"rule"`       // TODO: 字段校验规则 目前以gf框架的gvalid模块为主 https://goframe.org/util/gvalid/index
	Unique    string `json:"unique"`     // 非重复字段
}

func DefaultCollectionSettings() *CollectionSettings {
	return &CollectionSettings{
		MaxRecord: 0,
		AutoDump:  true,
		Rule:      nil,
		Unique:    "",
	}
}
