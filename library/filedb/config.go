package filedb

import (
	"flag"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/glog"
)

var (
	Dbname  string
	Dbpath  string
	WebHost string
	DB      *Database
)

func init() {
	flag.StringVar(&Dbname, "name", "default", "数据库名称,默认为: default")
	flag.StringVar(&Dbpath, "path", "db", "数据路径, 默认 ./db")
	flag.StringVar(&WebHost, "h", "0.0.0.0:8080", "监听地址,默认为0.0.0.0:8080")

	flag.Parse()
	if !gfile.Exists(Dbpath) {
		_ = gfile.Mkdir(Dbpath)
	}

	glog.Debugf("init db: %s path: %s", Dbname, Dbpath)
	DB = NewDatabase(Dbname, Dbpath)
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
