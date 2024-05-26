package user

import (
	"context"
	"dating-app/domain"
	"dating-app/domain/common/response"
	"dating-app/domain/user"
	userMocks "dating-app/domain/user/mocks"
	utilMocks "dating-app/infrastructure/utils/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type userTestParam struct {
	userRepository      *userMocks.UserRepository
	userCacheRepository *userMocks.UserCacheRepository
	userUsecase         user.UserUsecase
	context             context.Context
	hashUtils           *utilMocks.HashUtils
}

func setupUserTestParam(t *testing.T) userTestParam {
	userRepository := userMocks.NewUserRepository(t)
	userCacheRepository := userMocks.NewUserCacheRepository(t)
	hashUtils := utilMocks.NewHashUtils(t)
	userUsecase := NewUserUsecase(domain.Repositories{
		UserRepository:      userRepository,
		UserCacheRepository: userCacheRepository,
	}, "", domain.Utils{
		HashUtils: hashUtils,
	})

	return userTestParam{
		userRepository:      userRepository,
		userCacheRepository: userCacheRepository,
		userUsecase:         userUsecase,
		context:             context.Background(),
		hashUtils:           hashUtils,
	}
}

func Test_RegisterUser(t *testing.T) {
	params := setupUserTestParam(t)

	t.Run("return error on GetUserCountByEmail error", func(t *testing.T) {
		params.userRepository.On("GetUserCountByEmail", mock.Anything, mock.Anything).Return(0, assert.AnError).Once()

		res := params.userUsecase.RegisterUser(params.context, user.UserRegistrationRequest{})
		require.Equal(t, true, res.Error)
	})

	t.Run("return error on email taken", func(t *testing.T) {
		params.userRepository.On("GetUserCountByEmail", mock.Anything, mock.Anything).Return(1, nil).Once()

		res := params.userUsecase.RegisterUser(params.context, user.UserRegistrationRequest{})
		require.Equal(t, true, res.Error)
		require.Equal(t, response.ERROR_EMAIL_TAKEN, *res.Message)
	})

	t.Run("return error on InsertUser error", func(t *testing.T) {
		params.userRepository.On("GetUserCountByEmail", mock.Anything, mock.Anything).Return(0, nil).Once()
		params.userRepository.On("InsertUser", mock.Anything, mock.Anything).Return(assert.AnError).Once()
		params.hashUtils.On("HashPassword", mock.Anything, mock.Anything).Return("").Once()

		res := params.userUsecase.RegisterUser(params.context, user.UserRegistrationRequest{})
		require.Equal(t, true, res.Error)
	})

	t.Run("return success", func(t *testing.T) {
		params.userRepository.On("GetUserCountByEmail", mock.Anything, mock.Anything).Return(0, nil).Once()
		params.userRepository.On("InsertUser", mock.Anything, mock.Anything).Return(nil).Once()
		params.hashUtils.On("HashPassword", mock.Anything, mock.Anything).Return("").Once()

		res := params.userUsecase.RegisterUser(params.context, user.UserRegistrationRequest{})
		require.Equal(t, false, res.Error)
	})
}
