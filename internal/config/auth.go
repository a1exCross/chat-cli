package config

import (
	"errors"
	"net"
	"os"
)

const (
	AuthPortEnv = "AUTH_PORT"
	AuthHostEnv = "AUTH_HOST"
)

type authConfig struct {
	port string
	host string
}

func NewAuthConfig() (AuthConfig, error) {
	port := os.Getenv(AuthPortEnv)
	if len(port) == 0 {
		return nil, errors.New("environment variable 'AUTH_PORT' is not set")
	}

	host := os.Getenv(AuthHostEnv)
	if len(host) == 0 {
		return nil, errors.New("environment variable 'AUTH_HOST' is not set")
	}

	return &authConfig{
		port: port,
		host: host,
	}, nil
}

func (a *authConfig) AuthAddress() string {
	return net.JoinHostPort(a.host, a.port)
}
