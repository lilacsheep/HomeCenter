package requests

import (
	"homeproxy/app/models"
	"homeproxy/app/server"
	"homeproxy/library/common"
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

	_, err := g.DB().Model(&models.ProxyInstance{}).Save(&instance)
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		instance.RefreshCountry()
		if server.Mallory.Status {
			server.Mallory.AddInstances(instance.Url(), instance.Password, instance.PrivateKey, gconv.String(instance.Id))
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

	err = g.DB().Model(&models.ProxyInstance{}).Limit(offset, limit).Structs(&instances)
	if err != nil {
		if err == storm.ErrNotFound {
			response.SuccessWithDetail([]models.ProxyInstance{})
		} else {
			response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
		}
	} else {
		response.SuccessWithDetail(instances)
	}
	return
}

func NewQueryAllInstanceRequest() *QueryAllInstanceRequest {
	return &QueryAllInstanceRequest{}
}

type UpdateInstanceRequest struct {
	ID           int    `v:"id @required"`
	Address      string `v:"address     @required"`
	Username     string `v:"username    @required|length:4,32#请输入用户名称|用户名称长度非法"`
	Password     string `v:"password    @required-without:private_key"`
	PrivateKey   string `v:"private_key @required-without:password"`
	CountryCode  string `json:"country_code"`
	ForceCountry bool   `json:"force_country"`
}

func (self *UpdateInstanceRequest) change(instance *models.ProxyInstance) bool {
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
	if instance.CountryCode != self.CountryCode {
		data["country_code"] = self.CountryCode
	}
	if instance.ForceCountry != self.ForceCountry {
		data["force_country"] = self.ForceCountry
	}
	return len(data) != 0
}

func (self *UpdateInstanceRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	var (
		err      error
		instance = &models.ProxyInstance{}
	)
	err = g.DB().Model(&models.ProxyInstance{}).Where("id = ?", self.ID).Struct(instance)
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		if self.change(instance) {
			c := common.SearchCountryFromCode(self.CountryCode)
			instance.Address = self.Address
			instance.Username = self.Username
			instance.Password = self.Password
			instance.PrivateKey = self.PrivateKey
			instance.Status = true
			if self.ForceCountry {
				instance.Country = c.CN
				instance.CountryCode = self.CountryCode
				instance.ForceCountry = self.ForceCountry
				g.DB().Model(&models.ProxyInstance{}).Data(g.Map{
					"country": c.CN,
					"country_code": self.CountryCode,
					"force_country": true,
				}).Where("id = ?", self.ID).Update()
			} else {
				instance.ForceCountry = self.ForceCountry
				g.DB().Model(&models.ProxyInstance{}).Data(g.Map{"force_country": false}).Where("id = ?", self.ID).Update()
				instance.RefreshCountry()
			}

			if err != nil {
				response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
			} else {
				if server.Mallory.Status {
					server.Mallory.RemoveInstance(gconv.String(self.ID))
					server.Mallory.AddInstances(instance.Url(), instance.Password, instance.PrivateKey, gconv.String(instance.Id))
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
	ID int `v:"id @required"`
}

func (self *RemoveInstanceRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	var (
		err      error
		instance models.ProxyInstance
	)
	err = g.DB().Model(&models.ProxyInstance{}).Where("id = ?", self.ID).Struct(&instance)
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		_, err = g.DB().Model(&models.ProxyInstance{}).Where("id = ?", self.ID).Delete()
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
	_, err = g.DB().Model(&models.ProxyInstance{}).Data(g.Map{"status": false}).Where("id = ?", self.ID).Update()
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
	err = g.DB().Model(&models.ProxyInstance{}).Where("id = ?", self.ID).Struct(&instance)
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		_, err = g.DB().Model(&models.ProxyInstance{}).Data(g.Map{"status": true}).Where("id = ?", self.ID).Update()
		if err != nil {
			response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
		} else {
			if server.Mallory.Status {
				server.Mallory.AddInstances(instance.Url(), instance.Password, instance.PrivateKey, gconv.String(instance.Id))
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
