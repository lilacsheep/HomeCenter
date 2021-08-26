package requests

import (
	"homeproxy/app/models"
	"homeproxy/app/server"
	"homeproxy/library/filedb2"
	"net/http"

	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

type StartProxyServerRequest struct{}

func (self *StartProxyServerRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	err := server.Mallory.Start()
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
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
	Name       string `json:"name"`
	Port       int    `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Balance    bool   `json:"balance"`
	DNSAddr    string `json:"dns_addr"`
	AutoProxy  bool   `json:"auto_proxy"`
	ProxyMode  int    `json:"proxy_mode"`
	EnableAuth bool   `json:"enable_auth"`
}

func (self *UpdateProxyServerRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	server2, err := models.DefaultProxyServer()
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err)
	} else {
		if self.Port != 0 && server2.Port != self.Port {
			server2.Port = self.Port
		}
		if self.Username != "" && server2.Username == self.Username {
			server2.Username = self.Username
		}
		if self.Password != "" && server2.Password == self.Password {
			server2.Password = self.Password
		}
		if self.DNSAddr != "" && server2.DNSAddr != self.DNSAddr {
			server2.DNSAddr = self.DNSAddr
			server.Mallory.SwitchDNS(self.DNSAddr)
		}
		if self.Balance != server2.Balance {
			if self.Balance {
				server.Mallory.SetBalance(1)
			} else {
				server.Mallory.SetBalance(0)
			}
			server2.Balance = self.Balance
		}
		if self.ProxyMode != server2.ProxyMode {
			if server.Mallory.Status {
				server.Mallory.ProxyHandler.ProxyMode = self.ProxyMode
			}
			server2.ProxyMode = self.ProxyMode
		}
		if self.EnableAuth != server2.EnableAuth {
			if self.EnableAuth {
				server.Mallory.ProxyHandler.EnableAuth()
			} else {
				server.Mallory.ProxyHandler.DisableAuth()
			}
			server2.EnableAuth = self.EnableAuth
		}
		err := filedb2.DB.Set("settings", "server", server2)
		if err != nil {
			response.ErrorWithMessage(http.StatusServiceUnavailable, err.Error())
		} else {
			response.Success()
		}
	}
	return
}

func NewUpdateProxyServerRequest() *UpdateProxyServerRequest {
	return &UpdateProxyServerRequest{}
}

type InfoProxyServerRequest struct{}

func (self *InfoProxyServerRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	info, err := models.DefaultProxyServer()
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		data := garray.New(true)
		data.Append(g.Map{"key": "port", "name": "端口", "value": info.Port})
		data.Append(g.Map{"key": "status", "name": "状态", "value": server.Mallory.Status})
		data.Append(g.Map{"key": "dns_addr", "name": "DNS", "value": info.DNSAddr})
		data.Append(g.Map{"key": "balance", "name": "负载", "value": info.Balance})
		data.Append(g.Map{"key": "auto_start", "name": "启动", "value": info.AutoStart})
		data.Append(g.Map{"key": "proxy_mode", "name": "模式", "value": info.ProxyMode})
		data.Append(g.Map{"key": "enable_auth", "name": "认证", "value": info.EnableAuth})
		if server.Mallory.Error != nil && server.Mallory.Error != http.ErrServerClosed {
			data.Append(g.Map{"key": "error", "name": "错误", "value": server.Mallory.Error.Error()})
		}

		response.SuccessWithDetail(g.Map{
			"data":      data,
			"instances": server.Mallory.InstancesInfo(),
		})
	}
	return
}

func NewInfoProxyServerRequest() *InfoProxyServerRequest {
	return &InfoProxyServerRequest{}
}
