package api

import (
	"errors"
	"homeproxy/app/models"
	"homeproxy/app/services/requests"
	"homeproxy/app/services/tasks"
	"strings"

	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	"golang.org/x/crypto/ssh"
)

type SystemApi struct {
	BaseControllers
}

func (a *SystemApi) Info(r *ghttp.Request) {
	request := &requests.SystemInfoRequest{}
	a.DoRequest(request, r)
}

func (a *SystemApi) InfoWs(r *ghttp.Request) {

	ws, err := r.WebSocket()
	if err != nil {
		glog.Error(err)
		return
	}
	
	for {
		_, m, err := ws.ReadMessage()
		if err != nil {
			glog.Error(err)
			return
		}
		switch gconv.String(m) {
		case "status":
			h := tasks.History.Clone()
			ws.WriteJSON(g.Map{"type": "status", "data": h})
		default:
			continue
		}
	}
	
}
func (a *SystemApi) Processes(r *ghttp.Request) {
	request := &requests.ProcessesRquest{}
	a.DoRequest(request, r)
}

func (a *SystemApi) Process(r *ghttp.Request) {
	request := &requests.ProcessInfoRequest{}
	a.DoRequestValid(request, r)
}

func (a *SystemApi) Webssh(r *ghttp.Request) {
	type Message struct {
		Type string
		Data string
	}

	connInit := func(ws *ghttp.WebSocket) (*ssh.Client, int, int, error) {
		_, m, err := ws.ReadMessage()
		if err != nil {
			return nil, 0, 0, err
		}
		d := gjson.New(m)
		cols := d.GetInt("cols", 120)
		rows := d.GetInt("rows", 32)
		hostId := d.GetString("host")
		if hostId == "" {
			return nil, 0, 0, errors.New("未获取到主机id")
		}
		t := strings.Split(hostId, "-")
		host := models.Server{}
		err = g.DB().Model(&models.Server{}).Where("`id` = ?", t[1]).Struct(&host)
		if err != nil {
			return nil, 0, 0, err
		}
		cli, err := host.GetSshClient()
		if err != nil {
			return nil, 0, 0, err
		}
		return cli, cols, rows, nil
	}

	ws, err := r.WebSocket()
	if err != nil {
		glog.Error(err)
		r.Exit()
	}
	defer ws.Close()
	client, cols, rows, err := connInit(ws)

	if err != nil {
		msg := gjson.New(g.Map{"type": "error", "message": err.Error()})
		ws.WriteMessage(1, msg.MustToJson())
		return
	}
	ssConn, err := NewSshConn(cols, rows, client)
	if err != nil {
		msg := gjson.New(g.Map{"type": "error", "message": err.Error()})
		ws.WriteMessage(1, msg.MustToJson())
		return
	}
	defer ssConn.Close()
	msg := gjson.New(g.Map{"type": "success", "message": ""})

	ws.WriteMessage(1, msg.MustToJson())
	quitChan := make(chan bool, 3)

	// mc, err := sshm.NewSshConnMonitor(client, ws.Conn)
	// if err != nil {
	// 	msg := gjson.New(g.Map{"type": "error", "message": err.Error()})
	// 	ws.WriteMessage(1, msg.MustToJson())
	// 	return
	// }
	// go mc.Run()
	// defer mc.Close()
	// most messages are ssh output, not webSocket input
	go ssConn.ReceiveWsMsg(ws.Conn, nil, quitChan)
	go ssConn.SendComboOutput(ws.Conn, quitChan)
	go ssConn.SessionWait(quitChan)
	<-quitChan

}
