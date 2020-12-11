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
	server2, err := models.GetProxyServer()
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err)
	} else {
		data := g.Map{}
		if self.Name != "" && server2.Name != self.Name {
			server2.Name = self.Name
		}
		if self.Port != 0 && server2.Port != self.Port {
			server2.Port = self.Port
		}
		if self.Username != "" && server2.Username == self.Username {
			server2.Name = self.Username
		}
		if self.Password != "" && server2.Password == self.Password {
			server2.Password = self.Password
		}
		if self.Status != server2.Status {
			server2.Status = self.Status
			if self.Status {
				server.Mallory.SetBalance(1)
			} else {
				server.Mallory.SetBalance(0)
			}
		}
		if self.AutoProxy != server2.AutoProxy {
			data["auto_proxy"] = self.AutoProxy
		}
		if self.AllProxy != server2.AllProxy {
			data["all_proxy"] = self.AllProxy
			if server.Mallory.Status {
				server.Mallory.ProxyHandler.AllProxy = self.AllProxy
			}
		}
		err := filedb2.DB.Update(server2)
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
	info, err := models.GetProxyServer()
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		data := garray.New(true)
		data.Append(g.Map{"key": "port", "name": "端口", "value": info.Port})
		data.Append(g.Map{"key": "status", "name": "状态", "value": server.Mallory.Status})
		data.Append(g.Map{"key": "balance", "name": "负载", "value": info.Status})
		data.Append(g.Map{"key": "auto_start", "name": "启动", "value": info.AutoStart})
		data.Append(g.Map{"key": "all_proxy", "name": "模式", "value": info.AllProxy})
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
