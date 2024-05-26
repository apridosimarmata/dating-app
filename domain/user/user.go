package user

import (
	"context"
	"dating-app/domain/common/response"
)

var PROFILE_PIC_MAP = map[string]string{
	"M": "https://static.vecteezy.com/system/resources/previews/004/991/321/original/picture-profile-icon-male-icon-human-or-people-sign-and-symbol-vector.jpg",
	"F": "https://www.shutterstock.com/image-vector/person-gray-photo-placeholder-woman-600nw-1241538838.jpg",
}

type User struct {
	UID           string `json:"uid" gorm:"column=uid"`
	Name          string `json:"name" gorm:"column=name"`
	Email         string `json:"email" gorm:"column=email"`
	Gender        string `json:"gender" gorm:"column=gender"`
	ProfilePicUrl string `json:"profile_pic_url" gorm:"column=profile_pic_url"`
	Password      string `json:"password" gorm:"column=password"`
	DateOfBirth   string `json:"date_of_birth" gorm:"column=date_of_birth"`
	CreatedAt     string `json:"created_at" gorm:"column=created_at"`
	UpdatedAt     string `json:"updated_at" gorm:"column=updated_at"`
}

type UserRegistrationRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Gender      string `json:"gender"`
	DateOfBirth string `json:"date_of_birth"`
	Password    string `json:"password"`
}

type UserCacheDetails struct {
	Uid          string `json:"uid" gorm:"column=uid"`
	IsSubscriber bool   `json:"is_subscriber" gorm:"column=is_subscriber"`
}

type UserRepository interface {
	InsertUser(ctx context.Context, user User) (err error)
	GetUserCountByEmail(ctx context.Context, email string) (res int, err error)
	GetUserCountByUid(ctx context.Context, uid string) (res int, err error)
	GetUserMiniDetailsByUserUid(ctx context.Context, userUid string) (res *UserCacheDetails, err error)
	GetUserDetailsByEmail(ctx context.Context, email string) (res *User, err error)
}

type UserCacheRepository interface {
	SetUserCacheDetails(ctx context.Context, userUid string, user UserCacheDetails) (err error)
	GetUserCacheDetails(ctx context.Context, userUid string) (res *UserCacheDetails, err error)
}

type UserUsecase interface {
	RegisterUser(ctx context.Context, request UserRegistrationRequest) (res *response.Response[interface{}])
}
