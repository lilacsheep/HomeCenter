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
	ForceCountry *gcache.Cache
}

func (self *DnsCache) search(domain string) (string, error) {
	_, host := UrlSplit(domain)
	if v, _ := self.ForceCountry.Contains(host); v {
		c, _ := self.ForceCountry.Get(host)
		return c.(string), nil
	}
	if found, err := self.Cache.Contains(domain); err != nil {
		return "", err
	} else if !found {
		if country, err := self.lookupCountry(domain); err != nil {
			return "", err
		} else {
			return country, nil
		}
	} else {
		v, _ := self.Cache.Get(domain)
		return v.(string), nil
	}
}

func (self *DnsCache) LookupHost(domain string) (string, error) {
	addrs, err := self.DNS.LookupHost(context.Background(), domain)
	if err != nil {
		return "", err
	}
	return addrs[0], nil
}

func (self *DnsCache) IsLocal(domain string) (bool, error) {
	addr, err := self.LookupHost(domain)
	if err != nil {
		return false, err
	}
	ip := common.IPAddress(addr)
	if ip.Verify() {
		return !ip.IsPublic(), nil
	}
	return false, nil
}

func (self *DnsCache) lookupCountry(domain string) (string, error) {
	addr, err:= self.LookupHost(domain)
	if err != nil {
		return "", err
	}
	country, err := common.LookupCountry(addr)
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

func DefaultForceCountry() *gcache.Cache {
	cache := gcache.New()
	cache.Set("baidu.com", "CN", 0)
	return cache
}