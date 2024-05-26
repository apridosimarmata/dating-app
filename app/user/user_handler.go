package user

import (
	"dating-app/domain"
	"dating-app/domain/user"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type userHandler struct {
	userUsecase user.UserUsecase
}

func SetUserHandler(router *chi.Mux, usecases domain.Usecases) {
	userHandler := userHandler{
		userUsecase: usecases.UserUsecase,
	}

	router.Route("/api/v1/users", func(r chi.Router) {
		r.Use(usecases.AuthUsecase.AuthorizeRequestMiddleware)

		r.Post("/login", userHandler.Login)
		r.Post("/register", userHandler.Register)

	})
}

func (handler *userHandler) Login(w http.ResponseWriter, r *http.Request) {

}

func (handler *userHandler) Register(w http.ResponseWriter, r *http.Request) {

}
