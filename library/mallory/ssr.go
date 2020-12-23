package mallory

import (
	"context"
	"net/http"
	"sync"

	"github.com/mzz2017/shadowsocksR/client"
	"github.com/nadoo/glider/proxy"
)

func init() {
	// ssr 暂时不添加，被封的概率很高，还是暂时不支持比较好
}

type SSRClient struct {
	UUID   string
	Url    string
	Direct *Direct
	Conn   proxy.Dialer
	sf     Group
	l      sync.RWMutex
	Status bool
	Cancel context.CancelFunc
}

func (self *SSRClient) Connect(w http.ResponseWriter, r *http.Request) {
	self.Direct.Connect(w, r)
}

func (self *SSRClient) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	self.Direct.ServeHTTP(w, r)
}

func (self *SSRClient) Stop() error {
	self.Cancel()
	return nil
}

func (self *SSRClient) Init () (err error) {
	self.Conn, err = client.NewSSRDialer(self.Url, proxy.Default)
	if err != nil {
		return
	}
	self.Direct = &Direct{Tr: &http.Transport{Dial: self.Conn.Dial}}
	return nil
}

