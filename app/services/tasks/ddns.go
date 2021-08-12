package tasks

import (
	"fmt"
	"github.com/gogf/gf/os/gcron"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	"homeproxy/app/models"
	"homeproxy/library/ddns"
	"homeproxy/library/filedb2"
	"time"
)

func InitDDnsTask() {
	var roles []models.DDnsProviderSettings
	err := filedb2.DB.All(&roles)
	if err != nil {
		glog.Error("init ddns task error: %s", err.Error())
	} else {
		for _, role := range roles {
			if role.TimeInterval != "" && role.Status {
				glog.Infof("start ddns task %s, %s.%s", role.ID, role.SubDomain, role.Domain)
				DDnsSyncTask(role.ID)()
				gcron.AddSingleton(role.TimeInterval, DDnsSyncTask(role.ID), gconv.String(role.ID))
			}
		}
	}
}

func DDnsSyncTask(roleID int) func() {
	return func() {
		var (
			addresses []ddns.Address
			err       error
		)
		setting := models.DDnsProviderSettings{}
		err = filedb2.DB.One("ID", roleID, &setting)
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
			)
			History.Value = addr.String()
			History.Date = time.Now().Format("2006-01-02 15:04:05")

			switch setting.Provider {
			case "dnspod":
				provider = ddns.DefaultDnsPod(setting.DNSPodID, setting.DNSPodToken)
			default:
				glog.Error("unknown provider: ", setting.Provider)
				History.Error = fmt.Sprintf("unknown provider: %s", setting.Provider)
				History.Status = 1
				if len(setting.History) >= 5 {
					setting.History = setting.History[1:]
				}
				setting.History = append(setting.History, History)
				filedb2.DB.Update(&setting)
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
			if len(setting.History) >= 5 {
				setting.History = setting.History[1:]
			}
			setting.History = append(setting.History, History)
			filedb2.DB.Update(&setting)
		}
	}
}
