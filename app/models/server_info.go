package models

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gmeta"
	"golang.org/x/crypto/ssh"
)

type ServerGroup struct {
	DefaultModel
	Name       string `json:"name"`
	Remark     string `json:"remark"`
	gmeta.Meta `orm:"table:server_group"`
}

type Server struct {
	DefaultModel
	Name       string `json:"name"`
	Address    string `json:"address"`
	Port       int    `json:"port"`
	Group      int    `json:"group"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	PrivateKey string `json:"private_key"`
	Remark     string `json:"remark"`
	Status     bool   `json:"status"`
	UseProxy   bool   `json:"use_proxy"`
	ProxyID    int    `json:"proxy_id"`
	gmeta.Meta `orm:"table:server_host"`
}

func (s *Server) SshClient(conn net.Conn) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		Timeout:         time.Second * 5,
		User:            s.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	var auth []ssh.AuthMethod

	if s.Password == "" && s.PrivateKey == "" {
		return nil, errors.New("password and privatekey is empty")
	}
	if s.PrivateKey != "" {
		signer, err := ssh.ParsePrivateKey(gconv.Bytes(s.PrivateKey))
		if err != nil {
			return nil, err
		}
		auth = append(auth, ssh.PublicKeys(signer))
	} else {
		auth = append(auth, ssh.Password(s.Password))
	}
	config.Auth = auth
	addr := fmt.Sprintf("%s:%d", s.Address, s.Port)

	if conn != nil {
		c, chans, reqs, err := ssh.NewClientConn(conn, addr, config)
		if err != nil {
			return nil, err
		}
		return ssh.NewClient(c, chans, reqs), nil
	}
	c, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Server) GetProxyClient() (*ssh.Client, error) {
	if !s.UseProxy {
		return nil, errors.New("当前连接未配置代理")
	}
	glog.Info("use ssh proxy")
	var server = &Server{}
	err := g.DB().Model(Server{}).Where(`id = ?`, s.ProxyID).Scan(server)
	if err != nil {
		return nil, err
	}

	proxyClient, err := server.SshClient(nil)

	if err != nil {
		return nil, err
	}
	addr := fmt.Sprintf("%s:%d", s.Address, s.Port)

	c, err := proxyClient.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return s.SshClient(c)
}

func (s *Server) GetSshClient() (*ssh.Client, error) {
	glog.Info(s)
	if s.UseProxy {
		return s.GetProxyClient()
	}
	return s.SshClient(nil)
}
