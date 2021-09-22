package models

import (
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gmeta"
)

type DockerContanerOption struct {
	DefaultModel
	Name       string `json:"name"`
	Key        string `json:"key"`
	Value      string `json:"value"`
	gmeta.Meta `orm:"table:docker_container_option"`
}

type ContanerTemplate struct {
	Name          string
	Config        *container.Config         `json:"config"`
	HostConfig    *container.HostConfig     `json:"host_config"`
	NetworkConfig *network.NetworkingConfig `json:"network_config"`
}

func (self *ContanerTemplate) ReadFormOption(option DockerContanerOption) {
	switch option.Key {
	case "name":
		self.Name = option.Value
	case "environment":
		self.Config.Env = append(self.Config.Env, option.Value)
	case "network_mode":
		var nets = make(map[string]*network.EndpointSettings, 1)
		nets[option.Value] = &network.EndpointSettings{}
		self.NetworkConfig.EndpointsConfig = nets
	case "ports":
		// self.HostConfig.PortBindings
	case "command":
		self.Config.Cmd = append(self.Config.Cmd, option.Value)
	case "dns":
		self.HostConfig.DNS = append(self.HostConfig.DNS, option.Value)
	case "dns_search":
		self.HostConfig.DNSSearch = append(self.HostConfig.DNSSearch, option.Value)
	case "image":
		self.Config.Image = option.Value
	case "links":
		self.HostConfig.Links = append(self.HostConfig.Links, option.Value)
	case "volumes":
		self.HostConfig.VolumesFrom = append(self.HostConfig.VolumesFrom, option.Value)
	case "entrypoint":
		self.Config.Entrypoint = append(self.Config.Entrypoint, option.Value)
	case "extra_hosts":
		self.HostConfig.ExtraHosts = append(self.HostConfig.ExtraHosts, option.Value)
	case "labels":
		t := strings.SplitN(option.Value, ":", 1)
		self.Config.Labels[t[0]] = t[1]
	case "restart_policy":
		t := strings.SplitAfterN(option.Value, ":", 1)
		var (
			name string
			count int
		)
		if len(t) == 1 {
			name = t[0]
		} else {
			name = t[0]
			count = gconv.Int(t[1])
		}
		self.HostConfig.RestartPolicy = container.RestartPolicy{Name: name, MaximumRetryCount: count}
	}
}
