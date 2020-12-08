package mallory

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"golang.org/x/net/publicsuffix"
	"homeproxy/library/events"
	"homeproxy/library/filedb"
	"strings"
	"time"
)

const (
	ProxyRoleTable         = "proxy_role_list"
	ProxyRoleAnalysisTable = "proxy_role_analysis_url"
	ProxyVisitLogTable     = "proxy_visit_logs"
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
	log.Datetime = self.DateTime.Format("2006-01-02 15:04:05")
	_, err := filedb.DB.Insert(ProxyVisitLogTable, log)
	if err != nil {
		return err
	}
	return nil
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
	params := g.Map{"domain": domain}
	var data = ProxyRoleAnalysis{}
	err := filedb.DB.QueryOne(ProxyRoleAnalysisTable, &data, params)
	if err != nil {
		if err != filedb.ErrNoData {
			return err
		} else {
			data.Domain = domain
			data.Times = 1
			data.Error = self.Error
			_, err := filedb.DB.Insert(ProxyRoleAnalysisTable, data)
			if err != nil {
				return err
			}
		}
	} else {
		data.Times += 1
		data.Error = self.Error
		err := filedb.DB.UpdateById(ProxyRoleAnalysisTable, data.ID, data)
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
