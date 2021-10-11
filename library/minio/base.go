package minio

import (
	"context"
	"errors"
	"fmt"
	"homeproxy/library/docker"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

func CheckImage() (bool, error) {
	// containers, err := docker.Docker.ContainerList(context.Background(), types.ContainerListOptions{})
	// if err != nil {
	// 	return
	// }
	// for _, c := range containers {
	// 	c.Image
	// }
	return false, nil
}

func CreateMinioServer(target string, accessKey string, secretKey string, port int, public bool) (string, error) {
	if docker.Docker == nil {
		return "", errors.New("docker client is none, please check")
	}
	var (
		access = fmt.Sprintf("MINIO_ACCESS_KEY=%s", accessKey)
		secret = fmt.Sprintf("MINIO_SECRET_KEY=%s", secretKey)
		host = "127.0.0.1"
	)
	if target == "" {
		target = "/data/minio"
	}
	if public {
		host = "0.0.0.0"
	}
	if port == 0 {
		port = 9000
	}
	Config := container.Config{
		Image: "docker.io/minio/minio",
		Env: []string{access, secret},
		Cmd: []string{"server", "/data"},
	}
	HostConfig :=  container.HostConfig{
		PortBindings: nat.PortMap{nat.Port("9000/tcp"): []nat.PortBinding{{HostIP:host, HostPort: fmt.Sprintf("%d/tcp", port)}}},
		RestartPolicy: container.RestartPolicy{Name: "always"},
		Mounts: []mount.Mount{{Type: "bind", Source: "/data", Target: target}},
	}
	NetworkConfig := network.NetworkingConfig{}

	body, err := docker.Docker.ContainerCreate(context.Background(), &Config, &HostConfig, &NetworkConfig, nil, "minio")
	if err != nil {
		return "", err
	}
	
	return body.ID, nil
}