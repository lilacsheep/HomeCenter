package main

import (
	_ "homeproxy/app/server"
	_ "homeproxy/app/server/aria2"
	"homeproxy/boot"
	"homeproxy/library/config"
	"homeproxy/library/events"
	_ "homeproxy/packed"
	_ "homeproxy/router"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/grpool"
	"github.com/gogf/gf/os/gsession"
)

func main() {
	eventProcess := events.EventProcess{Pool: grpool.New(10)}
	go eventProcess.Run()

	s := g.Server()
	s.SetConfigWithMap(g.Map{
        "SessionMaxAge":  time.Minute * 30,
        "SessionStorage": gsession.NewStorageMemory(),
	})
	
	s.SetIndexFolder(true)
	s.SetClientMaxBodySize(2199023255552)
	s.SetRouteOverWrite(true)
	s.SetServerRoot("public")
	s.SetAddr(config.WebHost)
	s.AddStaticPath("/static", "public")
	// s.SetDumpRouterMap(false)

	boot.Setup()
	s.Run()
}
