package tasks

import (
	"homeproxy/app/server/aria2"

)

func ReloadAira2Manager() {
	if aria2.Manager == nil {
		aria2.InitClient()
	} else {
		if aria2.Manager.Change {
			aria2.Manager.Close()
			aria2.InitClient()
		}
	}
}
