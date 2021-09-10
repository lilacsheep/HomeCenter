package api

import (
	"bytes"
	"fmt"
	"homeproxy/app/services/requests"
	"time"

	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"golang.org/x/crypto/ssh"
)

type SystemApi struct {
	BaseControllers
}

func (a *SystemApi) Info(r *ghttp.Request) {
	request := &requests.SystemInfoRequest{}
	a.DoRequest(request, r)
}

func (a *SystemApi) Processes(r *ghttp.Request) {
	request := &requests.ProcessesRquest{}
	a.DoRequest(request, r)
}

func (a *SystemApi) Webssh(r *ghttp.Request) {
	ws, err := r.WebSocket()
	if err != nil {
		glog.Error(err)
		r.Exit()
	}
	defer ws.Close()

	cols := 120
	rows := 32

	client, err := NewSshClient()

	if err != nil {
		ws.WriteMessage(-1, []byte(err.Error()))
		return
	}

	ssConn, err := NewSshConn(cols, rows, client)
	if err != nil {
		ws.WriteMessage(-1, []byte(err.Error()))
		return
	}
	defer ssConn.Close()

	quitChan := make(chan bool, 3)

	var logBuff = new(bytes.Buffer)

	// most messages are ssh output, not webSocket input
	go ssConn.ReceiveWsMsg(ws.Conn, logBuff, quitChan)
	go ssConn.SendComboOutput(ws.Conn, quitChan)
	go ssConn.SessionWait(quitChan)
	<-quitChan

}

func NewSshClient() (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		Timeout:         time.Second * 5,
		User:            "root",
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	//if h.Type == "password" {
	config.Auth = []ssh.AuthMethod{ssh.Password("asdf3.14")}
	//} else {
	//	config.Auth = []ssh.AuthMethod{publicKeyAuthFunc(h.Key)}
	//}
	addr := fmt.Sprintf("%s:%d", "192.168.2.22", 22)
	c, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}
	return c, nil
}
