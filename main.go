package main

import (
	"flag"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/grpool"
	"homeproxy/app/models"
	_ "homeproxy/app/server"
	"homeproxy/library/events"
	_ "homeproxy/packed"
	_ "homeproxy/router"
)

var (
	host string
)

func init() {
	flag.StringVar(&models.Dbname, "name", "default", "数据库名称,默认为: default")
	flag.StringVar(&models.Dbpath, "path", "db", "数据路径, 默认 ./db")
	flag.StringVar(&host, "h", "0.0.0.0:80", "监听地址,默认为0.0.0.0:80")
	flag.Parse()
}

func main() {
	models.InitDB()
	eventProcess := events.EventProcess{Pool: grpool.New(10)}
	go eventProcess.Run()

	s := g.Server()
	s.SetIndexFolder(true)
	s.SetRouteOverWrite(true)
	s.SetServerRoot("public")
	s.SetAddr(host)
	s.SetDumpRouterMap(false)
	s.AddStaticPath("/static", "public")

	s.Run()
}
