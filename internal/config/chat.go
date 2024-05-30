package config

import (
	"errors"
	"net"
	"os"
)

const (
	ChatPortEnv = "CHAT_PORT"
	ChatHostEnv = "CHAT_HOST"
)

type chatConfig struct {
	port string
	host string
}

func NewChatConfig() (ChatConfig, error) {
	port := os.Getenv(ChatPortEnv)
	if len(port) == 0 {
		return nil, errors.New("environment variable 'CHAT_PORT' is not set")
	}

	host := os.Getenv(ChatHostEnv)
	if len(host) == 0 {
		return nil, errors.New("environment variable 'CHAT_HOST' is not set")
	}

	return &chatConfig{
		port: port,
		host: host,
	}, nil
}

func (a *chatConfig) ChatAddress() string {
	return net.JoinHostPort(a.host, a.port)
}
