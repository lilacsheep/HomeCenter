package requests

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"golang.org/x/net/publicsuffix"
	"homeproxy/app/models"
	"homeproxy/app/server"
	"homeproxy/library/filedb"
	"homeproxy/library/mallory"
	"net/http"
	"strings"
)

type AddUrlRoleRequest struct {
	Url        string `v:"url     @required"`
	InstanceID string
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
	var data []mallory.ProxyRole
	c, _ := filedb.DB.Collection(mallory.ProxyRoleTable)
	if sub == "" {
		sub = "*"
	}
	if err := c.Search(g.Map{"status": self.Status, "sub": sub, "domain": domain}, &data); err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
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
			c.Insert(&role)
			if server.Mallory.Status {
				server.Mallory.ProxyHandler.AddUrlRole(role.Sub, role.Domain, self.Status, self.InstanceID)
			}
			response.Success()
		}
	}
	return
}

func NewAddUrlRoleRequest() *AddUrlRoleRequest {
	return &AddUrlRoleRequest{}
}

type RemoveUrlRoleRequest struct {
	ID string `json:"id"`
}

func (self *RemoveUrlRoleRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	var role mallory.ProxyRole
	c, _ := filedb.DB.Collection(mallory.ProxyRoleTable)
	if err := c.GetAndRemove(self.ID, &role); err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
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

type QueryAllRoleRequest struct{}

func (self *QueryAllRoleRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	roles := models.AllRoles()
	response.DataTable(roles, len(roles))
	return
}

func NewQueryAllRoleRequest() *QueryAllRoleRequest {
	return &QueryAllRoleRequest{}
}

type ChangeRoleInstanceRequest struct {
	ID         string `v:"id     @required"`
	InstanceID string
	Status     bool
}

func (self *ChangeRoleInstanceRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	role := mallory.ProxyRole{}
	err := filedb.DB.GetById(mallory.ProxyRoleTable, self.ID, &role)
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		if server.Mallory.Status {
			server.Mallory.ProxyHandler.RemoveUrlRole(role.Sub, role.Domain, role.Status)
		}
		role.InstanceID = self.InstanceID
		role.Status = self.Status
		err = filedb.DB.UpdateById(mallory.ProxyRoleTable, self.ID, role)
		if err != nil {
			response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
		} else {
			if server.Mallory.Status {
				server.Mallory.ProxyHandler.AddUrlRole(role.Sub, role.Domain, role.Status, role.InstanceID)
			}
			response.Success()
		}
	}

	return
}

func NewChangeRoleInstanceRequest() *ChangeRoleInstanceRequest {
	return &ChangeRoleInstanceRequest{}
}
