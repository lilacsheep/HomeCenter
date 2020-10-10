package requests

import (
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"homeproxy/app/models"
	"homeproxy/app/server"
	"net/http"
)

type StartProxyServerRequest struct{}

func (self *StartProxyServerRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	err := server.Mallory.Start()
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err)
	} else {
		response.Success()
	}
	return
}

func NewStartProxyServerRequest() *StartProxyServerRequest {
	return &StartProxyServerRequest{}
}

type StopProxyServerRequest struct{}

func (self *StopProxyServerRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	if err := server.Mallory.Stop(); err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		response.Success()
	}
	return
}

func NewStopProxyServerRequest() *StopProxyServerRequest {
	return &StopProxyServerRequest{}
}

type UpdateProxyServerRequest struct {
	Name      string `json:"name"`
	Port      int    `json:"port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Status    bool   `json:"status"`
	AutoProxy bool   `json:"auto_proxy"`
	AllProxy  bool   `json:"all_proxy"`
}

func (self *UpdateProxyServerRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	proxy, err := models.GetProxyServer()
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err)
	} else {
		data := g.Map{}
		if self.Name != "" && proxy.Name != self.Name {
			data["name"] = self.Name
		}
		if self.Port != 0 && proxy.Port != self.Port {
			data["port"] = self.Port
		}
		if self.Username != "" && proxy.Username == self.Username {
			data["username"] = self.Username
		}
		if self.Password != "" && proxy.Password == self.Password {
			data["password"] = self.Password
		}
		if self.Status != proxy.Status {
			data["status"] = self.Status
			if self.Status {
				server.Mallory.SetBalance(1)
			} else {
				server.Mallory.SetBalance(0)
			}
		}
		if self.AutoProxy != proxy.AutoProxy {
			data["auto_proxy"] = self.AutoProxy
		}
		if self.AllProxy != proxy.AllProxy {
			data["all_proxy"] = self.AllProxy
			if server.Mallory.Status {
				server.Mallory.ProxyHandler.AllProxy = self.AllProxy
			}
		}
		c, _ := models.DB.Collection(models.ProxyServerTable)
		c.UpdateById(proxy.ID, data)
		response.Success()
	}
	return
}

func NewUpdateProxyServerRequest() *UpdateProxyServerRequest {
	return &UpdateProxyServerRequest{}
}

type InfoProxyServerRequest struct{}

func (self *InfoProxyServerRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	info, err := models.GetProxyServer()
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		data := garray.New(true)
		data.Append(g.Map{"key": "name", "name": "名称", "value": info.Name})
		data.Append(g.Map{"key": "port", "name": "端口", "value": info.Port})
		data.Append(g.Map{"key": "status", "name": "状态", "value": server.Mallory.Status})
		data.Append(g.Map{"key": "balance", "name": "负载", "value": info.Status})
		//TODO: 自动代理模式，
		//data.Append(g.Map{"key": "auto_proxy", "name": "代理", "value": info.AutoProxy})
		data.Append(g.Map{"key": "all_proxy", "name": "模式", "value": info.AllProxy})
		if server.Mallory.Status {
			data.Append(g.Map{"key": "instances", "name": "实例", "value": server.Mallory.Instances.Size()})
		}
		if server.Mallory.Error != nil && server.Mallory.Error != http.ErrServerClosed {
			data.Append(g.Map{"key": "error", "name": "错误", "value": server.Mallory.Error.Error()})
		}
		response.SuccessWithDetail(data)
	}

	return
}

func NewInfoProxyServerRequest() *InfoProxyServerRequest {
	return &InfoProxyServerRequest{}
}
