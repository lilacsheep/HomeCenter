package config

import (
	"flag"
)

var (
	Dbname  string
	Dbpath  string
	WebHost string
	DataDir string
)

func init() {
	flag.StringVar(&Dbname, "name", "default", "数据库名称,默认为: default")
	flag.StringVar(&Dbpath, "path", "db", "数据路径, 默认 ./db")
	flag.StringVar(&WebHost, "h", "0.0.0.0:8080", "监听地址,默认为0.0.0.0:8080")
	flag.StringVar(&DataDir, "d", "data", "数据上传路径")
	flag.Parse()
}