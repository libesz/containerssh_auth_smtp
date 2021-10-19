package config

import (
	"net/http"

	"github.com/containerssh/configuration/v2"
	"github.com/containerssh/log"
	"github.com/docker/docker/api/types/container"
)

func NewConfigReqHandler(logger log.Logger) (http.Handler, error) {
	return configuration.NewHandler(&configReqHandler{logger}, logger)
}

type configReqHandler struct {
	logger log.Logger
}

func (c configReqHandler) OnConfig(
	request configuration.ConfigRequest,
) (config configuration.AppConfig, err error) {
	c.logger.Info("Config request for: ", request.Username, "Session ID:", request.SessionID)
	if config.Docker.Execution.Launch.ContainerConfig == nil {
		config.Docker.Execution.Launch.ContainerConfig = &container.Config{}
	}

	config.Docker.Execution.Launch.ContainerConfig.WorkingDir = "/root"

	return config, err
}
