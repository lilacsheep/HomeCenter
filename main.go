package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/grpool"
	_ "homeproxy/app/server"
	"homeproxy/library/events"
	"homeproxy/library/filedb"
	_ "homeproxy/packed"
	_ "homeproxy/router"
)

func main() {
	eventProcess := events.EventProcess{Pool: grpool.New(10)}
	go eventProcess.Run()

	s := g.Server()
	s.SetIndexFolder(true)
	s.SetClientMaxBodySize(2199023255552)
	s.SetRouteOverWrite(true)
	s.SetServerRoot("public")
	s.SetAddr(filedb.WebHost)
	s.AddStaticPath("/static", "public")
	s.SetDumpRouterMap(false)
	s.Run()
}
