package sshm

import (
	"math"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

type cpuHistroy struct {
	Idle float64
	Busy float64
}

type Session struct {
	session *ssh.Session
	l       sync.Mutex
}

func (s *Session) Output(cmd string) ([]byte, error) {
	s.l.Lock()
	defer s.l.Unlock()	
	return s.session.Output(cmd)
}

func (c *cpuHistroy) Percent(idle, busy float64) float64 {

	if busy <= c.Busy {
		return 0
	}
	if (idle + busy) <= (c.Idle + c.Busy) {
		return 100
	}
	p := math.Min(100, math.Max(0, (busy-c.Busy)/((idle+busy)-(c.Idle+c.Busy))*100))

	c.Busy = busy
	c.Idle = idle
	return p
}

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func (m Message) Byte() ([]byte, error) {
	return gvar.New(m).MarshalJSON()
}

func NewSshConnMonitor(client *ssh.Client, ws *websocket.Conn) (*SshClientWs, error) {
	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	return &SshClientWs{
		Client:  client,
		Ws:      ws,
		Chan:    make(chan Message, 999999),
		exit:    false,
		session: &Session{session: session, l: sync.Mutex{}},
	}, nil
}

type SshClientWs struct {
	Client  *ssh.Client
	Ws      *websocket.Conn
	Chan    chan Message
	exit    bool
	session *Session
}

func (c *SshClientWs) Close() {
	c.exit = true
	close(c.Chan)
}

func (c *SshClientWs) Delay() error {
	for {
		if c.exit {
			return nil
		}
		t1 := time.Now()
		_, err := c.session.Output("echo 0")
		if err != nil {
			return err
		}
		delay := time.Now().Sub(t1).Milliseconds()

		c.Chan <- Message{Type: "delay", Data: delay}
		time.Sleep(time.Second)
	}
}

func (c *SshClientWs) Cpu() {
	histroy := &cpuHistroy{}

	cmd := `cat /proc/stat|grep cpu|awk '{print $2,$3,$4,$5,$6,$7,$8}'`

	for {
		if c.exit {
			return
		}
		out, err := c.session.Output(cmd)
		if err != nil {
			glog.Info(err)
			time.Sleep(time.Second)
			continue
		}
		o := string(out)
		var busy float64
		var idle float64
		for _, l := range strings.Split(o, "\n") {
			for i, n := range strings.Split(l, " ") {
				if i == 3 {
					idle += gconv.Float64(n)
				} else {
					busy += gconv.Float64(n)
				}
			}
		}
		p := histroy.Percent(idle, busy)
		msg := Message{Type: "cpu", Data: p}
		glog.Info(msg)
		c.Chan <- msg
		time.Sleep(time.Second)
	}
}

func (c *SshClientWs) Run() {
	go c.Delay()
	go c.Cpu()
	for {
		select {
		case msg, ok := <-c.Chan:
			if ok {
				d, err := msg.Byte()
				if err != nil {
					glog.Error(err)
				} else {
					err = c.Ws.WriteMessage(websocket.TextMessage, d)
					if err != nil {
						glog.Error(err)
					}
				}
			} else {
				return
			}
		}
	}
}
