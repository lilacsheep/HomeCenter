package router

import (
	"homeproxy/app/api"
	"homeproxy/app/middleware"
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
		"/other":      "/",
		"/login":      "/",
		"/users":      "/",
		"/ddns":       "/",
		"/containers": "/",
		"/webssh": "/",
	})
	s.Use(MiddlewareCORS)
	proxyInstanceApi := &api.ProxyInstanceApi{}
	proxyServerApi := &api.ProxyServerApi{}
	proxyRoleApi := &api.ProxyRoleApi{}
	auth := &api.AuthController{}
	common := &api.CommonApi{}
	downloadApi := &api.ProxyDownloadApi{}

	s.BindHandler("POST:/api/login", auth.LoginUser)
	s.BindHandler("Get:/download/:vkey", downloadApi.Download)
	s.Group("/api", func(group *ghttp.RouterGroup) {
		group.Middleware(AuthMiddleware)

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

		// proxy url role
		group.POST("/proxy/role/add", proxyRoleApi.AddRole)
		group.POST("/proxy/role/remove", proxyRoleApi.Remove)
		group.POST("/proxy/role/change", proxyRoleApi.Change)
		group.GET("/proxy/roles", proxyRoleApi.All)

		// download api
		group.Group("/download", func(group *ghttp.RouterGroup) {
			group.Middleware(middleware.CheckAria2Cli)
			group.POST("/create", downloadApi.Create)
			group.POST("/torrent", downloadApi.AddTorrent)
			group.GET("/tasks", downloadApi.Query)
			group.POST("/remove", downloadApi.Remove)
			group.POST("/task/pause", downloadApi.Pause)
			group.POST("/task/unpause", downloadApi.UnPause)
			group.POST("/task/status", downloadApi.TaskStatus)
			group.GET("/global/stat", downloadApi.GlobalStatInfo)
			group.GET("/global/options", downloadApi.Options)
			group.POST("/make/download", downloadApi.MakeDownloadUrl)
		})

		// download settings api
		group.GET("/download/settings", downloadApi.Settings)
		group.POST("/download/settings/update", downloadApi.UpdateSettings)

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

		// system
		systemMonApi := new(api.SystemApi)
		group.GET("/system/info", systemMonApi.Info)
		group.GET("/system/processes", systemMonApi.Processes)
		group.ALL("/system/webssh", systemMonApi.Webssh)

		// common api
		group.GET("/common/countrys", common.Countrys)

		// docker containers
		dockerContainerApi := new(api.ContainersController)
		group.Group("/docker/container", func(group *ghttp.RouterGroup) {
			group.Middleware(middleware.CheckDockerCli)
			group.POST("/list", dockerContainerApi.List)
			group.POST("/start", dockerContainerApi.Start)
			group.POST("/stop", dockerContainerApi.Stop)
			group.POST("/restart", dockerContainerApi.Restart)
			group.POST("/update", dockerContainerApi.Update)
			group.POST("/pause", dockerContainerApi.Pause)
			group.POST("/unpause", dockerContainerApi.Unpause)
			group.POST("/info", dockerContainerApi.Info)
			group.POST("/stats", dockerContainerApi.Stats)
			group.POST("/rename", dockerContainerApi.Rename)
		})

		dockerImageApi := new(api.ImageController)
		group.Group("/docker/image", func(group *ghttp.RouterGroup) {
			group.Middleware(middleware.CheckDockerCli)
			group.POST("/list", dockerImageApi.List)
			group.POST("/history", dockerImageApi.History)
			group.POST("/info", dockerImageApi.Info)
			group.POST("/pull", dockerImageApi.Pull)
			group.POST("/remove", dockerImageApi.Remove)
		})

		dockerVolumeApi := new(api.VolumeApi)
		group.Group("/docker/volume", func(group *ghttp.RouterGroup) {
			group.Middleware(middleware.CheckDockerCli)
			group.POST("/list", dockerVolumeApi.List)
			group.POST("/create", dockerVolumeApi.Create)
			group.POST("/remove", dockerVolumeApi.Remove)
		})

		serverApi := new(api.ServerApi)
		group.Group("/server", func(group *ghttp.RouterGroup) {
			group.POST("/list", serverApi.ServerList)
			group.POST("/create", serverApi.ServerCreate)
			group.POST("/update", serverApi.ServerUpdate)
			group.POST("/remove", serverApi.ServerDelete)
			group.POST("/group/list", serverApi.GroupList)
			group.POST("/group/remove", serverApi.GroupRemove)
			group.POST("/group/create", serverApi.GroupCreate)
			group.POST("/group/update", serverApi.GroupUpdate)
		})
	})
}
