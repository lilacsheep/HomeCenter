package models

import (
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/util/grand"
)

/*
-v 本地文件夹1:/downloads Aria2下载位置
-v 本地文件夹2:/config Aria2配置文件位置
-e PUID=1026 Linux用户UID
-e PGID=100 Linux用户GID
-e SECRET=yourtoken Aria2 token
-e CACHE=1024M Aria2磁盘缓存配置
-e PORT=6800 RPC通讯端口
-e WEBUI=true 启用WEBUI
-e WEBUI_PORT=8080 WEBUI端口
-e BTPORT=32516 DHT和BT监听端口
-e UT=true 启动容器时更新trackers
-e CTU= 启动容器时更新自定义trackes地址 https://cdn.jsdelivr.net/gh/XIU2/TrackersListCollection@master/best_aria2.txt
-e RUT=true	每天凌晨3点更新trackers
-e SMD=true 保存磁力链接为种子文件
-e FA= 磁盘预分配模式none,falloc,trunc,prealloc
-p 6800:6800 Aria2 RPC连接端口
-p 6881:6881 Aria2 tcp下载端口
-p 6881:6881/udp Aria2 p2p udp下载端口
--restart unless-stopped
*/

type DownloadSettings struct {
	Port         int    `json:"port"`          // RPC连接端口
	TcpPort      int    `json:"tcp_port"`      // tcp下载端口
	UdpPort      int    `json:"udp_port"`      // Aria2 p2p udp下载端口
	DownloadPath string `json:"download_path"` // Aria2下载位置
	ConfigPath   string `json:"config_path"`   // Aria2配置文件位置
	SECRET       string `json:"token"`         // Aria2 token
	WebUi        bool   `json:"webui"`         // 启用WEBUI
	WebUiPort    int    `json:"webui_port"`    // WEBUI端口
	BTPort       int    `json:"bt_port"`       // DHT和BT监听端口
	UT           bool   `json:"ut"`            // 启动容器时更新trackers
	CTU          string `json:"ctu"`           // 启动容器时更新自定义trackes地址
	RUT          bool   `json:"rut"`           // 每天凌晨3点更新trackers
	SMD          bool   `json:"smd"`           // 保存磁力链接为种子文件
	FA           string `json:"fa"`            // 磁盘预分配模式none,falloc,trunc,prealloc
	AutoClean    int    `json:"auto_clean"`    // 自动清理Bt下载后文件夹内内容，根据文件大小判断
	AutoStart    bool   `json:"auto_start"`    // 启动Aria2任务
	PublicVisit  bool   `json:"public_visit"`  // 允许公共访问
	ContainerId  string `json:"container_id"`
}

func (s *DownloadSettings) Env() []string {
	return []string{
		fmt.Sprintf("SECRET=%s", s.SECRET),
		"CACHE=1024M",
		fmt.Sprintf("PORT=%d", s.Port),
		fmt.Sprintf("WEBUI=%v", s.WebUi),
		fmt.Sprintf("WEBUI_PORT=%d", s.WebUiPort),
		fmt.Sprintf("BTPORT=%d", s.BTPort),
		fmt.Sprintf("UT=%v", s.UT),
		fmt.Sprintf("CTU=%v", s.CTU),
		fmt.Sprintf("RUT=%v", s.RUT),
		fmt.Sprintf("SMD=%v", s.SMD),
		fmt.Sprintf("FA=%v", s.FA),
	}
}

func (s *DownloadSettings) ContainerHostConfig() *container.HostConfig {
	var host string

	if s.PublicVisit {
		host = "0.0.0.0"
	} else {
		host = "127.0.0.1"
	}

	portMap := nat.PortMap{
		nat.Port("6800/tcp"):  []nat.PortBinding{{HostIP: host, HostPort: fmt.Sprintf("%d/tcp", s.Port)}},
		nat.Port("6881/tcp"):  []nat.PortBinding{{HostIP: host, HostPort: fmt.Sprintf("%d/tcp", s.TcpPort)}},
		nat.Port("32516/tcp"): []nat.PortBinding{{HostIP: host, HostPort: fmt.Sprintf("%d/tcp", s.BTPort)}},
		nat.Port("6881/udp"):  []nat.PortBinding{{HostIP: host, HostPort: fmt.Sprintf("%d/udp", s.UdpPort)}},
	}
	if s.WebUi {
		portMap[nat.Port("8080/tcp")] = []nat.PortBinding{{HostIP: host, HostPort: fmt.Sprintf("%d/tcp", s.WebUiPort)}}
	}

	downloadPath := gfile.Abs(s.DownloadPath)
	configPath := gfile.Abs(s.ConfigPath)
	if !gfile.Exists(downloadPath) {
		gfile.Mkdir(downloadPath)
	}

	if !gfile.Exists(configPath) {
		gfile.Mkdir(configPath)
	}
	mt := []mount.Mount{
		{Type: "bind", Source: downloadPath, Target: "/downloads"},
		{Type: "bind", Source: configPath, Target: "/config"},
	}
	HostConfig := container.HostConfig{
		PortBindings:  portMap,
		RestartPolicy: container.RestartPolicy{Name: "always"},
		Mounts:        mt,
	}
	return &HostConfig
}
func GetAria2Settings() (*DownloadSettings, error) {
	query := g.DB().Model(&GlobalConfig{})
	c, err := query.Where("`group` = ?", "aria2").Count()
	if err != nil {
		return nil, err
	}
	if c == 0 {
		query.Data(g.List{
			{"group": "aria2", "key": "port", "type": "int", "value": "6800", "desc": "RPC连接端口"},
			{"group": "aria2", "key": "tcp_port", "type": "int", "value": "6881", "desc": "tcp下载端口"},
			{"group": "aria2", "key": "udp_port", "type": "int", "value": "6881", "desc": "p2p udp下载端口"},
			{"group": "aria2", "key": "download_path", "type": "string", "value": "download", "desc": "Aria2下载位置"},
			{"group": "aria2", "key": "config_path", "type": "string", "value": "config/settings.config", "desc": "Aria2配置文件位置"},
			{"group": "aria2", "key": "token", "type": "string", "value": grand.S(16), "desc": "Aria2 token"},
			{"group": "aria2", "key": "webui", "type": "bool", "value": "false", "desc": "启用WEBUI"},
			{"group": "aria2", "key": "webui_port", "type": "int", "value": "8080", "desc": "启用WEBUI"},
			{"group": "aria2", "key": "bt_port", "type": "int", "value": "32516", "desc": "DHT和BT监听端口"},
			{"group": "aria2", "key": "ut", "type": "bool", "value": "true", "desc": "启动容器时更新trackers"},
			{"group": "aria2", "key": "ctu", "type": "string", "value": "https://cdn.jsdelivr.net/gh/XIU2/TrackersListCollection@master/best_aria2.txt", "desc": "启动容器时更新自定义trackes地址"},
			{"group": "aria2", "key": "rut", "type": "bool", "value": "true", "desc": "启动容器时更新自定义trackes地址"},
			{"group": "aria2", "key": "smd", "type": "string", "value": "none", "desc": "保存磁力链接为种子文件"},
			{"group": "aria2", "key": "fa", "type": "bool", "value": "true", "desc": "磁盘预分配模式none,falloc,trunc,prealloc"},
			{"group": "aria2", "key": "public_visit", "type": "bool", "value": "true", "desc": "允许公共访问"},
			{"group": "aria2", "key": "auto_start", "type": "bool", "value": "false", "desc": "启动Aria2任务"},
			{"group": "aria2", "key": "auto_clean", "type": "int", "value": "0", "desc": "自动清理Bt下载后文件夹内内容，根据文件大小判断"},
			{"group": "aria2", "key": "container_id", "type": "string", "value": "", "desc": "Aria2容器ID"},
		}).Save()
	}
	configs, err := GetConfigsMap("aria2")
	if err != nil {
		return nil, err
	}
	settings := &DownloadSettings{}
	err = gjson.New(configs).Struct(settings)
	return settings, err
}
