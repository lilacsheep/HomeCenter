package requests

import (
	"fmt"
	"homeproxy/app/models"
	"homeproxy/app/services/tasks"
	"homeproxy/library/ddns"
	"net/http"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gcron"
	"github.com/gogf/gf/util/gconv"
	"github.com/shirou/gopsutil/net"
)

type ChooseNetCardRequest struct{}

func (self *ChooseNetCardRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	cards, err := net.Interfaces()
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		response.SuccessWithDetail(cards)
	}
	return
}

func NewChooseNetCardRequest() *ChooseNetCardRequest {
	return &ChooseNetCardRequest{}
}

type GetDDnsSettingsRequest struct{}

func (self *GetDDnsSettingsRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	var data []models.DDnsProviderSettings
	err := g.DB().Model(&models.DDnsProviderSettings{}).Structs(&data)
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		if data == nil {
			response.DataTable([]models.DDnsProviderSettings{}, 0)
		} else {
			response.DataTable(data, len(data))
		}
	}
	return
}

func NewGetDDnsSettingsRequest() *GetDDnsSettingsRequest {
	return &GetDDnsSettingsRequest{}
}

type CreateDDnsSettingRequest struct {
	Domain       string `json:"domain"`
	SubDomain    string `json:"sub_domain"`
	Provider     string `json:"provider"`
	TimeInterval string `json:"time_interval"`
	Mode         bool   `json:"mode"`
	NetCard      int    `json:"net_card"`
	UpdatedOn    string `json:"updated_on"`
	RecordID     string `json:"record_id"`
	DNSPodID     string `json:"dnspod_id"`
	DNSPodToken  string `json:"dnspod_token"`
}

func (self *CreateDDnsSettingRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	setting := models.DDnsProviderSettings{}
	setting.Domain = self.Domain
	setting.SubDomain = self.SubDomain
	setting.Provider = self.Provider
	setting.UpdatedOn = time.Now().Format("2006-01-02 15:04:05")
	setting.RecordID = self.RecordID
	setting.TimeInterval = self.TimeInterval
	setting.Status = true
	if !self.Mode {
		setting.UsePublicIP = true
	} else {
		setting.NetCard = self.NetCard
	}
	switch setting.Provider {
	case "dnspod":
		setting.DNSPodID = self.DNSPodID
		setting.DNSPodToken = self.DNSPodToken
	}

	_, err := g.DB().Model(&models.DDnsProviderSettings{}).Save(&setting)
	if err != nil {
		response.SystemError(err)
	} else {
		tasks.DDnsSyncTask(setting.Id)()
		_, err = gcron.AddSingleton(setting.TimeInterval, tasks.DDnsSyncTask(setting.Id), gconv.String(setting.Id))
		if err != nil {
			response.SystemError(err)
		} else {
			response.SuccessWithDetail(setting.Id)
		}
	}
	return
}

func NewCreateDDnsSettingRequest() *CreateDDnsSettingRequest {
	return &CreateDDnsSettingRequest{}
}

type GetProviderRecordsRequest struct {
	Provider    string `json:"provider"`
	Domain      string `json:"domain" v:"domain  @required#请输入域名"`
	SubDomain   string `json:"sub_domain"`
	DNSPodID    string `json:"dnspod_id" v:"dnspod_id  @required#请输入ID"`
	DNSPodToken string `json:"dnspod_token" v:"dnspod_token  @required#请输入Token"`
}

func (self *GetProviderRecordsRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	switch self.Provider {
	case "dnspod":
		provider := ddns.DefaultDnsPod(self.DNSPodID, self.DNSPodToken)
		resp, err := provider.RecordList(self.Domain, self.SubDomain)
		if err != nil {
			response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
		} else {
			dnsPodResp := resp.(ddns.DNSPodResponse)
			if dnsPodResp.Status.Code == 1 {
				response.SuccessWithDetail(dnsPodResp.Records)
			} else if dnsPodResp.Status.Code == 10 {
				response.SuccessWithDetail([]ddns.Record{})
			} else {
				response.ErrorWithMessage(http.StatusInternalServerError, dnsPodResp.Status.Message)
			}
		}
	default:
		response.ErrorWithMessage(http.StatusInternalServerError, fmt.Sprintf("未知供应商: %s", self.Provider))
	}

	return
}

func NewGetProviderRecordsRequest() *GetProviderRecordsRequest {
	return &GetProviderRecordsRequest{}
}

type GetSettingsInfoRequest struct {
	ID int `json:"id" v:"required"`
}

func (self *GetSettingsInfoRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	info := g.Map{}
	entryInfo := g.Map{"name": "", "time": "", "status": -1}
	entry := gcron.Search(gconv.String(self.ID))
	if entry != nil {
		entryInfo["name"] = entry.Name
		entryInfo["time"] = entry.Time.Format("2006-01-02 15:04:05")
		entryInfo["status"] = entry.Status()
	}
	info["entry"] = entryInfo
	var setting models.DDnsProviderSettings
	err := g.DB().Model(&models.DDnsProviderSettings{}).Where("`id` = ?", self.ID).Struct(&setting)
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		info["setting"] = setting
		response.SuccessWithDetail(info)
	}
	return
}

func NewGetSettingsInfoRequest() *GetSettingsInfoRequest {
	return &GetSettingsInfoRequest{}
}

type StopRoleTaskRequest struct {
	ID int `json:"id" v:"required"`
}

func (self *StopRoleTaskRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	entry := gcron.Search(gconv.String(self.ID))
	if entry == nil {
		response.ErrorWithMessage(http.StatusInternalServerError, "任务不存在")
	} else {
		var (
			setting models.DDnsProviderSettings
			query   = g.DB().Model(&models.DDnsProviderSettings{})
		)
		err := query.Clone().Where("`id` = ?", self.ID).Struct(&setting)
		if err != nil {
			response.ErrorWithMessage(http.StatusInternalServerError, "任务不存在")
		} else {
			_, err := query.Clone().Data(g.Map{"status": false}).Where("`id` = ?", self.ID).Update()
			if err != nil {
				response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
			}
		}
		entry.Close()
		response.Success()
	}
	return
}

func NewStopRoleTaskRequest() *StopRoleTaskRequest {
	return &StopRoleTaskRequest{}
}

type StartRoleTaskRequest struct {
	ID int `json:"id" v:"required"`
}

func (self *StartRoleTaskRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	entry := gcron.Search(gconv.String(self.ID))
	if entry == nil {
		var (
			setting models.DDnsProviderSettings
			query   = g.DB().Model(&models.DDnsProviderSettings{})
		)
		err := query.Clone().Where("`id` = ?", self.ID).Struct(&setting)
		if err != nil {
			response.ErrorWithMessage(http.StatusInternalServerError, "任务不存在")
		} else {
			_, err = gcron.AddSingleton(setting.TimeInterval, tasks.DDnsSyncTask(setting.Id), gconv.String(setting.Id))
			if err != nil {
				response.ErrorWithMessage(http.StatusInternalServerError, "任务不存在")
			} else {
				_, err := query.Clone().Data(g.Map{"status": true}).Where("`id` = ?", self.ID).Update()
				if err != nil {
					response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
				} else {
					response.Success()
				}
			}
		}
	} else {
		entry.Start()
		response.Success()
	}
	return
}

func NewStartRoleTaskRequest() *StartRoleTaskRequest {
	return &StartRoleTaskRequest{}
}

type DeleteSettingRequest struct {
	ID int `json:"id" v:"required"`
}

func (self *DeleteSettingRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	entry := gcron.Search(gconv.String(self.ID))
	if entry != nil {
		entry.Close()
	}
	_, err := g.DB().Model(&models.DDnsProviderSettings{}).Where("`id` = ?", self.ID).Delete()
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		response.Success()
	}
	return
}

func NewDeleteSettingRequest() *DeleteSettingRequest {
	return &DeleteSettingRequest{}
}

type RefreshSettingRequest struct {
	ID           int    `json:"id" v:"required"`
	TimeInterval string `json:"time_interval" v:"required"`
}

func (self *RefreshSettingRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	entry := gcron.Search(gconv.String(self.ID))
	if entry != nil {
		entry.Close()
	}
	var (
		setting models.DDnsProviderSettings
		query   = g.DB().Model(&models.DDnsProviderSettings{})
	)
	err := query.Clone().Where("`id` = ?", self.ID).Struct(&setting)
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
		return
	}
	setting.TimeInterval = self.TimeInterval
	_, err = query.Clone().Data(g.Map{"time_interval": self.TimeInterval}).Where("`id` = ?", self.ID).Update()
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		_, err = gcron.AddSingleton(self.TimeInterval, tasks.DDnsSyncTask(self.ID), gconv.String(self.ID))
		if err != nil {
			response.ErrorWithMessage(http.StatusInternalServerError, "添加任务失败")
		} else {
			response.Success()
		}
	}
	return
}

func NewRefreshSettingRequest() *RefreshSettingRequest {
	return &RefreshSettingRequest{}
}
