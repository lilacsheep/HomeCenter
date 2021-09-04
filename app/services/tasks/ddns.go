package tasks

import (
	"fmt"
	"homeproxy/app/models"
	"homeproxy/library/ddns"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcron"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
)

func InitDDnsTask() {
	settings, err := models.AllDDnsSettings()
	if err != nil {
		glog.Error("init ddns task error: %s", err.Error())
	} else {
		for _, s := range settings {
			if s.TimeInterval != "" && s.Status {
				glog.Infof("start ddns task %s, %s.%s", s.Id, s.SubDomain, s.Domain)
				DDnsSyncTask(s.Id)()
				gcron.AddSingleton(s.TimeInterval, DDnsSyncTask(s.Id), gconv.String(s.Id))
			}
		}
	}
}

func DDnsSyncTask(roleID int) func() {
	return func() {
		var (
			addresses []ddns.Address
			err       error
			query     = g.DB().Model(&models.DDnsProviderSettings{})
		)
		setting := models.DDnsProviderSettings{}
		err = query.Where("`id` = ?", roleID).Struct(&setting)
		if err != nil {
			glog.Errorf("error: %s", err.Error())
			return
		}
		if setting.UsePublicIP {
			address, err := ddns.GetPublicAddress()
			if err != nil {
				glog.Errorf("error: %s", err.Error())
				return
			}
			addresses = append(addresses, address)
		} else {
			addr, err := ddns.GetPublicV6AddressByInterface(gconv.Int(setting.NetCard))
			if err != nil {
				glog.Info(err)
				return
			} else {
				addresses = append(addresses, addr...)
			}
		}
		for _, addr := range addresses {
			var (
				provider ddns.Provider
				History  = models.OperationRecord{}
				recordHistoryQuery = g.DB().Model(&models.OperationRecord{})
			)
			History.SettingId = setting.Id
			History.Value = addr.String()
			History.Date = time.Now().Format("2006-01-02 15:04:05")
			History.SettingId = setting.Id

			switch setting.Provider {
			case "dnspod":
				provider = ddns.DefaultDnsPod(setting.DNSPodID, setting.DNSPodToken)
			default:
				glog.Error("unknown provider: ", setting.Provider)
				History.Error = fmt.Sprintf("unknown provider: %s", setting.Provider)
				History.Status = 1
				recordHistoryQuery.Clone().Save(&History)
				return
			}
			switch setting.RecordID {
			case "":
				recordID, err := provider.RecordCreate(setting.Domain, setting.SubDomain, "600", addr)
				if err != nil {
					glog.Errorf("create record error: %s", err.Error())
					History.Error = err.Error()
					History.Status = 1
				} else {
					History.Status = 0
					setting.RecordID = recordID
				}
			default:
				err := provider.RecordModify(setting.Domain, setting.RecordID, setting.SubDomain, addr)
				if err != nil {
					glog.Errorf("create record error: %s", err.Error())
					History.Error = err.Error()
					History.Status = 1
				} else {
					History.Status = 0
				}
			}
			recordHistoryQuery.Clone().Save(&History)
		}
	}
}
