package config

import (
	"net/http"

	"github.com/containerssh/configuration/v2"
	"github.com/containerssh/log"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
)

const (
	contentMountPoint = "/content"
)

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

	config.Docker.Execution.Launch.ContainerConfig.WorkingDir = contentMountPoint

	volumeName, err := c.mapping.GetUserVolumeName(request.Username)
	if err != nil {
		return config, err
	}

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
