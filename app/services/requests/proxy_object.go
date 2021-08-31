package requests

import (
	"homeproxy/app/models"
	"net/http"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

type ListObjectRequest struct {
	Pagination
	Prefix string `json:"prefix"`
}

func (self *ListObjectRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	bucket := r.GetString("bucket")
	if bucket == "" {
		return *response.ErrorWithMessage(http.StatusBadRequest, "No bucket specified")
	}
	query := g.DB().Model(&models.ObjectInfo{})

	query = query.Where("`bucket` = ?", bucket)

	if self.Prefix != "" {
		query = query.Where("`name` like %%s%%", self.Prefix)
	}
	var objects []models.ObjectInfo

	if err := query.Structs(&objects); err != nil {
		return *response.SystemError(err)
	}
	return *response.SuccessWithDetail(objects)
}

type UploadObjectRequest struct {
	Key string `json:"key" v:"required"`
}

func (self *UploadObjectRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	var (
		hash       = r.Header.Get("FILE_HASH")
		query      = g.DB().Model(&models.ObjectInfo{})
		object     = models.ObjectInfo{}
		bucketName = r.GetString("bucket")
		file       = r.GetUploadFile("file")
		bucket     = models.Bucket{}
	)

	if bucketName == "" {
		return *response.ErrorWithMessage(http.StatusBadRequest, "No bucket specified")
	}

	if err := g.DB().Model(&models.Bucket{}).Where("`name` = ?", bucketName).Struct(&bucket); err != nil {
		return *response.SystemError(err)
	}

	hashQuery := query.Clone().Where("`hash` = ?", hash)

	if c, err := hashQuery.Clone().Count(); err != nil {
		return *response.SystemError(err)
	} else {
		if c != 0 {
			err = hashQuery.Struct(&object)
			if err != nil {
				return *response.SystemError(err)
			}
			_, err = query.Clone().Data(object.CopyNewRecord(bucket.Id, file.Filename, self.Key)).Insert()
			if err != nil {
				return *response.SystemError(err)
			}
			return *response.Success()
		}
	}

	buffer := make([]byte, 512)

	header, err := file.Open()
	if err != nil {
		return *response.SystemError(err)
	}

	_, err = header.Read(buffer)
	if err != nil {
		return *response.SystemError(err)
	}

	version, _ := models.NewVersion("")
	object.ContextType = http.DetectContentType(buffer)
	object.Key = self.Key
	object.Name = file.Filename
	object.Bucket = bucket.Id
	object.RealPath = ""
	object.Hash = hash
	object.Size = file.Size
	object.Version = version.String()

	_, err = query.Clone().Save(&object)
	if err != nil {
		return *response.SystemError(err)
	}
	return *response.Success()
}

type DeleteObjectRequest struct {
	Key  string `json:"key" v:"required"`
	Name string `json:"name" v:"required"`
}

func (self *DeleteObjectRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	bucket := r.GetString("bucket")
	if bucket == "" {
		return *response.ErrorWithMessage(http.StatusBadRequest, "No bucket specified")
	}
	query := g.DB().Model(&models.ObjectInfo{})

	if _, err := query.Where(g.Map{"bucket": bucket, "key": self.Key, "name": self.Name}).Delete(); err != nil {
		return *response.SystemError(err)
	}
	return *response.Success()
}
