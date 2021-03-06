package requests

import (
	"fmt"
	"homeproxy/app/models"
	"homeproxy/library/filedb2"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gfile"
)

type AllFilesystemNodesRequest struct{}

func (self *AllFilesystemNodesRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	var (
		nodes []models.ProxyFileSystemNode
		err   error
	)
	err = filedb2.DB.All(&nodes)
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		response.SuccessWithDetail(nodes)
	}
	return
}

func NewAllFileSystemNodesRequest() *AllFilesystemNodesRequest {
	return &AllFilesystemNodesRequest{}
}

type AllFilesystemFilesRequest struct {
	ID int `json:"id"`
}

func (self *AllFilesystemFilesRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	var node models.ProxyFileSystemNode
	err := filedb2.DB.One("ID", self.ID, &node)
	if err != nil {
		response.ErrorWithMessage(http.StatusNotFound, err.Error())
	} else {
		files := node.Files()
		var (
			t    = g.Map{"files": files, "dirs": []g.Map{}}
			dirs []g.Map
		)

		dirs = append(dirs, g.Map{
			"name": node.Name,
			"id":   node.ID,
			"path": node.Path,
		})
		t["dirs"] = dirs
		response.SuccessWithDetail(t)
	}
	return
}

func NewAllFilesystemFilesRequest() *AllFilesystemFilesRequest {
	return &AllFilesystemFilesRequest{}
}

type GetFilesystemFileInfoRequest struct {
	Path string
}

func (self *GetFilesystemFileInfoRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	if !gfile.Exists(self.Path) {
		response.ErrorWithMessage(http.StatusNotFound, "文件不存在")
	} else {
		ext := gfile.Ext(self.Path)
		switch strings.ToLower(ext) {
		case ".png", ".jpg", ".jpeg", ".gif":
			file, _ := os.Open(self.Path)
			conf, _, _ := image.DecodeConfig(file)
			var info = g.Map{}
			info["height"] = conf.Height
			info["width"] = conf.Width
			info["type"] = "img"
			response.SuccessWithDetail(info)
		case ".mp4", ".webm":
			var info = g.Map{}
			info["type"] = "video"
			response.SuccessWithDetail(info)
		default:
			response.ErrorWithMessage(http.StatusInternalServerError, "不支持")
		}
	}

	return
}

func NewGetFilesystemFileInfoRequest() *GetFilesystemFileInfoRequest {
	return &GetFilesystemFileInfoRequest{}
}

type DownloadFilesystemFileRequest struct {
	Path   string `json:"path"`
	NodeID int    `json:"node_id"`
}

func (self *DownloadFilesystemFileRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	var node models.ProxyFileSystemNode
	err := filedb2.DB.One("ID", self.NodeID, &node)

	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
		return
	}
	if !strings.HasPrefix(self.Path, gfile.Abs(node.Path)) {
		response.ErrorWithMessage(http.StatusInternalServerError, "非法请求")
		return
	}
	if gfile.Exists(self.Path) {
		ext := gfile.Ext(self.Path)
		name := gfile.Basename(self.Path)
		switch strings.ToLower(ext) {
		case ".mp4":
			r.Response.Header().Set("Content-Type", "video/mpeg4")
		case ".webm":
			r.Response.Header().Set("Content-Type", "video/webm")
		case ".jpg", ".jpeg":
			r.Response.Header().Set("Content-Type", "image/jpeg")
		case ".png":
			r.Response.Header().Set("Content-Type", "image/png")
		case ".gif":
			r.Response.Header().Set("Content-Type", "image/gif")
		default:
			r.Response.Header().Set("Content-Type", "application/octet-stream")
		}
		r.Response.Header().Set("Content-Disposition", fmt.Sprintf(`attachment;filename="%s"`, name))
		r.Response.ServeFile(self.Path)
	} else {
		response.ErrorWithMessage(http.StatusNotFound, "资源不存在")
	}
	return
}

type RemoveFilesystemFileRequest struct {
	NodeID int    `json:"node_id"`
	Path   string `json:"path"`
}

func (self *RemoveFilesystemFileRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	var node models.ProxyFileSystemNode
	err := filedb2.DB.One("ID", self.NodeID, &node)
	if err != nil {
		response.ErrorWithMessage(http.StatusNotFound, err.Error())
	} else {
		err = gfile.Remove(self.Path)
		if err != nil {
			response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
		} else {
			response.Success()
		}
	}
	return
}

func NewRemoveFilesystemFileRequest() *RemoveFilesystemFileRequest {
	return &RemoveFilesystemFileRequest{}
}

type CreateFilesystemNodeRequest struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

func (self *CreateFilesystemNodeRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	node := models.ProxyFileSystemNode{Path: gfile.Abs(self.Path), Name: self.Name, CreateAt: time.Now()}
	if err := filedb2.DB.Save(&node); err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		response.SuccessWithDetail(node.ID)
	}
	return
}

func NewCreateFilesystemNodeRequest() *CreateFilesystemNodeRequest {
	return &CreateFilesystemNodeRequest{}
}

type RemoveFilesystemNodeRequest struct {
	ID int `json:"id"`
}

func (self *RemoveFilesystemNodeRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	err := filedb2.DB.DeleteStruct(&models.ProxyFileSystemNode{ID: self.ID})
	if err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		response.Success()
	}
	return
}

func NewRemoveFilesystemNodeRequest() *RemoveFilesystemNodeRequest {
	return &RemoveFilesystemNodeRequest{}
}

type CreateFilesystemDirRequest struct {
	NodeID int    `json:"node_id"`
	Path   string `json:"path"`
}

func (self *CreateFilesystemDirRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	var node models.ProxyFileSystemNode
	if err := filedb2.DB.One("ID", self.NodeID, &node); err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		path := gfile.Join(node.Path, self.Path)
		if gfile.Exists(path) {
			response.ErrorWithMessage(http.StatusInternalServerError, "dir already exist")
		} else {
			err = gfile.Mkdir(path)
			if err != nil {
				response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
			} else {
				response.Success()
			}
		}
	}
	return
}

func NewCreateFilesystemDirRequest() *CreateFilesystemDirRequest {
	return &CreateFilesystemDirRequest{}
}

type RemoveFilesystemDirRequest struct {
	NodeID int    `json:"node_id"`
	Path   string `json:"path"`
}

func (self *RemoveFilesystemDirRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	var node models.ProxyFileSystemNode
	if err := filedb2.DB.One("ID", self.NodeID, &node); err != nil {
		response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
	} else {
		path := gfile.Join(node.Path, self.Path)
		if !gfile.Exists(path) {
			response.ErrorWithMessage(http.StatusInternalServerError, "dir not exist")
		} else {
			err := gfile.Remove(path)
			if err != nil {
				response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
			} else {
				response.Success()
			}
		}
	}
	return
}

func NewRemoveFilesystemDirRequest() *RemoveFilesystemDirRequest {
	return &RemoveFilesystemDirRequest{}
}

type UploadFilesystemFileRequest struct{}

func (self *UploadFilesystemFileRequest) Exec(r *ghttp.Request) (response MessageResponse) {
	file := r.GetUploadFile("file")
	path := r.GetString("path")
	nodeID := r.GetString("node_id")
	if nodeID == "" && path == "" {
		response.ErrorWithMessage(http.StatusInternalServerError, "未选择节点")
	} else {
		if path != "" {
			if gfile.Exists(path) {
				file.Save(path)
				response.Success()
			} else {
				response.ErrorWithMessage(http.StatusInternalServerError, "路径不存在")
			}
		} else {
			var node models.ProxyFileSystemNode
			if err := filedb2.DB.One("ID", nodeID, &node); err != nil {
				response.ErrorWithMessage(http.StatusInternalServerError, err.Error())
			} else {
				file.Save(node.Path)
				response.Success()
			}
		}
	}

	return
}

func NewUploadFilesystemFileRequest() *UploadFilesystemFileRequest {
	return &UploadFilesystemFileRequest{}
}
