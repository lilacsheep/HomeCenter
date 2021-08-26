package models

import (
	"context"
	"fmt"
	"homeproxy/library/common"
	"net"
	"strings"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gmeta"
)

type ProxyInstance struct {
	DefaultModel
	Protocol     int    `json:"protocol"`               // default 0(SSH)
	Address      string `json:"address"` //
	Username     string `json:"username"`
	Password     string `json:"password"`
	PrivateKey   string `json:"private_key"`
	Status       bool   `json:"status"`
	ForceCountry bool   `json:"force_country"`
	CountryCode  string `json:"country_code"`
	Country      string `json:"country"`
	Delay        int    `json:"delay"`
	gmeta.Meta  `orm:"table:instances"`
}

func (self *ProxyInstance) Url() string {
	return fmt.Sprintf("ssh://%s@%s", self.Username, self.Address)
}

func (self *ProxyInstance) IsChina() bool {
	return self.CountryCode == "CN"
}

func (self *ProxyInstance) RefreshCountry() {
	var (
		code string
	)
	if self.ForceCountry {
		return
	}
	v := strings.Split(self.Address, ":")
	glog.Debugf("start refresh instance: %s info", self.Address)
	if common.CheckIp(v[0]) {
		glog.Debugf("check ip address: %s", v[0])
		code, _ = common.LookupCountry(v[0])
	} else {
		serverInfo, _ := DefaultProxyServer()
		dns := &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: 10 * time.Second,
				}
				return d.DialContext(ctx, "udp", fmt.Sprintf("%s:53", serverInfo.DNSAddr))
			},
		}
		addrs, _ := dns.LookupHost(context.Background(), v[0])
		for _, addr := range addrs {
			glog.Debugf("domain: %s resolv addr: %s", v[0], addr)
			code, _ = common.LookupCountry(addr)
			break
		}
	}
	info := common.SearchCountryFromCode(code)
	if info != nil {
		self.CountryCode = code
		self.Country = info.CN
		g.DB().Model(&ProxyInstance{}).Data(g.Map{"country": info.CN, "country_code": code}).Where("id = ?", self.Id).Update()
	}
	glog.Debugf("refresh instance: %s info done, code: %s", v[0], code)
}

func UpdateProxyInstanceDelay(id int, delay int) error {
	_, err := g.DB().Model(&ProxyInstance{}).Data(g.Map{"delay": delay}).Where("id=?", id).Update()
	return err
}

func GetEnableProxyInstances() (instances []ProxyInstance, err error) {
	err = g.DB().Model(&ProxyInstance{}).Where("status", true).Structs(&instances)
	return 
}
