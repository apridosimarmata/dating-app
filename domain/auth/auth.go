package auth

import (
	"context"
	"dating-app/domain/common/response"
	"net/http"
)

type Tokens struct {
	Token        string  `json:"token"`
	RefreshToken *string `json:"refreshToken"`
}

type RefreshTokenRequest struct {
	RefreshToken *string `json:"refreshToken"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthUsecase interface {
	AuthenticateUser(ctx context.Context, request LoginRequest) (res *response.Response[Tokens])
	RefreshAccessToken(ctx context.Context, refreshToken string) (res *response.Response[Tokens])
	AuthorizeRequestMiddleware(next http.Handler) http.Handler
}
