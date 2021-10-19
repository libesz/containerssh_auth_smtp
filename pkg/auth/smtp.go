package auth

import (
	"crypto/tls"
	"errors"
	"net/http"
	"net/smtp"

	"github.com/containerssh/auth"
	"github.com/containerssh/log"
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

func NewSmtpAuthHandler(logger log.Logger, smtpEp string, smtpServerName string) http.Handler {
	return auth.NewHandler(smtpAuthHandler{
		logger:         logger,
		smtpEp:         smtpEp,
		smtpServerName: smtpServerName,
	}, logger)
}

type smtpAuthHandler struct {
	logger         log.Logger
	smtpEp         string
	smtpServerName string
}

func (h smtpAuthHandler) OnPassword(
	Username string,
	Password []byte,
	RemoteAddress string,
	ConnectionID string,
) (bool, error) {
	h.logger.Info("Login attempt with username:", Username, "Connection ID:", ConnectionID)

	client, err := smtp.Dial(h.smtpEp)
	if err != nil {
		return false, err
	}
	defer client.Close()
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

func (h smtpAuthHandler) OnPubKey(
	Username string,
	// PublicKey is the public key in the authorized key format.
	PublicKey string,
	RemoteAddress string,
	ConnectionID string,
) (bool, error) {
	return false, nil
}
