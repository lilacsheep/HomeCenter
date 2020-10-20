package models

import (
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/util/gutil"
	"homeproxy/app/server/download"
)

const (
	DownloadListTable = "proxy_download_list"
)

var Tasks *gmap.TreeMap

func init() {
	Tasks = gmap.NewTreeMap(gutil.ComparatorString, true)
}

// TODO: download tasks manager
//func InitDownloadTasks()  {
//	tasks, err := GetALLDownloadTasks()
//	if err != nil {
//		glog.Errorf("get download tasks error: %s", err.Error())
//	} else {
//
//	}
//}

func GetALLDownloadTasks() (tasks []download.Task, err error) {
	c, _ := DB.Collection(DownloadListTable)
	err = c.All(&tasks)
	return
}
