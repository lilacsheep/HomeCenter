package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/grpool"
	"homeproxy/app/api"
	_ "homeproxy/app/models"
	_ "homeproxy/app/server"
	"homeproxy/library/events"
)

func MiddlewareCORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

func main() {
	eventProcess := events.EventProcess{Pool: grpool.New(10)}
	go eventProcess.Run()

	s := g.Server()
	s.SetIndexFolder(true)
	s.SetRouteOverWrite(true)
	s.AddStaticPath("/static", "public")
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.WriteTpl("index.html")
	})

	s.SetRewriteMap(map[string]string{
		"/dashboard": "/",
	})
	proxyInstanceApi := &api.ProxyInstanceApi{}
	proxyServerApi := &api.ProxyServerApi{}
	proxyRoleApi := &api.ProxyRoleApi{}

	s.Group("/api", func(group *ghttp.RouterGroup) {
		group.Middleware(MiddlewareCORS)

		//group.POST("/login", new(api.AuthController).LoginUser)
		//group.POST("/auth/user/create", new(api.AuthController).CreateUser)
		//group.GET("/auth/self", new(api.AuthController).MySelf)
		// proxy instance api
		group.POST("/proxy/instance/create", proxyInstanceApi.Create)
		group.GET("/proxy/instances", proxyInstanceApi.Query)
		// proxy server api
		group.POST("/proxy/server/start", proxyServerApi.Start)
		group.POST("/proxy/server/stop", proxyServerApi.Stop)
		group.POST("/proxy/server/update", proxyServerApi.Update)
		group.GET("/proxy/server/info", proxyServerApi.Info)
		// proxy url role
		group.POST("/proxy/role/add", proxyRoleApi.AddRole)
		group.POST("/proxy/role/remove", proxyRoleApi.Remove)
		group.GET("/proxy/roles", proxyRoleApi.All)
	})
	s.Run()
}
