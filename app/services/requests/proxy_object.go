package requests

import (
	"homeproxy/app/models"
	"homeproxy/library/config"
	"net/http"
	"path"

	"github.com/gogf/gf/crypto/gmd5"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gfile"
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
		hash       = r.Header.Get("FILE_MD5")
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

	savePath := path.Join(gfile.Abs(config.DataDir), bucket.Name)

	version, _ := models.NewVersion("")
	object.ContextType = http.DetectContentType(buffer)
	object.Key = self.Key
	object.Name = file.Filename
	object.Bucket = bucket.Id
	object.Hash = hash
	object.Size = file.Size
	object.Version = version.String()

	_, err = query.Clone().Save(&object)
	if err != nil {
		return *response.SystemError(err)
	}
	newName, err := file.Save(savePath, true)
	if err != nil {
		return *response.SystemError(err)
	}

	realPath := path.Join(savePath, newName)

	s, err := gmd5.EncryptFile(realPath)
	if err != nil {
		return *response.SystemError(err)
	}
	if s != hash {
		query.Clone().Where("`id` = ?", object.Id).Delete()
		return *response.ErrorWithMessage(http.StatusBadRequest, "Inconsistent with expected MD5 value")
	}
	_, err = query.Clone().Data(g.Map{"real_path": path.Join(savePath, newName)}).Where("`id` = ?", object.Id).Update()
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
	query := g.DB().Model(&models.ObjectInfo{}).Where(g.Map{"bucket": bucket, "key": self.Key, "name": self.Name})

	if c, err := query.Clone().Count(); err != nil {
		return *response.SystemError(err)
	} else {
		if c == 1 {
			if _, err := query.Clone().Delete(); err != nil {
				return *response.SystemError(err)
			}
		} else {
			return *response.ErrorWithMessage(http.StatusNotFound, "object does not exist")
		}
	}

	return *response.Success()
}
