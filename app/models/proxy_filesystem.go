package models

import (
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/gogf/gf/os/gfile"
)


type FileInfo struct {
	Name        string     `json:"name"`
	Path        string     `json:"path"`
	Size        int64      `json:"size"`
	CreateAt    string     `json:"create_at"`
	IsDir       bool       `json:"is_dir"`
	Children    []FileInfo `json:"children"`
	HasChildren bool       `json:"has_children"`
}

type ProxyFileSystemNode struct {
	ID       int       `json:"id" storm:"id,increment"`
	Path     string    `json:"path" storm:"unique"`
	Name     string    `json:"name"`
	CreateAt time.Time `json:"create_at"`
}

func (self *ProxyFileSystemNode) walk(p string) (files []FileInfo) {
	if filepath.IsAbs(p) {
		infos, _ := ioutil.ReadDir(p)

		for _, file := range infos {
			if file.IsDir() {
				children := self.walk(gfile.Join(p, file.Name()))
				sub := FileInfo{
					Name:        file.Name(),
					Path:        gfile.Join(p, file.Name()),
					Size:        file.Size(),
					CreateAt:    file.ModTime().Format("2006-01-02 15:04:05"),
					IsDir:       true,
					HasChildren: len(children) != 0,
					Children:    children,
				}
				files = append(files, sub)
			} else {
				files = append(files, FileInfo{
					Name:        file.Name(),
					Path:        gfile.Join(p, file.Name()),
					Size:        file.Size(),
					CreateAt:    file.ModTime().Format("2006-01-02 15:04:05"),
					IsDir:       false,
					Children:    nil,
					HasChildren: false,
				})
			}
		}
	}
	return
}

func (self *ProxyFileSystemNode) Files() (files []FileInfo) {
	absPath := gfile.Abs(self.Path)
	return self.walk(absPath)
}
