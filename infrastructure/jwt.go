package infrastructure

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type JwtInterface interface {
	GenerateAccessToken(userUid string, jwtSecret string, isRefreshToken bool) (*string, error)
	ValidateToken(tokenString string, jwtSecret string) (userUid *string, ttl int64, err error)
}

type JWT struct {
}

func NewJwt() JwtInterface {
	return &JWT{}
}

func (_jwt *JWT) GenerateAccessToken(userUid string, jwtSecret string, isRefreshToken bool) (*string, error) {
	exp := time.Now().Add(time.Minute * 5).Unix()
	if isRefreshToken {
		exp = time.Now().Add(time.Hour * 24 * 30).Unix()
	}
	t := jwt.NewWithClaims(jwt.SigningMethodES256,
		jwt.MapClaims{
			"sub": userUid,
			"exp": exp,
		})

	signedToken, err := t.SignedString(jwtSecret)
	if err != nil {
		Log("got error on usecase.getSubscriberFeedProfiles() - GetProfileFeeds")
		return nil, err
	}

	return &signedToken, nil
}

func (_jwt *JWT) ValidateToken(tokenString string, jwtSecret string) (userUid *string, ttl int64, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, 0, err
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userUid := payload["sub"].(string)
		expStringEpoch := payload["exp"].(int64)
		ttl := expStringEpoch - time.Now().Unix()
		return &userUid, ttl, nil
	}

	return nil, 0, errors.New("invalid token")
}
