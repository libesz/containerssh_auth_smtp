package main

import (
	"net/http"
	"os"

	"github.com/libesz/containerssh_smtp_auth/pkg/auth"
	"github.com/libesz/containerssh_smtp_auth/pkg/config"

	"github.com/containerssh/log"
)

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
	mapping := config.NewMappingFileHandler(logger)
	err = mapping.Load(userVolumeMappingPath)
	if err != nil {
		panic(err.Error())
	}

	authHandler := auth.NewSmtpAuthHandler(logger, smtpEP, smtpServerName, mapping)
	configHandler, err := config.NewConfigReqHandler(logger, mapping)
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
