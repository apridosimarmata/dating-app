package utils

import (
	"crypto/sha1"
	"fmt"
)

type HashUtils interface {
	HashPassword(plainPassword string, salt string) string
}

type hashUtils struct {
}

func NewHashUtils() HashUtils {
	return &hashUtils{}
}

func (hashUtils *hashUtils) HashPassword(plainPassword string, salt string) string {
	var sha = sha1.New()
	sha.Write([]byte(plainPassword))

	var hash = sha.Sum(nil)

	return fmt.Sprintf("%x", hash)
}
