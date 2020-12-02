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

	c, err := filedb.DB.Collection(models.ProxyInstanceTable)
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err)
	} else {
		id, err := c.Insert(&instance)
		if err != nil && err == filedb.ErrUnique {
			response.ErrorWithMessage(http.StatusInternalServerError, "该记录已存在")
		} else {
			if server.Mallory.Status {
				server.Mallory.AddInstances(instance.Url(), instance.Password, instance.PrivateKey, id)
			}
			response.SuccessWithDetail(instance)
		}
	}
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
	if c, err = filedb.DB.Collection(models.ProxyInstanceTable); err != nil {
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

type UpdateInstanceRequest struct {
	ID         string `v:"id @required"`
	Address    string `v:"address     @required"`
	Username   string `v:"username    @required|length:4,32#请输入用户名称|用户名称长度非法"`
	Password   string `v:"password    @required-without:private_key"`
	PrivateKey string `v:"private_key @required-without:password"`
}

func (self *UpdateInstanceRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	var (
		err      error
		c        *filedb.Collection
		instance models.ProxyInstance
	)
	if c, err = filedb.DB.Collection(models.ProxyInstanceTable); err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		if err = c.GetById(self.ID, &instance); err != nil {
			response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
		} else {
			data := g.Map{}
			if instance.Address != self.Address {
				data["address"] = self.Address
				instance.Address = self.Address
			}
			if instance.Username != self.Username {
				data["username"] = self.Username
				instance.Username = self.Username
			}
			if instance.Password != self.Password {
				data["password"] = self.Password
				instance.Password = self.Password
			}
			if instance.PrivateKey != self.PrivateKey {
				data["private_key"] = self.PrivateKey
				instance.PrivateKey = self.PrivateKey
			}
			if len(data) == 0 {
				response.ErrorWithMessage(http.StatusInternalServerError, "没有改变")
			} else {
				err = c.UpdateById(self.ID, data)
				if err != nil {
					response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
				} else {
					if server.Mallory.Status {
						server.Mallory.RemoveInstance(self.ID)
						server.Mallory.AddInstances(instance.Url(), instance.Password, instance.PrivateKey, self.ID)
					}
					response.Success()
				}
			}
		}
	}
	return
}

func NewUpdateInstanceRequest() *UpdateInstanceRequest {
	return &UpdateInstanceRequest{}
}

type RemoveInstanceRequest struct {
	ID string `v:"id @required"`
}

func (self *RemoveInstanceRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	var (
		err error
		c   *filedb.Collection
	)
	if c, err = filedb.DB.Collection(models.ProxyInstanceTable); err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		c.RemoveById(self.ID)
		if server.Mallory.Status {
			server.Mallory.RemoveInstance(self.ID)
		}
		response.Success()
	}
	return
}

func NewRemoveInstanceRequest() *RemoveInstanceRequest {
	return &RemoveInstanceRequest{}
}

type RemoveInstanceFromPoolRequest struct {
	ID string `v:"id @required"`
}

func (self *RemoveInstanceFromPoolRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	var (
		err error
		c   *filedb.Collection
	)
	if c, err = filedb.DB.Collection(models.ProxyInstanceTable); err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		if err := c.UpdateById(self.ID, g.Map{"status": false}); err != nil {
			response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
		}
		server.Mallory.RemoveInstance(self.ID)
		response.Success()
	}
	return
}

func NewRemoveInstanceFromPoolRequest() *RemoveInstanceFromPoolRequest {
	return &RemoveInstanceFromPoolRequest{}
}

type AddInstanceIntoPoolRequest struct {
	ID string `v:"id @required"`
}

func (self *AddInstanceIntoPoolRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	var (
		err      error
		c        *filedb.Collection
		instance models.ProxyInstance
	)
	if c, err = filedb.DB.Collection(models.ProxyInstanceTable); err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		if err = c.GetById(self.ID, &instance); err != nil {
			response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
		} else {
			c.UpdateById(self.ID, g.Map{"status": true})
			if server.Mallory.Status {
				server.Mallory.AddInstances(instance.Url(), instance.Password, instance.PrivateKey, self.ID)
				response.Success()
			} else {
				response.ErrorWithMessage(http.StatusInternalServerError, "代理服务没有启动，但是会在代理启动时启动...")
			}
		}
	}
	return
}

func NewAddInstanceIntoPoolRequest() *AddInstanceIntoPoolRequest {
	return &AddInstanceIntoPoolRequest{}
}
