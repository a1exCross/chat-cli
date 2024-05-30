package app

import (
	"github.com/a1exCross/chat-cli/internal/client"
	"github.com/a1exCross/chat-cli/internal/client/auth"
	"github.com/a1exCross/chat-cli/internal/client/chat"
	"github.com/a1exCross/chat-cli/internal/config"
	"log"
)

type ServiceProvider struct {
	authConfig config.AuthConfig
	chatConfig config.ChatConfig

	authService client.AuthService
	chatService client.ChatService
}

func NewServiceProvider() *ServiceProvider {
	return &ServiceProvider{}
}

func (s *ServiceProvider) getAuthConfig() config.AuthConfig {
	if s.authConfig == nil {
		authConfig, err := config.NewAuthConfig()
		if err != nil {
			log.Fatalf("failed to get auth config")
		}

		s.authConfig = authConfig
	}

	return s.authConfig
}

func (s *ServiceProvider) getChatConfig() config.ChatConfig {
	if s.chatConfig == nil {
		chatConfig, err := config.NewChatConfig()
		if err != nil {
			log.Fatalf("failed to get chat config")
		}

		s.chatConfig = chatConfig
	}

	return s.chatConfig
}

func (s *ServiceProvider) GetAuthService() client.AuthService {
	if s.authService == nil {
		serv, err := auth.NewAuthService(s.getAuthConfig())
		if err != nil {
			log.Fatalf("failed to get auth service %s", err.Error())
		}

		s.authService = serv
	}

	return s.authService
}

func (s *ServiceProvider) GetChatService() client.ChatService {
	if s.chatConfig == nil {
		serv, err := chat.NewChatServer(s.getChatConfig())
		if err != nil {
			log.Fatalf("failed to get chat service %s", err.Error())
		}

		s.chatService = serv
	}

	return s.chatService
}
