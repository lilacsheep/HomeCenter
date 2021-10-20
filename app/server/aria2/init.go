package aria2

import (
	"homeproxy/app/models"
)

const (
	ConfigPathKey      = "conf-path"
	ConfigBTTrackerKey = "bt-tracker"
)

func InitClient() error {
	settings, err := models.GetAria2Settings()
	if err != nil {
		return err
	}
	Manager = &manager{Settings: settings}

	return Manager.Init()
}
