package api

import (
	"homeproxy/app/services/requests"

	"github.com/gogf/gf/net/ghttp"
)

type ServerApi struct {
	BaseControllers
}

func (api *ServerApi) ServerList(r *ghttp.Request) {
	request := requests.ServerListRequest{}
	api.DoRequestValid(&request, r)
}

func (api *ServerApi) ServerCreate(r *ghttp.Request) {
	request := requests.CreateServerRequest{}
	api.DoRequestValid(&request, r)
}

func (api *ServerApi) ServerDelete(r *ghttp.Request) {
	request := requests.ServerDeleteRequest{}
	api.DoRequestValid(&request, r)
}

func (api *ServerApi) ServerUpdate(r *ghttp.Request) {
	request := requests.ServerUpdateRequest{}
	api.DoRequestValid(&request, r)
}

func (api *ServerApi) GroupList(r *ghttp.Request) {
	request := requests.ServerGroupListRequest{}
	api.DoRequestValid(&request, r)
}

func (api *ServerApi) GroupCreate(r *ghttp.Request) {
	request := requests.CreateServerGroupRequest{}
	api.DoRequestValid(&request, r)
}

func (api *ServerApi) GroupRemove(r *ghttp.Request) {
	request := requests.RemoveServerGroupRequest{}
	api.DoRequestValid(&request, r)
}