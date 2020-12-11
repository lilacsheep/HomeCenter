package requests

import (
	"homeproxy/app/models"
	"homeproxy/app/server"
	"homeproxy/library/filedb2"
	"net/http"

	"github.com/asdine/storm/v3"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
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

	err := filedb2.DB.Save(&instance)
	if err != nil {
		if err == storm.ErrAlreadyExists {
			response.ErrorWithMessage(http.StatusInternalServerError, "该记录已存在")
		} else {
			response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
		}
	} else {
		if server.Mallory.Status {
			server.Mallory.AddInstances(instance.Url(), instance.Password, instance.PrivateKey, gconv.String(instance.ID))
		}
		response.SuccessWithDetail(instance)
	}
	return
}

func NewCreateProxyInstanceRequest() *CreateProxyInstanceRequest {
	return &CreateProxyInstanceRequest{}
}

type QueryAllInstanceRequest struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

func (self *QueryAllInstanceRequest) Pagination() (int, int) {
	var (
		limit  = 10
		page   = 1
		offset = 0
	)
	if self.Limit != 0 {
		limit = self.Limit
	}
	if self.Page != 0 {
		page = self.Page
	}
	offset = (page - 1) * limit
	return offset, limit
}

func (self *QueryAllInstanceRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	var (
		err       error
		instances []models.ProxyInstance
	)
	offset, limit := self.Pagination()
	err = filedb2.DB.Select().Limit(limit).Skip(offset).Find(&instances)
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

type UpdateInstanceRequest struct {
	ID         int    `v:"id @required"`
	Address    string `v:"address     @required"`
	Username   string `v:"username    @required|length:4,32#请输入用户名称|用户名称长度非法"`
	Password   string `v:"password    @required-without:private_key"`
	PrivateKey string `v:"private_key @required-without:password"`
}

func (self *UpdateInstanceRequest) change(instance models.ProxyInstance) bool {
	data := g.Map{}
	if instance.Address != self.Address {
		data["address"] = self.Address
	}
	if instance.Username != self.Username {
		data["username"] = self.Username
	}
	if instance.Password != self.Password {
		data["password"] = self.Password
	}
	if instance.PrivateKey != self.PrivateKey {
		data["private_key"] = self.PrivateKey
	}
	return len(data) != 0
}

func (self *UpdateInstanceRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	var (
		err      error
		instance models.ProxyInstance
	)
	err = filedb2.DB.One("ID", self.ID, &instance)
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		if self.change(instance) {
			err = filedb2.DB.Update(&models.ProxyInstance{
				ID: self.ID,
				Address: self.Address,
				Username: self.Username,
				Password: self.Password,
			})
			if err != nil {
				response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
			} else {
				if server.Mallory.Status {
					server.Mallory.RemoveInstance(gconv.String(self.ID))
					server.Mallory.AddInstances(instance.Url(), instance.Password, instance.PrivateKey, gconv.String(instance.ID))
				}
				response.Success()
			}
		} else {
			response.ErrorWithMessage(http.StatusInternalServerError, "没有改变")
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
		instance models.ProxyInstance
	)
	err = filedb2.DB.One("ID", self.ID, &instance)
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		err = filedb2.DB.DeleteStruct(&instance)
		if err != nil {
			response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
		} else {
			response.Success()
		}
	}
	return
}

func NewRemoveInstanceRequest() *RemoveInstanceRequest {
	return &RemoveInstanceRequest{}
}

type RemoveInstanceFromPoolRequest struct {
	ID int `v:"id @required"`
}

func (self *RemoveInstanceFromPoolRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	var (
		err error
	)
	err = filedb2.DB.UpdateField(models.ProxyInstance{ID: self.ID}, "status", false)
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		if server.Mallory.Status {
			server.Mallory.RemoveInstance(gconv.String(self.ID))
		}
		response.Success()
	}
	return
}

func NewRemoveInstanceFromPoolRequest() *RemoveInstanceFromPoolRequest {
	return &RemoveInstanceFromPoolRequest{}
}

type AddInstanceIntoPoolRequest struct {
	ID int `v:"id @required"`
}

func (self *AddInstanceIntoPoolRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	var (
		err      error
		instance models.ProxyInstance
	)
	err = filedb2.DB.One("ID", self.ID, &instance)
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		err = filedb2.DB.UpdateField(&instance, "Status", true)
		if err != nil {
			response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
		} else {
			if server.Mallory.Status {
				server.Mallory.AddInstances(instance.Url(), instance.Password, instance.PrivateKey, gconv.String(instance.ID))
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
