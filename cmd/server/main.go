package main

import (
	"net/http"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/libesz/containerssh_smtp_auth/pkg/auth"
	"github.com/libesz/containerssh_smtp_auth/pkg/config"

	"github.com/containerssh/log"
)

type Config struct {
	ListenOn              string `envconfig:"LISTEN_ON" required:"true"`
	SmtpEp                string `envconfig:"SMTP_EP" required:"true"`
	SmtpServerName        string `envconfig:"SMTP_SERVER_NAME" required:"true"`
	UserVolumeMappingPath string `envconfig:"USER_VOLUME_MAPPING_PATH"`
}

func main() {
	logger, err := log.NewLogger(log.Config{Format: log.FormatText, Destination: log.DestinationStdout, Level: log.LevelDebug})
	if err != nil {
		panic(err.Error())
	}
	logger.Info("ContainerSSH SMTP authenticator started up")
	var envConfig Config
	err = envconfig.Process("", &envConfig)
	if err != nil {
		logger.Critical("Invalid environment configuration: ", err.Error())
		os.Exit(1)
	}
	logger.Info("Configuration: ", envConfig)

	var mapping config.MappingFileHandler
	if envConfig.UserVolumeMappingPath != "" {
		mapping = config.NewMappingFileHandler(logger)
		err = mapping.Load(envConfig.UserVolumeMappingPath)
		if err != nil {
			panic(err.Error())
		}
		configHandler, err := config.NewConfigReqHandler(logger, mapping)
		if err != nil {
			panic(err.Error())
		}

		http.Handle("/config", configHandler)
	} else {
		logger.Info("USER_VOLUME_MAPPING_PATH not defined. Running in auth-only mode.")
	}

	authHandler := auth.NewSmtpAuthHandler(logger, envConfig.SmtpEp, envConfig.SmtpServerName, mapping)
	http.Handle("/auth/", authHandler)

	err = http.ListenAndServe(envConfig.ListenOn, nil)
	if err != nil {
		panic(err.Error())
	}
	logger.Info("Exiting")
}
