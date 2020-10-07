package requests

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/guid"
	"golang.org/x/net/publicsuffix"
	"homeproxy/app/models"
	"homeproxy/app/server"
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
	querySet := g.DB().Table(models.ProxyRoleTable)
	exist := false

	c, _ := querySet.Count("status = ? and (domain = ? and sub = '*')", self.Status, domain)
	if c != 0 {
		exist = true
	} else {
		if sub == "" {
			sub = "*"
		}
		n, _ := querySet.Count("status = ? and (domain = ? and sub = ?)", self.Status, domain, sub)
		if n > 0 {
			exist = true
		}
	}
	if exist {
		response.ErrorWithMessage(http.StatusInternalServerError, "规则已经存在")
	} else {
		role := mallory.ProxyRole{
			ID:         guid.S(),
			InstanceID: self.InstanceID,
			Status:     self.Status,
			Sub:        sub,
			Domain:     domain,
		}
		_, err := g.DB().Table(models.ProxyRoleTable).Insert(&role)
		if err != nil {
			response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
		} else {
			if server.Mallory.Status {
				server.Mallory.ProxyHandler.AddUrlRole(role.Sub, role.Domain, self.Status)
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
	ID     string
	Sub    string
	Domain string `v:"domain     @required"`
}

func (self *RemoveUrlRoleRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	_, err := g.DB().Table(models.ProxyRoleTable).Delete("sub = ? and domain = ?", self.Sub, self.Domain)
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		if server.Mallory.Status {
			server.Mallory.ProxyHandler.RemoveUrlRole(self.Sub, self.Domain)
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
