// Package mallory implements a simple http proxy support direct and GAE remote fetcher
package mallory

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/os/glog"
	"golang.org/x/net/publicsuffix"
	"net/http"
	"strings"
	"sync"
)

const (
	SmartSrv = iota
	NormalSrv
)

type AccessType bool

func (t AccessType) String() string {
	if t {
		return "PROXY"
	} else {
		return "DIRECT"
	}
}

type ProxyRole struct {
	ID         string `orm:"id"`
	InstanceID string `orm:"instance_id"`
	Status     bool   `orm:"status"`
	Sub        string `orm:"sub"`
	Domain     string `orm:"domain"`
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
	// username
	Username string
	// Password
	Password string
	// Auto Proxy
	AutoProxy bool
	// for serve http
	mutex sync.RWMutex
}

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

func (self *Server) AddUrlRole(sub, domain string, status bool) {
	if status {
		if v, found := self.BlockedHosts.Search(domain); found {
			if sub != "" {
				self.BlockedHosts.Set(domain, v.(*garray.StrArray).Append(sub))
			}
		} else {
			subList := garray.NewStrArray(true)
			if sub != "" {
				subList = subList.Append(sub)
			}
			self.BlockedHosts.Set(domain, subList)
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

func (self *Server) RemoveUrlRole(sub, domain string) {
	if v, found := self.BlockedHosts.Search(domain); found {
		if sub != "" {
			v.(*garray.StrArray).RemoveValue(sub)
		} else {
			self.BlockedHosts.Remove(domain)
		}
	}
}

func (self *Server) Blocked(host string) bool {
	blocked := false
	sub, domain := self.UrlSplit(host)

	if value, f := self.BlockedHosts.Search(domain); f {
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
//
func (self *Server) BasicAuth(w http.ResponseWriter, r *http.Request) {
	basicAuthPrefix := "Basic "

	// 获取 request header
	auth := r.Header.Get("Authorization")
	// 如果是 http basic auth
	if strings.HasPrefix(auth, basicAuthPrefix) {
		// 解码认证信息
		payload, err := base64.StdEncoding.DecodeString(
			auth[len(basicAuthPrefix):],
		)
		if err == nil {
			pair := bytes.SplitN(payload, []byte(":"), 2)
			if len(pair) == 2 && bytes.Equal(pair[0], []byte(self.Username)) &&
				bytes.Equal(pair[1], []byte(self.Password)) {
				// 执行被装饰的函数
				//self.ServeProxy(w, r)
				return
			}
		}
	}

	// 认证失败，提示 401 Unauthorized
	// Restricted 可以改成其他的值
	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
	// 401 状态码
	w.WriteHeader(http.StatusUnauthorized)
}

func (self *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if self.Black(r.URL.Host) {
		return
	}

	use := self.Blocked(r.URL.Host) && r.URL.Host != ""
	glog.Infof("%s", r.RemoteAddr)
	glog.Infof("[%s] %d %s %s %s", AccessType(use), self.Mode, r.Method, r.RequestURI, r.Proto)

	if use {
		SSH, err := self.Balance.DoBalance()
		if err != nil {
			glog.Errorf("get proxy connect error: %s", err.Error())
		} else {
			if r.Method == "CONNECT" && SSH.Status {
				SSH.Connect(w, r)
			} else if r.URL.IsAbs() && SSH.Status {
				r.RequestURI = ""
				RemoveHopHeaders(r.Header)
				SSH.ServeHTTP(w, r)
			} else {
				glog.Infof("%s is not a full URL path", r.RequestURI)
			}
		}
	} else {
		if r.Method == "CONNECT" {
			if err := self.Direct.Connect(w, r); err != nil {
				glog.Errorf("connect %s error %s", r.URL.Host, err)
				//if self.AutoProxy {
				//	glog.Errorf("visit %s reproxy error %s", r.URL.Host, err)
				//	SSH, err := self.Balance.DoBalance()
				//	if err == nil {
				//		SSH.Connect(w, r)
				//	} else {
				//		glog.Errorf("get proxy connect error: %s", err.Error())
				//	}
				//} else {
				//	glog.Errorf("connect %s error %s", r.URL.Host, err)
				//}
			}
		} else if r.URL.IsAbs() {
			r.RequestURI = ""
			RemoveHopHeaders(r.Header)
			if err := self.Direct.ServeHTTP(w, r); err != nil {
				glog.Errorf("visit %s error %s", r.URL.Host, err)
				//if self.AutoProxy {
				//	glog.Errorf("visit %s reproxy error %s", r.URL.Host, err)
				//	SSH, err := self.Balance.DoBalance()
				//	if err == nil {
				//		SSH.ServeHTTP(w, r)
				//	} else {
				//		glog.Errorf("get proxy connect error: %s", err.Error())
				//	}
				//} else {
				//	glog.Errorf("visit %s error %s", r.URL.Host, err)
				//}
			}
		} else {
			glog.Infof("%s is not a full URL path", r.RequestURI)
		}
	}
}
