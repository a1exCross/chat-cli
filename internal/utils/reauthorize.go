package utils

import (
	"context"
	"fmt"
	"github.com/a1exCross/chat-cli/internal/app"
	"github.com/a1exCross/chat-cli/internal/config"
	"github.com/a1exCross/chat-cli/internal/model"
	"time"
)

const AccessDeniedError = "access denied"

func WithAuthInterceptor(ctx context.Context, delegate func(ctx context.Context) error) {
	err := delegate(ctx)
	if err != nil && err.Error() == AccessDeniedError {
		err = TryReauthorize(app.NewServiceProvider(), config.LoadExecConfig())
		if err != nil {
			fmt.Println("reauthorize failed:", err)
		}

		err = delegate(ctx)
	}
}

func TryReauthorize(di *app.ServiceProvider, cfg model.UserInfoConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	accessToken, err := di.GetAuthService().GetAccessToken(ctx, cfg.RefreshToken)
	if err != nil {
		return fmt.Errorf("get access token: %w", err)
	}

	refreshToken, err := di.GetAuthService().GetRefreshToken(ctx, cfg.RefreshToken)
	if err != nil {
		return fmt.Errorf("get refresh token: %w", err)
	}

	config.SaveExecConfig(model.UserInfoConfig{
		UserID:       cfg.UserID,
		Username:     cfg.Username,
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	})

	return nil
}
