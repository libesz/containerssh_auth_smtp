package config

import (
	"net/http"
	"path/filepath"
	"strings"

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

	volumeNames, err := c.mapping.GetUserVolumeNames(request.Username)
	if err != nil {
		return config, err
	}
	volumeNamePrefix := c.mapping.GetVolumePrefix()

	config.Docker.Execution.Launch.ContainerName = containerNamePrefix + "-" + strings.Replace(strings.Replace(request.Username, "@", "-", -1), ".", "-", -1) + "-" + generate(containerNameRandomTailSize)

	config.Docker.Execution.Launch.ContainerConfig.WorkingDir = contentMountPoint

	if config.Docker.Execution.Launch.HostConfig == nil {
		config.Docker.Execution.Launch.HostConfig = &container.HostConfig{}
	}
	if config.Docker.Execution.Launch.HostConfig.Mounts == nil {
		config.Docker.Execution.Launch.HostConfig.Mounts = []mount.Mount{}
	}
	for _, volumeName := range volumeNames {
		volumeNameWithoutPrefix := strings.Replace(volumeName, volumeNamePrefix, "", 1)
		config.Docker.Execution.Launch.HostConfig.Mounts = append(config.Docker.Execution.Launch.HostConfig.Mounts, mount.Mount{
			Type:          "volume",
			Source:        volumeName,
			Target:        filepath.Join(contentMountPoint, volumeNameWithoutPrefix),
			VolumeOptions: &mount.VolumeOptions{},
		})
	}

	return config, err
}
