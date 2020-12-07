package tasks

import (
	"fmt"
	"github.com/gogf/gf/os/gcron"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	"homeproxy/app/models"
	"homeproxy/library/ddns"
	"homeproxy/library/filedb"
	"time"
)

func init() {
	InitDDnsTask()
}
func InitDDnsTask() {
	var roles []models.DDnsProviderSettings
	err := filedb.DB.QueryAll(models.DDnsProviderSettingsTable, &roles)
	if err != nil {
		glog.Error("init ddns task error: %s", err.Error())
	} else {
		for _, role := range roles {
			if role.TimeInterval != "" && role.Status {
				glog.Infof("start ddns task %s, %s.%s", role.ID, role.SubDomain, role.Domain)
				DDnsSyncTask(role.ID)()
				gcron.AddSingleton(role.TimeInterval, DDnsSyncTask(role.ID), role.ID)
			}
		}
	}
}

func DDnsSyncTask(roleID string) func() {
	return func() {
		var (
			addresses []ddns.Address
			err       error
		)
		setting := models.DDnsProviderSettings{}
		err = filedb.DB.GetById(models.DDnsProviderSettingsTable, roleID, &setting)
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

			switch setting.Provider {
			case "dnspod":
				provider = ddns.DefaultDnsPod(setting.DNSPodID, setting.DNSPodToken)
			default:
				glog.Error("unknown provider: ", setting.Provider)
				History.Error = fmt.Sprintf("unknown provider: %s", setting.Provider)
				History.Value = addr.String()
				History.Status = 1
				History.Date = time.Now().Format("2006-01-02 15:04:05")
				setting.History = append(setting.History, History)
				filedb.DB.UpdateById(models.DDnsProviderSettingsTable, setting.ID, setting)
				return
			}
			switch setting.RecordID {
			case "":
				recordID, err := provider.RecordCreate(setting.Domain, setting.SubDomain, "600", addr)
				if err != nil {
					glog.Errorf("create record error: %s", err.Error())
					History.Error = err.Error()
					History.Value = addr.String()
					History.Status = 1
					History.Date = time.Now().Format("2006-01-02 15:04:05")
					setting.History = append(setting.History, History)
					filedb.DB.UpdateById(models.DDnsProviderSettingsTable, setting.ID, setting)
				} else {
					History.Error = ""
					History.Value = addr.String()
					History.Status = 0
					History.Date = time.Now().Format("2006-01-02 15:04:05")
					setting.RecordID = recordID
					setting.History = append(setting.History, History)
					filedb.DB.UpdateById(models.DDnsProviderSettingsTable, setting.ID, setting)
				}
			default:
				err := provider.RecordModify(setting.Domain, setting.RecordID, setting.SubDomain, addr)
				if err != nil {
					glog.Errorf("create record error: %s", err.Error())
					History.Error = err.Error()
					History.Value = addr.String()
					History.Status = 1
					History.Date = time.Now().Format("2006-01-02 15:04:05")
					setting.History = append(setting.History, History)
					if err := filedb.DB.UpdateById(models.DDnsProviderSettingsTable, setting.ID, setting); err != nil {
						glog.Errorf("ddns save data error: %s", err.Error())
					}
				} else {
					History.Error = ""
					History.Value = addr.String()
					History.Status = 0
					History.Date = time.Now().Format("2006-01-02 15:04:05")
					setting.History = append(setting.History, History)
					if err := filedb.DB.UpdateById(models.DDnsProviderSettingsTable, setting.ID, setting); err != nil {
						glog.Errorf("ddns save data error: %s", err.Error())
					}
				}
			}
		}
	}
}
