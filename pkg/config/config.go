package config

import (
	"net/http"

	"crypto/rand"
	"unsafe"

	"github.com/containerssh/configuration/v2"
	"github.com/containerssh/log"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
)

var alphabet = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

const (
	contentMountPoint           = "/content"
	containerNamePrefix         = "cssh-client"
	containerNameRandomTailSize = 5
)

func generate(size int) string {
	b := make([]byte, size)
	rand.Read(b)
	for i := 0; i < size; i++ {
		b[i] = alphabet[b[i]/5]
	}
	return *(*string)(unsafe.Pointer(&b))
}

func NewConfigReqHandler(logger log.Logger, mapping MappingFileHandler) (http.Handler, error) {
	return configuration.NewHandler(&configReqHandler{logger, mapping}, logger)
}

type configReqHandler struct {
	logger  log.Logger
	mapping MappingFileHandler
}

func (c configReqHandler) OnConfig(
	request configuration.ConfigRequest,
) (config configuration.AppConfig, err error) {
	c.logger.Info("Config request for: ", request.Username, "Session ID:", request.SessionID)
	if config.Docker.Execution.Launch.ContainerConfig == nil {
		config.Docker.Execution.Launch.ContainerConfig = &container.Config{}
	}

	volumeName, err := c.mapping.GetUserVolumeName(request.Username)
	if err != nil {
		return config, err
	}

	config.Docker.Execution.Launch.ContainerName = containerNamePrefix + "-" + volumeName + "-" + generate(containerNameRandomTailSize)

	config.Docker.Execution.Launch.ContainerConfig.WorkingDir = contentMountPoint

	if config.Docker.Execution.Launch.HostConfig == nil {
		config.Docker.Execution.Launch.HostConfig = &container.HostConfig{}
	}
	if config.Docker.Execution.Launch.HostConfig.Mounts == nil {
		config.Docker.Execution.Launch.HostConfig.Mounts = []mount.Mount{}
	}
	config.Docker.Execution.Launch.HostConfig.Mounts = append(config.Docker.Execution.Launch.HostConfig.Mounts, mount.Mount{
		Type:          "volume",
		Source:        volumeName,
		Target:        contentMountPoint,
		VolumeOptions: &mount.VolumeOptions{},
	})

	return config, err
}
