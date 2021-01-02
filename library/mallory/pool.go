package mallory

import (
	"context"
	"errors"
	"homeproxy/library/common"
	"math/rand"
	"net"
	"time"

	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/os/gcache"
)

var (
	ErrNoInstances = errors.New("no instance")
)

type Balance interface {
	DoBalance(instances *gmap.TreeMap, domain ...string) (interface{}, error)
}

type RandomBalance struct{}

func (self *RandomBalance) DoBalance(instances *gmap.TreeMap, domain ...string) (interface{}, error) {
	if instances.IsEmpty() {
		return nil, ErrNoInstances
	} else if instances.Size() == 1 {
		return instances.Values()[0], nil
	} else {
		index := rand.Intn(instances.Size())
		return instances.Values()[index], nil
	}
}

type RoundRobinBalance struct {
	curIndex int
}

func (self *RoundRobinBalance) DoBalance(instances *gmap.TreeMap, domain ...string) (interface{}, error) {
	if instances.IsEmpty() {
		return nil, ErrNoInstances
	} else if instances.Size() == 1 {
		return instances.Values()[0], nil
	} else {
		if self.curIndex >= instances.Size() {
			self.curIndex = 0
		}
		self.curIndex++
		return instances.Values()[self.curIndex], nil
	}
}

// dns 路由负载
type DNSBalance struct {
	dns   *net.Resolver
	cache *gcache.Cache
}

func (self *DNSBalance) DoBalance(instances *gmap.TreeMap, domains ...string) (interface{}, error) {
	var (
		code string
	)
	if instances.IsEmpty() {
		return nil, ErrNoInstances
	} else if instances.Size() == 1 {
		return instances.Values()[0], nil
	} else {
		domain := domains[0]
		if v, _ := self.cache.Get(domain); v == nil {
			addrs, _ := self.dns.LookupHost(context.Background(), domain)
			code, _ = common.LookupCountry(addrs[0])
			self.cache.Set(domain, code, time.Minute * 10)
		} else {
			code = v.(string)
		}
		instances.Iterator(func(key, value interface{}) bool {
			return true
		})

		return nil, nil
	}
}

func NewDNSBalance(dns *net.Resolver) *DNSBalance {
	return &DNSBalance{
		dns: dns,
		cache: gcache.New(600),
	}
}

func NewRandomBalance() *RandomBalance {
	return &RandomBalance{}
}

func NewRoundRobinBalance() *RoundRobinBalance {
	return &RoundRobinBalance{
		curIndex: 0,
	}
}
