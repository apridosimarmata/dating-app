package presentation

import (
	"context"
	"dating-app/domain/common"
	"dating-app/infrastructure/utils"
	"fmt"

	authApp "dating-app/app/auth"

	activityApp "dating-app/app/activity"
	feedApp "dating-app/app/feed"
	subscriptionApp "dating-app/app/subscription"
	userApp "dating-app/app/user"

	"dating-app/domain"
	"dating-app/infrastructure"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis"
)

func InitServer() chi.Router {
	ctx := context.Background()
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	config := infrastructure.GetConfig()

	fmt.Println(fmt.Sprintf("config got: %v", config))

	// databases
	postgresDb, err := infrastructure.NewPostgresConn(config)
	if err != nil {
		panic("got error on infrastructure.NewPostgresConn()")

	}

	mongoDb, err := infrastructure.NewMongoDBClient(config)
	if err != nil {
		panic("got error on infrastructure.NewMongoDBClient()")
	}

	redisClient := infrastructure.NewRedisClient(ctx, config)
	cache := infrastructure.NewCache(redisClient)

	// redsync for distributed mutual exclusion
	pool := goredis.NewPool(&redisClient)
	mutexProvider := redsync.New(pool)

	repositories := domain.Repositories{
		FeedRepository:         feedApp.NewFeedRepository(postgresDb),
		FeedCacheRepository:    feedApp.NewFeedCacheRepository(cache),
		ActivityRepository:     activityApp.NewActivityRepository(mongoDb),
		UserRepository:         userApp.NewUserRepository(postgresDb),
		UserCacheRepository:    userApp.NewUserCacheRepository(cache),
		SubscriptionRepository: subscriptionApp.NewSubscriptionRepository(postgresDb),
	}

	utils := domain.Utils{
		HashUtils: utils.NewHashUtils(),
	}

	usecases := domain.Usecases{
		AuthUsecase: authApp.NewAuthUsecase(repositories, common.Secret{
			JwtSecret: config.JWT_SECRET,
		}, infrastructure.NewJwt(), utils),
		FeedUsecase:         feedApp.NewFeedUsecase(repositories, infrastructure.NewMutexProvider(*mutexProvider)),
		UserUsecase:         userApp.NewUserUsecase(repositories, "todo", utils),
		SubscriptionUsecase: subscriptionApp.NewSubscriptionUsecase(repositories),
	}

	authApp.SetAuthHandler(router, usecases)
	feedApp.SetFeedHandler(router, usecases)
	userApp.SetUserHandler(router, usecases)
	subscriptionApp.SetsubscriptionHandler(router, usecases)

	return router
}

func StopServer() {

}
