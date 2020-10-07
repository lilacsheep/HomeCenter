package requests

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/guid"
	"homeproxy/app/models"
	"homeproxy/app/server"
	"net/http"
)

type CreateProxyInstanceRequest struct {
	Address    string `v:"address     @required"`
	Username   string `v:"username    @required|length:4,32#请输入用户名称|用户名称长度非法"`
	Password   string `v:"password    @required-without:private_key"`
	PrivateKey string `v:"private_key @required-without:password"`
}

func (self *CreateProxyInstanceRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	instance := models.ProxyInstance{}
	instance.ID = guid.S()
	instance.Address = self.Address
	instance.Username = self.Username
	instance.Password = self.Password
	instance.PrivateKey = self.PrivateKey
	instance.Status = true

	_, err := g.DB().Table(models.ProxyInstanceTable).Insert(&instance)
	if err != nil {
		glog.Errorf("add instance error: %s", err.Error())
		glog.Errorf("add instance error: %v", instance)
		response.ErrorWithMessage(http.StatusInternalServerError, err)
	} else {
		if server.Mallory.Status {
			server.Mallory.ProxyHandler.Balance.AddInstances(instance.Url(), instance.PrivateKey)
		}
		response.SuccessWithDetail(instance)
	}
	return
}

func NewCreateProxyInstanceRequest() *CreateProxyInstanceRequest {
	return &CreateProxyInstanceRequest{}
}

type QueryAllInstanceRequest struct{}

func (self *QueryAllInstanceRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	querySet := g.DB().Table(models.ProxyInstanceTable)
	var instances []models.ProxyInstance
	err := querySet.Structs(&instances)
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		response.SuccessWithDetail(instances)
	}
	return
}

func NewQueryAllInstanceRequest() *QueryAllInstanceRequest {
	return &QueryAllInstanceRequest{}
}
