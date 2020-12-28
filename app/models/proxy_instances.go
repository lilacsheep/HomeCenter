package models

import (
	"context"
	"fmt"
	"homeproxy/library/common"
	"homeproxy/library/filedb2"
	"net"
	"strings"
	"time"

	"github.com/asdine/storm/v3/q"
	"github.com/gogf/gf/os/glog"
)

type ProxyInstance struct {
	ID          int    `json:"id" storm:"id,increment"`
	Protocol    int    `json:"protocol"`               // default 0(SSH)
	Address     string `json:"address" storm:"unique"` //
	Username    string `json:"username"`
	Password    string `json:"password"`
	PrivateKey  string `json:"private_key"`
	Status      bool   `json:"status"`
	CountryCode string `json:"country_code"`
	Country     string `json:"country"`
	Delay       int    `json:"delay"`
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
	v := strings.Split(self.Address, ":")
	glog.Debugf("start refresh instance: %s info", self.Address)
	if common.CheckIp(v[0]) {
		glog.Debugf("check ip address: %s", v[0])
		code, _ = common.LookupCountry(v[0])
	} else {
		serverInfo, _ := GetProxyServer()
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
		filedb2.DB.Update(self)
	}
	glog.Debugf("refresh instance: %s info done, code: %s", v[0], code)
}

func UpdateProxyInstanceDelay(id int, delay int) error {
	return filedb2.DB.UpdateField(&ProxyInstance{ID: id}, "Delay", delay)
}

func GetEnableProxyInstances() (instances []ProxyInstance, err error) {
	err = filedb2.DB.Select(q.Eq("Status", true)).Find(&instances)
	return
}
