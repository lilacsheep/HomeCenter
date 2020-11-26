package router

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"homeproxy/app/api"
)

func MiddlewareCORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

func init() {
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.WriteTpl("index.html")
	})
	s.SetRewriteMap(map[string]string{
		"/dashboard":  "/",
		"/roles":      "/",
		"/monitor":    "/",
		"/download":   "/",
		"/filesystem": "/",
		"/message":    "/",
		"/other":      "/",
	})
	proxyInstanceApi := &api.ProxyInstanceApi{}
	proxyServerApi := &api.ProxyServerApi{}
	proxyRoleApi := &api.ProxyRoleApi{}

	s.Group("/api", func(group *ghttp.RouterGroup) {
		group.Middleware(MiddlewareCORS)
		// proxy instance api
		group.POST("/proxy/instance/create", proxyInstanceApi.Create)
		group.POST("/proxy/instance/update", proxyInstanceApi.Update)
		group.POST("/proxy/instance/remove", proxyInstanceApi.Remove)
		group.POST("/proxy/instance/pool/remove", proxyInstanceApi.RemoveFromPool)
		group.POST("/proxy/instance/pool/add", proxyInstanceApi.IntoPool)
		group.GET("/proxy/instances", proxyInstanceApi.Query)
		// proxy server api
		group.POST("/proxy/server/start", proxyServerApi.Start)
		group.POST("/proxy/server/stop", proxyServerApi.Stop)
		group.POST("/proxy/server/update", proxyServerApi.Update)
		group.GET("/proxy/server/info", proxyServerApi.Info)
		group.GET("/proxy/server/monitor", proxyServerApi.Monitor)
		// proxy url role
		group.POST("/proxy/role/add", proxyRoleApi.AddRole)
		group.POST("/proxy/role/remove", proxyRoleApi.Remove)
		group.GET("/proxy/roles", proxyRoleApi.All)

		// download api
		downloadApi := &api.ProxyDownloadApi{}
		group.POST("/download/create", downloadApi.Create)
		group.GET("/download/tasks", downloadApi.Query)
		group.POST("/download/remove", downloadApi.Remove)
		group.POST("/download/cancel", downloadApi.Cancel)
		group.POST("/download/start", downloadApi.Start)

		// download settings api
		group.GET("/download/settings", downloadApi.Settings)
		group.POST("/download/settings/update", downloadApi.UpdateSettings)

		// filesystem api
		filesystemApi := &api.ProxyFilesystemApi{}
		group.GET("/filesystem/nodes", filesystemApi.Nodes)
		group.POST("/filesystem/files", filesystemApi.Files)
		group.POST("/filesystem/file/remove", filesystemApi.RemoveFile)
		group.POST("/filesystem/file/upload", filesystemApi.UploadFile)
		group.POST("/filesystem/node/create", filesystemApi.CreateNode)
		group.POST("/filesystem/node/remove", filesystemApi.RemoveNode)
		group.POST("/filesystem/dir/create", filesystemApi.CreateDir)
		group.POST("/filesystem/dir/remove", filesystemApi.RemoveDir)
		group.GET("/filesystem/download", filesystemApi.DownloadFile)
		group.POST("/filesystem/file/info", filesystemApi.FileInfo)

		// message api
		messageApi := new(api.ProxyMessageApi)
		group.GET("/messages", messageApi.All)

		// ddns api
		ddnsApi := new(api.ProxyDDnsApi)
		group.GET("/ddns/settings", ddnsApi.GetSettings)
		group.GET("/ddns/netcards", ddnsApi.ChooseCards)
		group.POST("/ddns/setting/create", ddnsApi.CreateSetting)
		group.POST("/ddns/records", ddnsApi.GetRecords)
		group.POST("/ddns/setting", ddnsApi.GetSettingInfo)
		group.POST("/ddns/setting/start", ddnsApi.StartSetting)
		group.POST("/ddns/setting/stop", ddnsApi.StopSetting)
		group.POST("/ddns/setting/remove", ddnsApi.DeleteSetting)
		group.POST("/ddns/setting/refresh", ddnsApi.RefreshSetting)
	})
}
