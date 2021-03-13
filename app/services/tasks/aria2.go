package tasks

import (
	"homeproxy/app/server/aria2"

	"github.com/gogf/gf/os/gcron"
)

func init() {
	gcron.AddSingleton("*/2 * * * * *", ReloadAira2Manager)
}

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
