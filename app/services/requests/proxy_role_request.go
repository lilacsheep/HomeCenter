package requests

import (
	"fmt"
	"homeproxy/app/models"
	"homeproxy/app/server"
	"homeproxy/library/mallory"
	"net/http"
	"strings"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
	"golang.org/x/net/publicsuffix"
)

type AddUrlRoleRequest struct {
	Url        string `v:"url     @required"`
	InstanceID int
	Status     bool
}

func (self *AddUrlRoleRequest) UrlSplit(url string) (string, string) {
	host := mallory.HostOnly(url)
	domain, _ := publicsuffix.EffectiveTLDPlusOne(host)
	subDomain := ""
	t := strings.Split(host, fmt.Sprintf(".%s", domain))
	subDomain = t[0]
	if subDomain == domain {
		return "", domain
	}
	return subDomain, domain
}

func (self *AddUrlRoleRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	sub, domain := self.UrlSplit(self.Url)
	var (
		data  []models.ProxyRole
		query = g.DB().Model(&models.ProxyRole{})
	)
	if sub == "" {
		sub = "*"
	}
	err := query.Clone().Where("`status` = ? AND `sub` = ? AND `domain` = ? ", self.Status, sub, domain).Structs(&data)
	if err != nil {
		response.SystemError(err)
	} else {
		if len(data) > 0 {
			response.ErrorWithMessage(http.StatusInternalServerError, "规则已经存在")
		} else {
			role := mallory.ProxyRole{
				InstanceID: self.InstanceID,
				Status:     self.Status,
				Sub:        sub,
				Domain:     domain,
			}
			_, err = query.Save(&role)
			if err != nil {
				response.SystemError(err)
			} else {
				if server.Mallory.Status {
					server.Mallory.ProxyHandler.AddUrlRole(role.Sub, role.Domain, self.Status, gconv.String(self.InstanceID))
				}
				response.Success()
			}
		}
	}
	return
}

func NewAddUrlRoleRequest() *AddUrlRoleRequest {
	return &AddUrlRoleRequest{}
}

type RemoveUrlRoleRequest struct {
	ID int `json:"id"`
}

func (self *RemoveUrlRoleRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	var (
		role  mallory.ProxyRole
		query = g.DB().Model(&mallory.ProxyRole{}).Where("`id` = ?", self.ID)
	)
	if err := query.Clone().Struct(&role); err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		query.Delete()
		if server.Mallory.Status {
			server.Mallory.ProxyHandler.RemoveUrlRole(role.Sub, role.Domain, role.Status)
		}
		response.Success()
	}
	return
}

func NewRemoveUrlRoleRequest() *RemoveUrlRoleRequest {
	return &RemoveUrlRoleRequest{}
}

type QueryAllRoleRequest struct {
	*Pagination
	Filter string `json:"filter"`
}

func (self *QueryAllRoleRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	var (
		data  []mallory.ProxyRole
		query = g.DB().Model(&mallory.ProxyRole{})
	)
	if self.Filter != "" {
		query = query.Where("`domain` LIKE %%%s%%", self.Filter)
	}
	c, _ := query.Count(&mallory.ProxyRole{})
	err := query.Limit(self.OffsetLimit()...).Structs(&data)
	if err != nil {
		response.SystemError(err)
	} else {
		response.DataTable(data, c)
	}
	return
}

func NewQueryAllRoleRequest() *QueryAllRoleRequest {
	return &QueryAllRoleRequest{}
}

type ChangeRoleInstanceRequest struct {
	ID         int `v:"id     @required"`
	InstanceID int `json:"instance_id"`
	Status     bool
}

func (self *ChangeRoleInstanceRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	var (
		query = g.DB().Model(&mallory.ProxyRole{}).Where("`id` = ?", self.ID)
		role  = mallory.ProxyRole{}
	)
	err := query.Clone().Struct(&role)
	if err != nil {
		response.SystemError(err)
	} else {
		if server.Mallory.Status {
			server.Mallory.ProxyHandler.RemoveUrlRole(role.Sub, role.Domain, role.Status)
		}
		_, err = query.Clone().Data(g.Map{"instance_id": self.InstanceID, "status": self.Status}).Update()
		if err != nil {
			response.SystemError(err)
		} else {
			if server.Mallory.Status {
				server.Mallory.ProxyHandler.AddUrlRole(role.Sub, role.Domain, role.Status, gconv.String(role.InstanceID))
			}
			response.Success()
		}
	}
	return
}

func NewChangeRoleInstanceRequest() *ChangeRoleInstanceRequest {
	return &ChangeRoleInstanceRequest{}
}
