package mallory

import (
	"context"
	"github.com/gogf/gf/os/glog"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

type SSH struct {
	UUID string
	// global config file
	Cfg struct {
		RemoteServer string // e.g.  ssh://user:passwd@192.168.1.1:1122
		PrivateKey   string
	}
	// connect URL
	URL *url.URL
	// SSH client
	Client *ssh.Client
	// SSH client config
	CliCfg *ssh.ClientConfig
	// direct fetcher
	Direct *Direct
	// only re-dial once
	sf     Group
	l      sync.RWMutex
	Status bool
	Cancel context.CancelFunc
}

func (self *SSH) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	self.Direct.ServeHTTP(w, r)
}

func (self *SSH) Connect(w http.ResponseWriter, r *http.Request) {
	self.Direct.Connect(w, r)
}

func (self *SSH) SSHDail(network, addr string) (c net.Conn, err error) {
	self.l.RLock()
	cli := self.Client
	self.l.RUnlock()

	c, err = cli.Dial(network, addr)
	if err == nil {
		return
	}

	glog.Infof("dial %s failed: %s, reconnecting ssh server %s...", addr, err, self.URL.Host)

	clif, err := self.sf.Do(network+addr, func() (interface{}, error) {
		return ssh.Dial("tcp", self.URL.Host, self.CliCfg)
	})
	if err != nil {
		glog.Errorf("connect ssh server %s failed: %s", self.URL.Host, err)
		return
	}
	cli = clif.(*ssh.Client)

	self.l.Lock()
	self.Client = cli
	self.l.Unlock()

	return cli.Dial(network, addr)
}

func (self *SSH) Renew() {
	for {
		var err error
		self.Client, err = ssh.Dial("tcp", self.URL.Host, self.CliCfg)
		if err != nil {
			glog.Errorf("ssh proxy connect err: %s", err.Error())
			time.Sleep(5 * time.Second)
		} else {
			self.Direct = &Direct{
				Tr: &http.Transport{Dial: self.SSHDail},
			}
			self.Status = true
			glog.Info("ssh proxy connect success")
			break
		}
	}
}

func (self *SSH) Stop() error {
	self.Cancel()
	return self.Client.Close()
}
