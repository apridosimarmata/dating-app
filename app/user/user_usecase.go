package user

import (
	"context"
	"dating-app/domain"
	"dating-app/domain/common/response"
	"dating-app/domain/feed"
	"dating-app/domain/user"
	"dating-app/infrastructure"
	"time"

	"dating-app/infrastructure/utils"

	"github.com/google/uuid"
)

type userUsecase struct {
	userRepository      user.UserRepository
	feedCacheRepository feed.FeedCacheRepository

	passwordSalt string
	hashUtils    utils.HashUtils
}

func NewUserUsecase(repositories domain.Repositories, passwordSalt string, utils domain.Utils) user.UserUsecase {
	return &userUsecase{
		userRepository:      repositories.UserRepository,
		feedCacheRepository: repositories.FeedCacheRepository,
		passwordSalt:        passwordSalt,
		hashUtils:           utils.HashUtils,
	}
}

func (usecase *userUsecase) RegisterUser(ctx context.Context, request user.UserRegistrationRequest) (res *response.Response[interface{}]) {
	rowsCount, err := usecase.userRepository.GetUserCountByEmail(ctx, request.Email)
	if err != nil {
		infrastructure.Log("got error on usecase.userRepository.GetUserCountByEmail() - RegisterUser")
		return &response.Response[interface{}]{
			Error:   true,
			Message: &response.ERROR_INTERNAL_SERVER_ERROR,
		}
	}

	if rowsCount > 0 {
		return &response.Response[interface{}]{
			Error:   true,
			Message: &response.ERROR_EMAIL_TAKEN,
		}
	}

	userUid, err := uuid.NewV6()
	if err != nil {
		infrastructure.Log("got error on uuid.NewV6() - RegisterUser")
		return &response.Response[interface{}]{
			Error:   true,
			Message: &response.ERROR_INTERNAL_SERVER_ERROR,
		}
	}

	now := time.Now()
	_user := user.User{
		UID:           userUid.String(),
		Name:          request.Name,
		Email:         request.Email,
		Gender:        request.Gender,
		ProfilePicUrl: user.PROFILE_PIC_MAP[request.Gender],
		DateOfBirth:   request.DateOfBirth,
		CreatedAt:     now.Format(time.RFC3339),
		UpdatedAt:     now.Format(time.RFC3339),
		Password:      usecase.hashUtils.HashPassword(request.Password, usecase.passwordSalt),
	}

	cacheProfile := map[string]feed.FeedProfile{}
	cacheProfile[_user.UID] = feed.FeedProfile{
		UserUid:       _user.UID,
		Name:          _user.Name,
		Gender:        _user.Gender,
		DateOfBirth:   _user.DateOfBirth,
		ProfilePicUrl: _user.ProfilePicUrl,
	}
	err = usecase.feedCacheRepository.SetProfiles(ctx, cacheProfile)
	if err != nil {
		infrastructure.Log("got error on usecase.feedCacheRepository.SetProfiles() - RegisterUser")
		return &response.Response[interface{}]{
			Error:   true,
			Message: &response.ERROR_INTERNAL_SERVER_ERROR,
		}
	}

	err = usecase.userRepository.InsertUser(ctx, _user)
	if err != nil {
		infrastructure.Log("got error on usecase.userRepository.InsertUser() - RegisterUser")
		return &response.Response[interface{}]{
			Error:   true,
			Message: &response.ERROR_INTERNAL_SERVER_ERROR,
		}
	}

	return &response.Response[interface{}]{
		Error: false,
	}
}
