package mallory

import (
	"context"
	"net/http"
	"sync"

	"github.com/mzz2017/shadowsocksR/client"
	"github.com/nadoo/glider/proxy"
)

func init() {

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

func (self *SSRClient) Init () (err error) {
	self.Conn, err = client.NewSSRDialer(self.Url, proxy.Default)
	if err != nil {
		return
	}
	self.Direct = &Direct{Tr: &http.Transport{Dial: self.Conn.Dial}}
	return nil
}

