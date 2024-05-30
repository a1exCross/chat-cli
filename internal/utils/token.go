package utils

import (
	"context"
	"errors"
	"github.com/a1exCross/chat-cli/internal/model"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/metadata"
	"time"
)

const authPrefix = "Bearer "

func NewOutgoingContextWithToken(ctx context.Context, accessToken string) context.Context {
	md := metadata.New(map[string]string{"Authorization": authPrefix + accessToken})

	return metadata.NewOutgoingContext(ctx, md)
}

func GetUserIDFromAccessToken(accessToken string) (int64, error) {
	claims, err := getUserClaims(accessToken)
	if err != nil {
		return -1, err
	}

	return claims.UserID, nil
}

func getUserClaims(accessToken string) (*model.UserClaims, error) {
	var parser jwt.Parser
	token, _, err := parser.ParseUnverified(accessToken, &model.UserClaims{})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*model.UserClaims)
	if !ok {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func TokenHasExpired(accessToken string) bool {
	claims, err := getUserClaims(accessToken)
	if err != nil {
		return true
	}

	if claims.ExpiresAt != 0 && claims.ExpiresAt > time.Now().Unix() {
		return false
	}

	return true
}
