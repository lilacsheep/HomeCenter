package mallory

import (
	"fmt"
	"homeproxy/library/events"
	"homeproxy/library/filedb2"
	"strings"
	"time"

	"github.com/asdine/storm/v3"
	"golang.org/x/net/publicsuffix"
)

type VisitLogEvent struct {
	Address  string
	Host     string
	DateTime time.Time
}

func (self *VisitLogEvent) DoEvent() error {
	log := ProxyVisitLog{}
	log.Address = self.Address
	log.Host = self.Host
	log.Datetime = self.DateTime
	return filedb2.DB.Save(&log)
}

type DomainErrorEvent struct {
	Domain string
	Error  string
}

func (self *DomainErrorEvent) urlSplit(url string) (string, string) {
	host := HostOnly(url)
	domain, _ := publicsuffix.EffectiveTLDPlusOne(host)
	subDomain := ""
	t := strings.Split(host, fmt.Sprintf(".%s", domain))
	subDomain = t[0]
	if subDomain == domain {
		return "", domain
	}
	return subDomain, domain
}

func (self *DomainErrorEvent) DoEvent() error {
	_, domain := self.urlSplit(self.Domain)
	var data = ProxyRoleAnalysis{}
	err := filedb2.DB.One("Domain", domain, &data)
	if err != nil {
		if err != storm.ErrNotFound {
			return err
		} else {
			data.Domain = domain
			data.Times = 1
			data.Error = self.Error
			err = filedb2.DB.Save(&data)
			if err != nil {
				return err
			}
		}
	} else {
		data.Times += 1
		data.Error = self.Error
		err = filedb2.DB.Update(&data)
		if err != nil {
			return err
		}
	}
	return nil
}

func errorEvent(host string, err error) {
	event := &DomainErrorEvent{host, err.Error()}
	events.EventChan <- event
}

func visitLogEvent(address, host string) {
	event := &VisitLogEvent{Address: address, Host: host, DateTime: time.Now()}
	events.EventChan <- event
}
