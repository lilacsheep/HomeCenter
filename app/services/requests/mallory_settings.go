package requests

import (
	"github.com/gogf/gf/net/ghttp"
	"homeproxy/library/mallory"
)

type MallorySettingsQueryRequest struct {
	mallory.ConfigFile
}

func (request *MallorySettingsQueryRequest) Exec(r *ghttp.Request) (response MessageResponse) {

	return
}
