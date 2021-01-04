package mallory

import (
	"context"
	"homeproxy/library/common"
	"net"
	"time"

	"github.com/gogf/gf/os/gcache"
	"github.com/gogf/gf/os/glog"
)

type DnsCache struct {
	DNS   *net.Resolver
	Cache *gcache.Cache
}

func (self *DnsCache) search(domain string) (string, error) {
	if found, err := self.Cache.Contains(domain); err != nil {
		return "", err
	} else if !found {
		if country, err := self.lookupHost(domain); err != nil {
			return "", err
		} else {
			return country, nil
		}
	} else {
		v, _ := self.Cache.Get(domain)
		return v.(string), nil
	}
}

func (self *DnsCache) lookupHost(domain string) (string, error) {
	addrs, err := self.DNS.LookupHost(context.Background(), domain)
	if err != nil {
		return "", err
	}
	country, err := common.LookupCountry(addrs[0])
	if err != nil {
		return "", err
	}
	self.Cache.Set(domain, country, time.Second * 600)
	return country, nil
}

func (self *DnsCache) IsChina(domain string) bool {
	country, err := self.search(domain)
	glog.Debugf("domain: %s country: %s", domain, country)
	if err != nil {
		return false
	} else {
		return country == "CN"
	}
}
