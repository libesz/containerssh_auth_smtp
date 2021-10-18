package main

import (
	"crypto/tls"
	"net/http"
	"os"

	"github.com/containerssh/auth"
	"github.com/containerssh/configuration/v2"
	"github.com/docker/docker/api/types/container"

	"github.com/containerssh/log"

	"errors"
	"net/smtp"
)

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unkown fromServer")
		}
	}
	return nil, nil
}

type myAuthHandler struct {
	logger         log.Logger
	smtpEp         string
	smtpServerName string
}

func (h *myAuthHandler) OnPassword(
	Username string,
	Password []byte,
	RemoteAddress string,
	ConnectionID string,
) (bool, error) {
	h.logger.Info("Login attempt with username:", Username, "Connection ID:", ConnectionID)
	//return false, nil

	client, err := smtp.Dial(h.smtpEp)
	defer client.Close()
	if err != nil {
		return false, err
	}
	err = client.StartTLS(&tls.Config{ServerName: h.smtpServerName})
	if err != nil {
		return false, err
	}

	err = client.Auth(LoginAuth(Username, string(Password)))
	if err != nil {
		h.logger.Info("Login failed with username:", Username, "Connection ID:", ConnectionID, "Error:", err)
		return false, nil
	}
	client.Close()
	h.logger.Info("Login succeeded with username:", Username, "Connection ID:", ConnectionID)
	return true, nil
}

func (h *myAuthHandler) OnPubKey(
	Username string,
	// PublicKey is the public key in the authorized key format.
	PublicKey string,
	RemoteAddress string,
	ConnectionID string,
) (bool, error) {
	return false, nil
}

type myConfigReqHandler struct {
	logger log.Logger
}

func (m *myConfigReqHandler) OnConfig(
	request configuration.ConfigRequest,
) (config configuration.AppConfig, err error) {
	if config.Docker.Execution.Launch.ContainerConfig == nil {
		config.Docker.Execution.Launch.ContainerConfig = &container.Config{}
	}
	//config.Docker.Execution.Launch.ContainerConfig.Image = "alpine"
	if config.Docker.Execution.Launch.ContainerConfig.Labels == nil {
		config.Docker.Execution.Launch.ContainerConfig.Labels = map[string]string{}
	}
	config.Docker.Execution.Launch.ContainerConfig.Labels["test"] = "test2"

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

	authHandler := auth.NewHandler(&myAuthHandler{logger, smtpEP, smtpServerName}, logger)
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
