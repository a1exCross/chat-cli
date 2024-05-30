package client

import (
	"context"
	"github.com/a1exCross/chat-cli/internal/client/chat/proto/chat_v1"
	"github.com/a1exCross/chat-cli/internal/model"
)

type AuthService interface {
	Login(ctx context.Context, userInfo model.UserInfo) (string, error)
	Authorize(ctx context.Context, refreshToken string) (string, error)
	CreateUser(ctx context.Context, params model.UserCreateParams) (int64, error)
	DeleteUser(ctx context.Context, userID int64) error
	GetAccessToken(ctx context.Context, refreshToken string) (string, error)
	GetRefreshToken(ctx context.Context, oldRefreshToken string) (string, error)
}

type ChatService interface {
	Create(ctx context.Context, usernames []string) (int64, error)
	Delete(ctx context.Context, chatID int64) error
	Connect(ctx context.Context, chatID, userID int64) (chat_v1.ChatV1_ConnectClient, error)
	SendMessage(ctx context.Context, params model.SendMessageParams) error
	ListChats(ctx context.Context, username string) ([]model.Chat, error)
}
