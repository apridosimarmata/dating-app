package auth

import (
	"context"
	"dating-app/domain"
	"dating-app/domain/auth"
	"dating-app/domain/common"
	"dating-app/domain/common/response"
	"dating-app/domain/user"
	"dating-app/infrastructure"
	"dating-app/infrastructure/utils"
	"net/http"
	"strings"
)

type authUsecase struct {
	userRepository user.UserRepository
	passwordSalt   string
	jwtSecret      string
	_jwt           infrastructure.JwtInterface
	hashUtils      utils.HashUtils
	APIKey         string
}

func NewAuthUsecase(repositories domain.Repositories, secret common.Secret, jwt infrastructure.JwtInterface, utils domain.Utils) auth.AuthUsecase {
	return &authUsecase{
		userRepository: repositories.UserRepository,
		passwordSalt:   secret.PsasswordSalt,
		jwtSecret:      secret.JwtSecret,
		_jwt:           jwt,
		hashUtils:      utils.HashUtils,
		APIKey:         secret.APIKey,
	}
}

func (usecase *authUsecase) AuthenticateUser(ctx context.Context, request auth.LoginRequest) (res *response.Response[auth.Tokens]) {
	user, err := usecase.userRepository.GetUserDetailsByEmail(ctx, request.Email)
	if err != nil {
		infrastructure.Log("got error on usecase.userRepository.GetUserDetailsByEmail() - AuthenticateUser")
		return &response.Response[auth.Tokens]{
			Error:   true,
			Message: &response.ERROR_INTERNAL_SERVER_ERROR,
		}
	}

	if user == nil {
		return &response.Response[auth.Tokens]{
			Error:   true,
			Message: &response.ERROR_USER_NOT_FOUND,
		}
	}

	if usecase.hashUtils.HashPassword(request.Password, usecase.passwordSalt) != user.Password {
		return &response.Response[auth.Tokens]{
			Error:   true,
			Message: &response.ERROR_UNAUTHORIZED,
		}
	}

	accessToken, err := usecase._jwt.GenerateAccessToken(user.UID, usecase.jwtSecret, false)
	if err != nil {
		infrastructure.Log("got error on usecase.generateAccessToke() - AuthenticateUser")
		return &response.Response[auth.Tokens]{
			Error:   true,
			Message: &response.ERROR_INTERNAL_SERVER_ERROR,
		}
	}

	refreshToken, err := usecase._jwt.GenerateAccessToken(user.UID, usecase.jwtSecret, true)
	if err != nil {
		infrastructure.Log("got error on usecase.generateAccessToke() - AuthenticateUser")
		return &response.Response[auth.Tokens]{
			Error:   true,
			Message: &response.ERROR_INTERNAL_SERVER_ERROR,
		}
	}

	return &response.Response[auth.Tokens]{
		Error: false,
		Data: &auth.Tokens{
			Token:        *accessToken,
			RefreshToken: refreshToken,
		},
	}
}
func (usecase *authUsecase) RefreshAccessToken(ctx context.Context, refreshToken string) (res *response.Response[auth.Tokens]) {
	var _refreshToken *string

	userUid, ttl, err := usecase._jwt.ValidateToken(refreshToken, usecase.jwtSecret)
	if err != nil {
		return &response.Response[auth.Tokens]{
			Error:   true,
			Message: &response.ERROR_UNAUTHORIZED,
		}
	}

	accessToken, err := usecase._jwt.GenerateAccessToken(*userUid, usecase.jwtSecret, false)
	if err != nil {
		infrastructure.Log("got error on usecase.generateAccessToke() - AuthenticateUser")
		return &response.Response[auth.Tokens]{
			Error:   true,
			Message: &response.ERROR_INTERNAL_SERVER_ERROR,
		}
	}

	if (ttl / 86400) < 1 {
		_refreshToken, err = usecase._jwt.GenerateAccessToken(*userUid, usecase.jwtSecret, true)
		if err != nil {
			infrastructure.Log("got error on usecase.generateAccessToke() - AuthenticateUser")
			return &response.Response[auth.Tokens]{
				Error:   true,
				Message: &response.ERROR_INTERNAL_SERVER_ERROR,
			}
		}
	}

	return &response.Response[auth.Tokens]{
		Error: false,
		Data: &auth.Tokens{
			Token:        *accessToken,
			RefreshToken: _refreshToken,
		},
	}
}

func (usecase *authUsecase) AuthorizeRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			response := &response.Response[*interface{}]{
				Error:   true,
				Message: &response.ERROR_INTERNAL_SERVER_ERROR,
			}
			response.WriteResponse(w)
			return
		}

		userUid, _, err := usecase._jwt.ValidateToken(authHeader[1], usecase.jwtSecret)
		if err != nil {
			response := &response.Response[auth.Tokens]{
				Error:   true,
				Message: &response.ERROR_UNAUTHORIZED,
			}
			response.WriteResponse(w)
			return

		}

		ctx = context.WithValue(ctx, "user_uid", userUid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (usecase *authUsecase) AuthorizePartnerRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			response := &response.Response[*interface{}]{
				Error:   true,
				Message: &response.ERROR_INTERNAL_SERVER_ERROR,
			}
			response.WriteResponse(w)
			return
		}

		if authHeader[1] != usecase.APIKey {
			response := &response.Response[*interface{}]{
				Error:   true,
				Message: &response.ERROR_UNAUTHORIZED,
			}
			response.WriteResponse(w)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
