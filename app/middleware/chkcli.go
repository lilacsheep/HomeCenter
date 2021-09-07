package middleware

import (
	"homeproxy/app/server/aria2"
	"homeproxy/library/docker"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func CheckDockerCli(r *ghttp.Request) {
	if docker.Docker == nil {
		r.Response.WriteJsonExit(g.Map{"code": 0, "message": "失败", "detail": "docker 客户端未初始化"})
	}
	r.Middleware.Next()
}

func CheckAria2Cli(r *ghttp.Request) {
	if aria2.Manager == nil {
		r.Response.WriteJsonExit(g.Map{"code": 0, "message": "失败", "detail": "ariar2 客户端未初始化"})
	}
	r.Middleware.Next()
}
