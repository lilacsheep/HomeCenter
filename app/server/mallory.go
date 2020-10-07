package server

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gutil"
	"homeproxy/app/models"
	"homeproxy/library/mallory"
	"net/http"
	"time"
)

var Mallory *MalloryManger

func init() {
	Mallory = &MalloryManger{}
}

type MalloryManger struct {
	HttpServer   *http.Server    `json:"-"`
	ProxyHandler *mallory.Server `json:"-"`
	Error        error
	Status       bool
}

func (self *MalloryManger) Init() error {
	// init
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
		self.ProxyHandler.AddUrlRole(p.Sub, p.Domain, p.Status)
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
	self.ProxyHandler.Balance = mallory.NewRandomBalance()

	// add ssh instance
	for _, instance := range instances {
		self.ProxyHandler.Balance.AddInstances(instance.Url(), instance.PrivateKey)
	}

	// set http server Handler
	self.HttpServer.Handler = self.ProxyHandler
	self.HttpServer.Addr = fmt.Sprintf(":%d", info.Port)
	return nil
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
				self.Error = err
				glog.Errorf("proxy error exit: %s", err.Error())
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
			self.ProxyHandler.Balance.ReleaseInstances()
			self.Status = false
			self.HttpServer = nil
			self.ProxyHandler = nil
		}
	} else {
		return errors.New("http proxy server already stop")
	}
	return nil
}
