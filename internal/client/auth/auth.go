package auth

import (
	"context"
	"fmt"
	"github.com/a1exCross/chat-cli/internal/client"
	"github.com/a1exCross/chat-cli/internal/client/auth/proto/access_v1"
	"github.com/a1exCross/chat-cli/internal/client/auth/proto/auth_v1"
	"github.com/a1exCross/chat-cli/internal/client/auth/proto/user_v1"
	"github.com/a1exCross/chat-cli/internal/config"
	"github.com/a1exCross/chat-cli/internal/model"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type authService struct {
	authClient   auth_v1.AuthV1Client
	userClient   user_v1.UserV1Client
	accessClient access_v1.AccessV1Client
	config       config.AuthConfig
}

func NewAuthService(config config.AuthConfig) (client.AuthService, error) {
	conn, err := grpc.Dial(config.AuthAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("did not connect: %v", err)
	}

	auth := auth_v1.NewAuthV1Client(conn)
	user := user_v1.NewUserV1Client(conn)
	access := access_v1.NewAccessV1Client(conn)

	return &authService{
		authClient:   auth,
		accessClient: access,
		userClient:   user,
		config:       config,
	}, nil
}

func (a *authService) Login(ctx context.Context, userInfo model.UserInfo) (string, error) {
	res, err := a.authClient.Login(ctx, &auth_v1.LoginRequest{
		Username: userInfo.Username,
		Password: userInfo.Password,
	})

	if err != nil {
		return "", fmt.Errorf("login error: %v", err)
	}

	return res.GetRefreshToken(), err
}

func (a *authService) Authorize(ctx context.Context, refreshToken string) (string, error) {
	res, err := a.authClient.GetAccessToken(ctx, &auth_v1.GetAccessTokenRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		return "", fmt.Errorf("authorize error: %v", err)
	}

	return res.GetAccessToken(), nil
}

func (a *authService) CreateUser(ctx context.Context, params model.UserCreateParams) (int64, error) {
	resp, err := a.userClient.Create(ctx, &user_v1.CreateRequest{
		Info: &user_v1.UserInfo{
			Username: params.Username,
			Name:     params.Name,
			Role:     user_v1.UserRole(params.Role),
			Email:    params.Email,
		},
		Pass: &user_v1.UserPassword{
			Password:        params.Password,
			PasswordConfirm: params.PasswordConfirm,
		},
	})
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %v", err)
	}

	return resp.GetId(), nil
}

func (a *authService) DeleteUser(ctx context.Context, userID int64) error {
	_, err := a.userClient.Delete(ctx, &user_v1.DeleteRequest{
		Id: userID,
	})
	if err != nil {
		return fmt.Errorf("delete user error: %v", err)
	}

	return nil
}

func (a *authService) GetAccessToken(ctx context.Context, refreshToken string) (string, error) {
	res, err := a.authClient.GetAccessToken(ctx, &auth_v1.GetAccessTokenRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		return "", fmt.Errorf("get access token error: %v", err)
	}

	return res.GetAccessToken(), nil
}

func (a *authService) GetRefreshToken(ctx context.Context, oldRefreshToken string) (string, error) {
	res, err := a.authClient.GetRefreshToken(ctx, &auth_v1.GetRefreshTokenRequest{
		OldRefreshToken: oldRefreshToken,
	})
	if err != nil {
		return "", fmt.Errorf("get refresh token error: %v", err)
	}

	return res.GetRefreshToken(), nil
}
