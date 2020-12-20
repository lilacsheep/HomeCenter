package requests

import (
	"fmt"
	"homeproxy/app/server"
	"homeproxy/library/filedb2"
	"homeproxy/library/mallory"
	"net/http"
	"strings"

	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
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
	var data []mallory.ProxyRole
	if sub == "" {
		sub = "*"
	}
	query := filedb2.DB.Select(q.Eq("Status", self.Status), q.Eq("Sub", sub), q.Eq("Domain", domain))
	err := query.Find(&data)
	if err != nil && err != storm.ErrNotFound {
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
			err = filedb2.DB.Save(&role)
			if err != nil {
				response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
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
	var role mallory.ProxyRole

	if err := filedb2.DB.One("ID", self.ID, &role); err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		filedb2.DB.DeleteStruct(&role)
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
	var data []mallory.ProxyRole
	offset, limit := self.Next()
	var query storm.Query
	if self.Filter != "" {
		query = filedb2.DB.Select(q.Re("Domain", self.Filter))
	} else {
		query = filedb2.DB.Select()
	}
	c, _ := query.Count(&mallory.ProxyRole{})

	err := query.Skip(offset).Limit(limit).Find(&data)
	if err != nil {
		if err == storm.ErrNotFound {
			response.DataTable([]mallory.ProxyRole{}, 0)
		} else {
			response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
		}
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
	role := mallory.ProxyRole{}
	err := filedb2.DB.One("ID", self.ID, &role)
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		if server.Mallory.Status {
			server.Mallory.ProxyHandler.RemoveUrlRole(role.Sub, role.Domain, role.Status)
		}
		role.InstanceID = self.InstanceID
		role.Status = self.Status
		err = filedb2.DB.Update(&role)
		if err != nil {
			response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
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
