package docker

import (
	"github.com/gogf/gf/os/glog"
	"github.com/moby/moby/client"
)

var (
	// Docker 全局客户端
	Docker *client.Client
)

// InitDockerClient 初始化docker 客户端
func InitDockerClient() {
	var err error
	Docker, err = client.NewClient(client.DefaultDockerHost, "", nil, nil)
	if err != nil {
		glog.Warningf("加载docker客户端失败: %s", err.Error())
	}
}
