package requests

import (
	"github.com/gogf/gf/net/ghttp"
)

type ApplyUploadObjectRequest struct {
	Key      string `json:"key"`
	Hash     string `json:"hash"`
	Filename string `json:"filename"`
}

func (self *ApplyUploadObjectRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	return
}

type UploadObjectRequest struct{}

func (self *UploadObjectRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	// file := r.GetUploadFile("file")
	// key := r.GetString("key")
	// version := r.GetString("version")

	// obj := &models.ObjectInfo{}

	// filedb2.DB.Find("hash")
	return
}
