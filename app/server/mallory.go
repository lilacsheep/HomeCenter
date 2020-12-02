package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/guid"
	"github.com/gogf/gf/util/gutil"
	"golang.org/x/crypto/ssh"
	"homeproxy/app/models"
	"homeproxy/library/mallory"
	"net"
	"net/http"
	"net/url"
	"os/user"
	"time"
)

var Mallory *MalloryManger

func init() {
	Mallory = &MalloryManger{}
	server, err := models.GetProxyServer()
	if err != nil {
		panic(err)
	} else {
		if server.AutoStart {
			glog.Info("proxy auto start, please wait...")
			Mallory.Start()
			glog.Info("[proxy] start ok...")
		}
	}
}

type MalloryManger struct {
	HttpServer   *http.Server    `json:"-"`
	ProxyHandler *mallory.Server `json:"-"`
	BalanceType  int
	Instances    *gmap.TreeMap
	Error        error
	Status       bool
}

func (self *MalloryManger) Init() error {
	// init
	self.Instances = gmap.NewTreeMap(gutil.ComparatorString, true)
	self.Status = false
	self.HttpServer = nil
	self.ProxyHandler = nil

	// init http server
	self.HttpServer = &http.Server{}
	self.HttpServer.SetKeepAlivesEnabled(false)

	// init Handler
	self.ProxyHandler = &mallory.Server{
		Direct:       mallory.NewDirect(30 * time.Second),
		BlockedHosts: gmap.NewTreeMap(gutil.ComparatorString, true),
		BlackHosts:   gmap.NewTreeMap(gutil.ComparatorString, true),
	}
	// init url role
	for _, p := range models.AllRoles() {
		self.ProxyHandler.AddUrlRole(p.Sub, p.Domain, p.Status, p.InstanceID)
	}

	// get Proxy Server
	info, err := models.GetProxyServer()
	if err != nil {
		return err
	}

	// get enable instances
	instances, err := models.GetEnableProxyInstances()
	if err != nil {
		return err
	}
	// set handler setting
	self.ProxyHandler.Port = info.Port
	self.ProxyHandler.Username = info.Username
	self.ProxyHandler.Password = info.Password
	self.ProxyHandler.AutoProxy = info.AutoProxy
	self.ProxyHandler.AllProxy = info.AllProxy

	// set proxy Balance
	switch self.BalanceType {
	case 0:
		self.ProxyHandler.Balance = mallory.NewRandomBalance()
	case 1:
		self.ProxyHandler.Balance = mallory.NewRoundRobinBalance()
	}

	// add ssh instance
	for _, instance := range instances {
		self.AddInstances(instance.Url(), instance.Password, instance.PrivateKey, instance.ID)
	}
	self.ProxyHandler.Instances = self.Instances

	// set http server Handler
	self.HttpServer.Handler = self.ProxyHandler
	self.HttpServer.Addr = fmt.Sprintf(":%d", info.Port)

	return nil
}

func (self *MalloryManger) SetBalance(balanceType int) {
	if balanceType != self.BalanceType {
		self.BalanceType = balanceType
		switch balanceType {
		case 0:
			self.ProxyHandler.Balance = mallory.NewRandomBalance()
		case 1:
			self.ProxyHandler.Balance = mallory.NewRoundRobinBalance()
		}
	}
}

func (self *MalloryManger) AddInstances(remoteUrl, Password, PrivateKey string, id ...string) {
	var (
		uuid string
		err  error
	)
	if len(id) > 0 {
		uuid = id[0]
	} else {
		uuid = guid.S()
	}

	Instance := &mallory.SSH{
		UUID: uuid,
		Cfg: struct {
			RemoteServer string
			PrivateKey   string
		}{RemoteServer: remoteUrl, PrivateKey: PrivateKey},
		CliCfg: &ssh.ClientConfig{},
	}
	Instance.URL, err = url.Parse(Instance.Cfg.RemoteServer)
	if err != nil {
		glog.Errorf("Error parsing link, %s", err.Error())
		return
	}
	if Instance.URL.User != nil {
		Instance.CliCfg.User = Instance.URL.User.Username()
	} else {
		u, err := user.Current()
		if err != nil {
			glog.Errorf("Error parsing link, %s", err.Error())
			return
		}
		// u.Name is the full name, should not be used
		Instance.CliCfg.User = u.Username
	}

	// host key break
	Instance.CliCfg.HostKeyCallback = func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}
	// ssh auth
	if Instance.Cfg.PrivateKey == "" {
		if Instance.URL.User == nil || Password == "" {
			glog.Errorf("%s not found user or password", Instance.URL)
			return
		}
		Instance.CliCfg.Auth = append(Instance.CliCfg.Auth, ssh.Password(Password))
	} else {
		signer, err := ssh.ParsePrivateKey([]byte(PrivateKey))
		if err != nil {
			glog.Errorf("ParsePrivateKey %s failed:%s", Instance.Cfg.PrivateKey, err)
			return
		}
		Instance.CliCfg.Auth = append(Instance.CliCfg.Auth, ssh.PublicKeys(signer))
	}

	// init Client , first time to dial to remote server, make sure it is available
	Instance.Client, err = ssh.Dial("tcp", Instance.URL.Host, Instance.CliCfg)
	if err != nil {
		glog.Errorf("connect err: %s, %s", Instance.URL.Host, err)
		return
	}
	Instance.Direct = &mallory.Direct{
		Tr: &http.Transport{Dial: Instance.SSHDail},
	}
	Instance.Status = true
	// set instance keepalive
	ctx, Cancel := context.WithCancel(context.Background())
	Instance.Cancel = Cancel
	go self.InstanceKeepAlive(Instance, ctx)

	// add instance in proxy
	self.Instances.Set(uuid, Instance)
}

func (self *MalloryManger) InstanceKeepAlive(instance *mallory.SSH, ctx context.Context) {
	for {
		select {
		case _ = <-ctx.Done():
			return
		default:
			t1 := time.Now()
			_, _, err := instance.Client.SendRequest("keepalive", true, nil)
			if err != nil {
				glog.Errorf("ssh proxy connect err: %s ", err.Error())
				instance.Status = false
				instance.Renew()
			} else {
				_ = models.UpdateProxyInstanceDelay(instance.UUID, gconv.Int(time.Now().Sub(t1).Milliseconds()))
			}
		}
		time.Sleep(2 * time.Second)
	}
}

func (self *MalloryManger) RemoveInstance(uuid string) {
	value := self.Instances.Remove(uuid)
	if instance, ok := value.(*mallory.SSH); ok {
		if err := instance.Stop(); err != nil {
			glog.Errorf("stop instance error: %s", err.Error())
		}
	}
}

func (self *MalloryManger) ReleaseInstances() {
	for _, key := range self.Instances.Keys() {
		value := self.Instances.Remove(key)
		instance := value.(*mallory.SSH)
		if err := instance.Stop(); err != nil {
			glog.Errorf("stop instance error: %s", err.Error())
		}
	}
}

func (self *MalloryManger) InstancesInfo() []interface{} {
	if !self.Status {
		return []interface{}{}
	}
	var data = garray.New(true)
	self.Instances.Iterator(func(key, value interface{}) bool {
		if s, ok := value.(*mallory.SSH); ok {
			data.Append(g.Map{
				"id":      s.UUID,
				"address": s.URL.Hostname(),
				"status":  s.Status,
			})
		}
		return true
	})
	return data.Slice()
}

func (self *MalloryManger) Start() error {
	if !self.Status {
		if err := Mallory.Init(); err != nil {
			return err
		}
		self.Status = true

		go func() {
			err := self.HttpServer.ListenAndServe()
			if err != nil {
				self.Status = false
				if err != http.ErrServerClosed {
					self.Error = err
					glog.Errorf("proxy error exit: %s", err.Error())
				}
			}
		}()

	} else {
		return errors.New("proxy http server already started")
	}
	return nil
}

func (self *MalloryManger) Stop() error {
	if self.Status {
		err := self.HttpServer.Close()
		if err != nil {
			glog.Errorf("http proxy server stop err: %s", err.Error())
			return err
		} else {
			self.ReleaseInstances()
			self.Status = false
			self.HttpServer = nil
			self.ProxyHandler = nil
		}
	} else {
		return errors.New("http proxy server already stop")
	}
	return nil
}
