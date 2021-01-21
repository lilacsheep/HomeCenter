// Package mallory implements a simple http proxy support direct and GAE remote fetcher
package mallory

import (
	"fmt"
	"homeproxy/library/common"
	"net/http"
	"strings"
	"sync"

	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/os/glog"
	"golang.org/x/net/publicsuffix"
)

type AccessType bool

func (t AccessType) String() string {
	if t {
		return "PROXY"
	} else {
		return "DIRECT"
	}
}


type Server struct {
	// SmartSrv or NormalSrv
	Mode int
	// config file
	Port int
	// direct fetcher
	Direct *Direct
	// ssh fetcher, to connect remote proxy server
	Balance Balance
	// a cache
	BlockedHosts *gmap.TreeMap
	// black list
	BlackHosts *gmap.TreeMap
	// Proxy mode
	ProxyMode int // 1 全代理模式 2 规则代理模式 3 DNS 代理模式
	// for serve http
	mutex sync.RWMutex
	// ssh instances
	Instances *gmap.TreeMap
	// custom dns
	DNSCache *DnsCache
	// Authentication
	Authentication Authentication
}

// split url
func (self *Server) UrlSplit(url string) (string, string) {
	host := HostOnly(url)
	domain, _ := publicsuffix.EffectiveTLDPlusOne(host)
	subDomain := ""
	t := strings.Split(host, fmt.Sprintf(".%s", domain))
	subDomain = t[0]
	if subDomain == domain {
		return "", domain
	}
	return subDomain, domain
}

func (self *Server) AddUrlRole(sub, domain string, status bool, instances ...string) {
	instanceId := ""
	if len(instances) > 0 {
		instanceId = instances[0]
	}
	if status {
		if v, found := self.BlockedHosts.Search(domain); found {
			if sub != "" {
				v.(*gmap.StrStrMap).Set(sub, instanceId)
			}
		} else {
			subMapping := gmap.NewStrStrMap(true)
			if sub != "" {
				subMapping.Set(sub, instanceId)
			}
			self.BlockedHosts.Set(domain, subMapping)
		}
	} else {
		if v, found := self.BlackHosts.Search(domain); found {
			if sub != "" {
				self.BlackHosts.Set(domain, v.(*garray.StrArray).Append(sub))
			}
		} else {
			subList := garray.NewStrArray(true)
			if sub != "" {
				subList = subList.Append(sub)
			}
			self.BlackHosts.Set(domain, subList)
		}
	}
}

func (self *Server) RemoveUrlRole(sub, domain string, status bool) {
	if status {
		if v, found := self.BlockedHosts.Search(domain); found {
			if sub != "" {
				v.(*gmap.StrStrMap).Remove(sub)
			} else {
				self.BlockedHosts.Remove(domain)
			}
		}
	} else {
		if v, found := self.BlackHosts.Search(domain); found {
			if sub != "" {
				v.(*garray.StrArray).RemoveValue(sub)
			} else {
				self.BlackHosts.Remove(domain)
			}
		}
	}
}

func (self *Server) Blocked(host string) (bool, string) {
	blocked := false
	instance := ""
	sub, domain := self.UrlSplit(host)

	if value, f := self.BlockedHosts.Search(domain); f {
		value.(*gmap.StrStrMap).Iterator(func(k string, v string) bool {
			if k == sub || k == "*" {
				blocked = true
				instance = v
				return false
			}
			return true
		})
	}
	return blocked, instance
}

func (self *Server) Black(host string) bool {
	blocked := false
	sub, domain := self.UrlSplit(host)

	if value, f := self.BlackHosts.Search(domain); f {
		value.(*garray.StrArray).Iterator(func(k int, v string) bool {
			if sub == v || v == "*" {
				blocked = true
				return false
			}
			return true
		})
	}
	return blocked
}

// ServeHTTP proxy accepts requests with following two types:
//  - CONNECT
//    Generally, this method is used when the client want to connect server with HTTPS.
//    In fact, the client can do anything he want in this CONNECT way...
//    The request is something like:
//      CONNECT www.google.com:443 HTTP/1.1
//    Only has the host and port information, and the proxy should not do anything with
//    the underlying data. What the proxy can do is just exchange data between client and server.
//    After accepting this, the proxy should response
//      HTTP/1.1 200 OK
//    to the client if the connection to the remote server is established.
//    Then client and server start to exchange data...
//
//  - non-CONNECT, such as GET, POST, ...
//    In this case, the proxy should redo the method to the remote server.
//    All of these methods should have the absolute URL that contains the host information.
//    A GET request looks like:
//      GET weibo.com/justmao945/.... HTTP/1.1
//    which is different from the normal http request:
//      GET /justmao945/... HTTP/1.1
//    Because we can be sure that all of them are http request, we can only redo the request
//    to the remote server and copy the reponse to client.
func (self *Server) local(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		if err := self.Direct.Connect(w, r); err != nil {
			// add error visit log
			errorEvent(r.URL.Hostname(), err)
		}
	} else if r.URL.IsAbs() {
		r.RequestURI = ""
		RemoveHopHeaders(r.Header)
		if err := self.Direct.ServeHTTP(w, r); err != nil {
			// add error visit log
			errorEvent(r.URL.Hostname(), err)
		}
	} else {
		glog.Infof("%s is not a full URL path", r.RequestURI)
	}
}

func (self *Server) overseas(w http.ResponseWriter, r *http.Request, instanceIDs ...string) {
	var (
		instanceID string
		connect    *SSH
	)
	if len(instanceIDs) > 0 {
		instanceID = instanceIDs[0]
	}
	if v, ok := self.Instances.Search(instanceID); ok {
		connect = v.(*SSH)
	} else {
		instance, err := self.Balance.DoBalance(self.Instances)
		if err != nil {
			glog.Errorf("get proxy connect error: %s", err.Error())
		} else {
			connect = instance.(*SSH)
		}
	}
	if connect == nil {
		glog.Errorf("get proxy connect is nil")
	} else {
		if r.Method == http.MethodConnect && connect.Status {
			connect.Connect(w, r)
		} else if r.URL.IsAbs() && connect.Status {
			r.RequestURI = ""
			RemoveHopHeaders(r.Header)
			connect.ServeHTTP(w, r)
		} else {
			glog.Infof("%s is not a full URL path", r.RequestURI)
		}
	}
}

func (self *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if self.Black(r.URL.Host) {
		return
	}
	var (
		use        bool
		instanceId string
		t          bool
	)
	if self.Authentication != nil && !self.Authentication.Auth(w, r) {
		return
	}
	ip := common.IPAddress(r.URL.Hostname())
	if !ip.IsPublic() {
		self.local(w, r)
	} else {
		// dns 路由模式
		switch self.ProxyMode {
		case 4: // 规则+DNS智能代理
			t, instanceId = self.Blocked(r.URL.Host)
			use = t && r.URL.Host != "" && self.Instances.Size() != 0
			if use {
				self.overseas(w, r, instanceId)
			} else {
				if v, _ := self.DNSCache.IsLocal(r.URL.Hostname()); v {
					self.local(w, r)
				} else {
					if self.DNSCache.IsChina(r.URL.Hostname()) {
						self.local(w, r)
					} else {
						self.overseas(w, r, instanceId)
					}
				}
			}
		case 3: // DNS智能代理
			if v, _ := self.DNSCache.IsLocal(r.URL.Hostname()); v {
				self.local(w, r)
			} else {
				if self.DNSCache.IsChina(r.URL.Hostname()) {
					self.local(w, r)
				} else {
					self.overseas(w, r, instanceId)
				}
			}
		case 2: // 规则代理
			t, instanceId = self.Blocked(r.URL.Host)
			use = t && r.URL.Host != "" && self.Instances.Size() != 0
			if use {
				self.overseas(w, r, instanceId)
			} else {
				self.local(w, r)
			}
		case 1: // 只代理
			self.overseas(w, r, instanceId)
		default: // 只走本地代理
			self.local(w, r)
		}
	}
}


func (self *Server) EnableAuth() {
	self.Authentication = BasicAuth()
}

func (self *Server) DisableAuth() {
	self.Authentication = nil
}