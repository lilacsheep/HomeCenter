package api

import (
	"bytes"
	"fmt"
	"homeproxy/app/models"
	"homeproxy/app/services/requests"
	"strings"
	"time"

	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
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

	_, m, err := ws.ReadMessage()
	if err != nil {
		glog.Error(err)
		ws.WriteMessage(-1, []byte(err.Error()))
		return
	}
	d := gjson.New(m)
	cols := d.GetInt("cols", 120)
	rows := d.GetInt("rows", 32)
	hostId := d.GetString("host")
	if hostId == "" {
		glog.Error(hostId)
		ws.WriteMessage(-1, []byte("为获取到主机id"))
		return
	}
	t := strings.Split(hostId, "-")
	host := models.Server{}

	err = g.DB().Model(&models.Server{}).Where("`id` = ?", t[1]).Struct(&host)
	if err != nil {
		glog.Error(err)
		ws.WriteMessage(-1, []byte(err.Error()))
		return
	}
	client, err := NewSshClient(host.Address, host.Port)
	if err != nil {
		glog.Error(err)
		ws.WriteMessage(-1, []byte(err.Error()))
		return
	}

	
	ssConn, err := NewSshConn(cols, rows, client)
	if err != nil {
		ws.WriteMessage(-1, []byte(err.Error()))
		return
	}
	defer ssConn.Close()

	glog.Info(cols, rows)
	quitChan := make(chan bool, 3)

	var logBuff = new(bytes.Buffer)

	// most messages are ssh output, not webSocket input
	go ssConn.ReceiveWsMsg(ws.Conn, logBuff, quitChan)
	go ssConn.SendComboOutput(ws.Conn, quitChan)
	go ssConn.SessionWait(quitChan)
	<-quitChan

}

func NewSshClient(host string, port int) (*ssh.Client, error) {
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
	addr := fmt.Sprintf("%s:%d", host, port)
	c, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}
	return c, nil
}
