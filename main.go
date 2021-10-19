package main

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/containerssh/configuration/v2"
	"github.com/docker/docker/api/types/container"
	"github.com/libesz/containerssh_smtp_auth/pkg/auth"

	"github.com/containerssh/log"
)

type myConfigReqHandler struct {
	logger log.Logger
}

func (m *myConfigReqHandler) OnConfig(
	request configuration.ConfigRequest,
) (config configuration.AppConfig, err error) {
	m.logger.Info("Config request for: ", request.Username, "Session ID:", request.SessionID)
	if config.Docker.Execution.Launch.ContainerConfig == nil {
		config.Docker.Execution.Launch.ContainerConfig = &container.Config{}
	}

	config.Docker.Execution.Launch.ContainerConfig.WorkingDir = "/root"

	return config, err
}

func main() {
	logger, err := log.NewLogger(log.Config{Format: log.FormatText, Destination: log.DestinationStdout, Level: log.LevelDebug})
	if err != nil {
		panic(err.Error())
	}
	logger.Info("ContainerSSH SMTP authenticator started up")
	listenOn := os.Getenv("LISTEN_ON")
	if listenOn == "" {
		panic("LISTEN_ON not defined")
	}
	smtpEP := os.Getenv("SMTP_EP")
	if smtpEP == "" {
		panic("SMTP_EP not defined")
	}
	smtpServerName := os.Getenv("SMTP_SERVER_NAME")
	if smtpServerName == "" {
		panic("SMTP_SERVER_NAME not defined")
	}

	userVolumeMappingPath := os.Getenv("USER_VOLUME_MAPPING_PATH")
	if userVolumeMappingPath == "" {
		panic("USER_VOLUME_MAPPING_PATH not defined")
	}
	mappingRawContent, err := ioutil.ReadFile(userVolumeMappingPath)
	if err != nil {
		panic("Mapping file read error: " + err.Error())
	}
	logger.Info(mappingRawContent)

	authHandler := auth.NewSmtpAuthHandler(logger, smtpEP, smtpServerName)
	configHandler, err := configuration.NewHandler(&myConfigReqHandler{logger}, logger)
	if err != nil {
		panic(err.Error())
	}

	http.Handle("/auth/", authHandler)
	http.Handle("/config", configHandler)

	err = http.ListenAndServe(listenOn, nil)
	if err != nil {
		panic(err.Error())
	}
	logger.Info("Exiting")
}
