package mallory

import (
	"errors"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/guid"
	"github.com/gogf/gf/util/gutil"
	"golang.org/x/crypto/ssh"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os/user"
)

type Balance interface {
	DoBalance() (*SSH, error)
	AddInstances(remoteUrl, PrivateKey string)
	RemoveInstance(uuid string)
	ReleaseInstances()
}

type RandomBalance struct {
	Instances *gmap.TreeMap
}

func (self *RandomBalance) DoBalance() (ssh *SSH, err error) {
	if self.Instances.IsEmpty() {
		err = errors.New("no instance")
		return
	} else if self.Instances.Size() == 1 {
		ssh = self.Instances.Values()[0].(*SSH)
		glog.Infof("get instance %s , index 0", ssh.Cfg.RemoteServer)
		return
	} else {
		index := rand.Intn(self.Instances.Size())
		ssh = self.Instances.Values()[index].(*SSH)
		glog.Infof("get instance %s , index %d", ssh.Cfg.RemoteServer, index)
		return
	}
}

func (self *RandomBalance) AddInstances(remoteUrl, PrivateKey string) {
	var (
		uuid = guid.S()
		err  error
	)
	Instance := &SSH{
		UUID: uuid,
		Cfg: struct {
			RemoteServer string
			PrivateKey   string
		}{RemoteServer: remoteUrl, PrivateKey: PrivateKey},
		CliCfg:   &ssh.ClientConfig{},
		StopChan: make(chan bool, 1),
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
		if Instance.URL.User == nil {
			glog.Errorf("%s not found user", Instance.URL)
			return
		}
		if pass, ok := Instance.URL.User.Password(); ok {
			Instance.CliCfg.Auth = append(Instance.CliCfg.Auth, ssh.Password(pass))
		}
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
	Instance.Direct = &Direct{
		Tr: &http.Transport{Dial: Instance.SSHDail},
	}
	Instance.Status = true
	// add
	glog.Infof("add %s ", Instance.URL)
	go Instance.KeepAlive()
	self.Instances.Set(uuid, Instance)
}

func (self *RandomBalance) RemoveInstance(uuid string) {
	value := self.Instances.Remove(uuid)
	instance := value.(*SSH)
	if err := instance.Stop(); err != nil {
		glog.Errorf("stop instance error: %s", err.Error())
	}
}

func (self *RandomBalance) ReleaseInstances() {
	for _, key := range self.Instances.Keys() {
		value := self.Instances.Remove(key)
		instance := value.(*SSH)
		if err := instance.Stop(); err != nil {
			glog.Errorf("stop instance error: %s", err.Error())
		}
	}
}

type RoundRobinBalance struct {
	Instances *gmap.TreeMap
	curIndex  int
}

func (self *RoundRobinBalance) DoBalance() (ssh *SSH, err error) {
	if self.Instances.IsEmpty() {
		err = errors.New("no instance")
		return
	} else if self.Instances.Size() == 1 {
		ssh = self.Instances.Values()[0].(*SSH)
		return
	} else {
		if self.curIndex >= self.Instances.Size() {
			self.curIndex = 0
		}
		ssh = self.Instances.Values()[self.curIndex].(*SSH)
		self.curIndex++
		return
	}
}

func (self *RoundRobinBalance) AddInstances(remoteUrl, PrivateKey string) {
	var (
		uuid = guid.S()
		err  error
	)
	Instance := &SSH{
		UUID: uuid,
		Cfg: struct {
			RemoteServer string
			PrivateKey   string
		}{RemoteServer: remoteUrl, PrivateKey: PrivateKey},
		CliCfg:   &ssh.ClientConfig{},
		StopChan: make(chan bool, 1),
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
		if Instance.URL.User == nil {
			glog.Errorf("%s not found user", Instance.URL)
			return
		}
		if pass, ok := Instance.URL.User.Password(); ok {
			Instance.CliCfg.Auth = append(Instance.CliCfg.Auth, ssh.Password(pass))
		}
	} else {
		signer, err := ssh.ParsePrivateKey([]byte(Instance.Cfg.PrivateKey))
		if err != nil {
			glog.Errorf("ParsePrivateKey %s failed:%s", Instance.Cfg.PrivateKey, err)
			return
		}
		Instance.CliCfg.Auth = append(Instance.CliCfg.Auth, ssh.PublicKeys(signer))
	}

	// init Client , first time to dial to remote server, make sure it is available
	Instance.Client, err = ssh.Dial("tcp", Instance.URL.Host, Instance.CliCfg)
	if err != nil {
		glog.Errorf("connect err: %s", err)
		return
	}
	Instance.Direct = &Direct{
		Tr: &http.Transport{Dial: Instance.SSHDail},
	}
	Instance.Status = true
	// add
	glog.Infof("add %s ", Instance.URL)
	go Instance.KeepAlive()
	self.Instances.Set(uuid, Instance)
}

func (self *RoundRobinBalance) RemoveInstance(uuid string) {
	value := self.Instances.Remove(uuid)
	instance := value.(*SSH)
	if err := instance.Stop(); err != nil {
		glog.Errorf("stop instance error: %s", err.Error())
	}
}

func (self *RoundRobinBalance) ReleaseInstances() {
	for _, key := range self.Instances.Keys() {
		value := self.Instances.Remove(key)
		instance := value.(*SSH)
		if err := instance.Stop(); err != nil {
			glog.Errorf("stop instance error: %s", err.Error())
		}
	}
}

func NewRandomBalance() *RandomBalance {
	return &RandomBalance{
		Instances: gmap.NewTreeMap(gutil.ComparatorString, true),
	}
}

func NewRoundRobinBalance() *RoundRobinBalance {
	return &RoundRobinBalance{
		Instances: gmap.NewTreeMap(gutil.ComparatorString, true),
	}
}
