package config

import (
	"net/http"

	"github.com/containerssh/configuration/v2"
	"github.com/containerssh/log"
	"github.com/docker/docker/api/types/container"
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

	config.Docker.Execution.Launch.ContainerConfig.WorkingDir = "/content"

	volumeName, err := c.mapping.GetUserVolumeName(request.Username)
	if err != nil {
		return config, err
	}
	if config.Docker.Execution.Launch.ContainerConfig.Volumes == nil {
		config.Docker.Execution.Launch.ContainerConfig.Volumes = map[string]struct{}{}
	}
	config.Docker.Execution.Launch.ContainerConfig.Volumes[volumeName+":/content"] = struct{}{}
	return config, err
}
