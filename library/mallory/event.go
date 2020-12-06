package mallory

import (
	"github.com/gogf/gf/frame/g"
	"homeproxy/library/events"
	"homeproxy/library/filedb"
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

func (self *DomainErrorEvent) DoEvent() error {
	params := g.Map{"domain": self.Domain, "error": self.Error}
	var data = ProxyRoleAnalysis{}
	err := filedb.DB.QueryOne(ProxyRoleAnalysisTable, &data, params)
	if err != nil {
		if err != filedb.ErrNoData {
			return err
		} else {
			data.Domain = self.Domain
			data.Times = 1
			data.Error = self.Error
			_, err := filedb.DB.Insert(ProxyRoleAnalysisTable, data)
			if err != nil {
				return err
			}
		}
	} else {
		data.Times += 1
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
