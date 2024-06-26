package infrastructure

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

func NewPostgresConn(config Config) (conn *gorm.DB, err error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", config.POSTGRES_HOST, config.POSTGRES_USER, config.POSTGRES_PASSWORD, config.POSTGRES_DB, config.POSTGRES_PORT)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
