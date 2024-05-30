package chat

import (
	"context"
	"fmt"
	"github.com/a1exCross/chat-cli/internal/client"
	"github.com/a1exCross/chat-cli/internal/client/chat/converter"
	"github.com/a1exCross/chat-cli/internal/client/chat/proto/chat_v1"
	"github.com/a1exCross/chat-cli/internal/config"
	"github.com/a1exCross/chat-cli/internal/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ChatService struct {
	client chat_v1.ChatV1Client
	config config.ChatConfig
}

func NewChatServer(config config.ChatConfig) (client.ChatService, error) {
	conn, err := grpc.Dial(config.ChatAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to conect chat server: %w", err)
	}

	cl := chat_v1.NewChatV1Client(conn)

	return &ChatService{
		client: cl,
	}, nil
}

func (s *ChatService) Create(ctx context.Context, usernames []string) (int64, error) {
	req, err := s.client.Create(ctx, &chat_v1.CreateRequest{
		Usernames: usernames,
	})
	if err != nil {
		return 0, err
	}

	return req.GetId(), nil
}

func (s *ChatService) Delete(ctx context.Context, chatID int64) error {
	_, err := s.client.Delete(ctx, &chat_v1.DeleteRequest{
		Id: chatID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete chat: %w", err)
	}

	return nil
}

func (s *ChatService) Connect(ctx context.Context, chatID, userID int64) (chat_v1.ChatV1_ConnectClient, error) {
	res, err := s.client.Connect(ctx, &chat_v1.ConnectRequest{
		Id:     chatID,
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *ChatService) SendMessage(ctx context.Context, params model.SendMessageParams) error {
	_, err := s.client.SendMessage(ctx, &chat_v1.SendMessageRequest{
		Id: params.ChatID,
		Message: &chat_v1.Message{
			From:      params.From,
			Text:      params.Text,
			Timestamp: timestamppb.New(params.CreatedAt),
		},
	})

	return err
}

func (s *ChatService) ListChats(ctx context.Context, username string) ([]model.Chat, error) {
	res, err := s.client.ListChats(ctx, &chat_v1.ListRequest{
		Username: username,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list chats: %w", err)
	}

	return converter.ProtoToChats(res.GetChats()), nil
}
