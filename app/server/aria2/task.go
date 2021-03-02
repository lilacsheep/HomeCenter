package aria2

import (
	"strings"

	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gcron"
	"github.com/gogf/gf/os/glog"
	"gopkg.in/ini.v1"
)


const (
	AutoUpdateBTTrackerName = "AutoUpdateBTTracker"
	TrackersAllUrl = "https://raw.githubusercontent.com/ngosang/trackerslist/master/trackers_all.txt"
)

func AutoUpdateBTTracker() {
	client := ghttp.NewClient()
	response, err := client.Get(TrackersAllUrl)
	if err != nil {
		glog.Errorf("get trackers server error: %v", err)
	} else {
		defer response.Close()
		var servers []string
		data := response.ReadAllString()
		for _, server := range strings.Split(data, "\n") {
			if server != "" {
				servers = append(servers, server)
			}
		}
		if len(servers) != 0 {
			configPath, err := Manager.GetOption(ConfigPathKey)
			if err != nil {
				glog.Errorf("get aria2 config path error: %v", err)
				return
			}
			cfg, err := ini.Load(configPath)
			if err != nil {
				glog.Errorf("Fail to read file: %v", err)
			}
			cfg.Section("").Key(ConfigBTTrackerKey).SetValue(strings.Join(servers, ","))
			err = cfg.SaveTo(configPath)
			if err != nil {
				glog.Errorf("save aria2 config error: %v", err)
				return
			}
			glog.Info("update trackers server successfully")
		}
	}
}

func NewAutoUpdateBTTracker(time_interval string) {
	entry := gcron.Search(AutoUpdateBTTrackerName)
	if entry != nil {
		entry.Close()
	}
	if time_interval != "" { 
		gcron.AddSingleton(time_interval, AutoUpdateBTTracker, AutoUpdateBTTrackerName)
	}
}