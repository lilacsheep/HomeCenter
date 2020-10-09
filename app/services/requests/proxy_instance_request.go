package requests

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"homeproxy/app/models"
	"homeproxy/app/server"
	"homeproxy/library/filedb"
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
	instance.Address = self.Address
	instance.Username = self.Username
	instance.Password = self.Password
	instance.PrivateKey = self.PrivateKey
	instance.Status = true

	c, err := models.DB.Collection(models.ProxyInstanceTable)
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err)
	} else {
		c.Insert(&instance)
	}
	if server.Mallory.Status {
		server.Mallory.ProxyHandler.Balance.AddInstances(instance.Url(), instance.PrivateKey)
	}
	response.SuccessWithDetail(instance)
	return
}

func NewCreateProxyInstanceRequest() *CreateProxyInstanceRequest {
	return &CreateProxyInstanceRequest{}
}

type QueryAllInstanceRequest struct{}

func (self *QueryAllInstanceRequest) Exec(r *ghttp.Request) (response MessageResponse) {

	var (
		c         *filedb.Collection
		err       error
		instances []models.ProxyInstance
	)
	if c, err = models.DB.Collection(models.ProxyInstanceTable); err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		if err = c.Search(g.Map{}, &instances); err != nil {
			response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
		} else {
			response.SuccessWithDetail(instances)
		}
	}
	return
}

func NewQueryAllInstanceRequest() *QueryAllInstanceRequest {
	return &QueryAllInstanceRequest{}
}
