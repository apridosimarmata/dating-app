package auth

import (
	"context"
	"dating-app/domain"
	"dating-app/domain/auth"
	"dating-app/domain/common"
	"dating-app/domain/user"

	"dating-app/domain/common/response"
	userMocks "dating-app/domain/user/mocks"
	infraMocks "dating-app/infrastructure/mocks"
	utilMocks "dating-app/infrastructure/utils/mocks"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type authUsecaseTestParams struct {
	userRepository *userMocks.UserRepository
	_jwt           *infraMocks.JwtInterface
	authUsecase    auth.AuthUsecase
	context        context.Context
	hashUtils      *utilMocks.HashUtils
}

func setupAuthUsecaseTestParams(t *testing.T) authUsecaseTestParams {
	userRepository := userMocks.NewUserRepository(t)
	_jwt := infraMocks.NewJwtInterface(t)
	hashUtils := utilMocks.NewHashUtils(t)
	authUsecase := NewAuthUsecase(domain.Repositories{
		UserRepository: userRepository,
	}, common.Secret{}, _jwt, domain.Utils{
		HashUtils: hashUtils,
	})

	return authUsecaseTestParams{
		authUsecase:    authUsecase,
		_jwt:           _jwt,
		context:        context.Background(),
		userRepository: userRepository,
		hashUtils:      hashUtils,
	}
}

func Test_AuthenticateUser(t *testing.T) {
	params := setupAuthUsecaseTestParams(t)

	access_token := "access_token"
	refresh_token := "refresh_token"

	t.Run("return error on GetUserDetailsByEmail error", func(t *testing.T) {
		params.userRepository.On("GetUserDetailsByEmail", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()

		res := params.authUsecase.AuthenticateUser(params.context, auth.LoginRequest{})
		require.Equal(t, true, res.Error)
	})

	t.Run("return error on user found", func(t *testing.T) {
		params.userRepository.On("GetUserDetailsByEmail", mock.Anything, mock.Anything).Return(nil, nil).Once()

		res := params.authUsecase.AuthenticateUser(params.context, auth.LoginRequest{})
		require.Equal(t, true, res.Error)
		require.Equal(t, response.ERROR_USER_NOT_FOUND, *res.Message)
	})

	t.Run("return error on GenerateAccessToken error", func(t *testing.T) {
		params.userRepository.On("GetUserDetailsByEmail", mock.Anything, mock.Anything).Return(&user.User{}, nil).Once()
		params.hashUtils.On("HashPassword", mock.Anything, mock.Anything).Return("").Once()
		params._jwt.On("GenerateAccessToken", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()

		res := params.authUsecase.AuthenticateUser(params.context, auth.LoginRequest{})
		require.Equal(t, true, res.Error)
	})

	t.Run("return error on GenerateAccessToken [refresh token] error", func(t *testing.T) {
		params.userRepository.On("GetUserDetailsByEmail", mock.Anything, mock.Anything).Return(&user.User{}, nil).Once()
		params.hashUtils.On("HashPassword", mock.Anything, mock.Anything).Return("").Once()
		params._jwt.On("GenerateAccessToken", mock.Anything, mock.Anything, mock.Anything).Return(&access_token, nil).Once()
		params._jwt.On("GenerateAccessToken", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()

		res := params.authUsecase.AuthenticateUser(params.context, auth.LoginRequest{})
		require.Equal(t, true, res.Error)
	})

	t.Run("return success", func(t *testing.T) {
		params.userRepository.On("GetUserDetailsByEmail", mock.Anything, mock.Anything).Return(&user.User{}, nil).Once()
		params.hashUtils.On("HashPassword", mock.Anything, mock.Anything).Return("").Once()
		params._jwt.On("GenerateAccessToken", mock.Anything, mock.Anything, mock.Anything).Return(&access_token, nil).Once()
		params._jwt.On("GenerateAccessToken", mock.Anything, mock.Anything, mock.Anything).Return(&refresh_token, nil).Once()

		res := params.authUsecase.AuthenticateUser(params.context, auth.LoginRequest{})
		require.Equal(t, false, res.Error)
		require.Equal(t, "access_token", res.Data.Token)
		require.Equal(t, "refresh_token", *res.Data.RefreshToken)
	})
}
