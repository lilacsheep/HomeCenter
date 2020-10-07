package events

import (
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/grpool"
)

var EventChan chan EventMessageInterface

func init() {
	EventChan = make(chan EventMessageInterface, 99999)
}

type EventMessageInterface interface {
	DoEvent() error
}

type EventProcess struct {
	Pool *grpool.Pool
}

func (self *EventProcess) Run() {
	for {
		select {
		case event, ok := <-EventChan:
			if ok {
				self.Pool.Add(func() {
					err := event.DoEvent()
					if err != nil {
						glog.Errorf("do event %v %s", event, err)
					}
				})
			}
		}
	}
}
