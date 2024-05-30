package model

import "github.com/dgrijalva/jwt-go"

type UserCreateParams struct {
	Username        string
	Password        string
	PasswordConfirm string
	Email           string
	Role            int8
	Name            string
}

type UserInfoConfig struct {
	Username     string `json:"username"`
	UserID       int64  `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserInfo struct {
	Username string
	Password string
}

type UserClaims struct {
	jwt.StandardClaims
	UserID   int64
	Username string
	Role     int8
}
