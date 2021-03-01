package router

import (
	"homeproxy/app/api"
	"net/http"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func MiddlewareCORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

func AuthMiddleware(r *ghttp.Request) {
	if user := r.Session.Get("user"); user == nil {
		r.Response.WriteStatus(http.StatusUnauthorized)
	} else {
		r.Middleware.Next()
	}
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
		"/login":      "/",
		"/user":      "/",
	})
	proxyInstanceApi := &api.ProxyInstanceApi{}
	proxyServerApi := &api.ProxyServerApi{}
	proxyRoleApi := &api.ProxyRoleApi{}
	auth := &api.AuthController{}
	common := &api.CommonApi{}

	s.BindHandler("POST:/api/login", auth.LoginUser)
	s.Group("/api", func(group *ghttp.RouterGroup) {
		group.Middleware(MiddlewareCORS, AuthMiddleware)
		// user auth
		group.POST("/auth/self", auth.Self)
		group.GET("/auth/logout", auth.Logout)
		group.POST("/auth/create/user", auth.CreateUser)
		group.POST("/auth/change/self/password", auth.ChangeSelfPassword)
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
		group.POST("/proxy/role/change", proxyRoleApi.Change)
		group.GET("/proxy/roles", proxyRoleApi.All)
		group.GET("/proxy/logs", proxyRoleApi.RoleVisitLogs)
		group.POST("/proxy/log/add", proxyRoleApi.AddLogToRole)

		// download api
		downloadApi := &api.ProxyDownloadApi{}
		group.POST("/download/create", downloadApi.Create)
		group.POST("/download/torrent", downloadApi.AddTorrent)
		group.GET("/download/tasks", downloadApi.Query)
		group.POST("/download/remove", downloadApi.Remove)
		group.POST("/download/task/pause", downloadApi.Pause)
		group.POST("/download/task/unpause", downloadApi.UnPause)
		group.POST("/download/task/status", downloadApi.TaskStatus)
		group.GET("/download/global/stat", downloadApi.GlobalStatInfo)
		group.GET("/download/global/options", downloadApi.Options)

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

		// logs api
		group.GET("/logs", proxyServerApi.Logs)

		// common api
		group.GET("/common/countrys", common.Countrys)
	})
}
