package auth

import (
	"dating-app/domain"
	"dating-app/domain/auth"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type authHandler struct {
	authUsecase auth.AuthUsecase
}

func SetAuthHandler(router *chi.Mux, usecases domain.Usecases) {
	authHandler := authHandler{
		authUsecase: usecases.AuthUsecase,
	}

	router.Route("/api/v1/auth", func(r chi.Router) {
		r.Get("/login", authHandler.AuthenticateUser)
		r.Get("/refresh", authHandler.AuthenticateUser)
	})
}

func (handler *authHandler) AuthenticateUser(w http.ResponseWriter, r *http.Request) {

	response := handler.authUsecase.AuthenticateUser(r.Context(), auth.LoginRequest{})

	response.WriteResponse(w)
}

func (handler *authHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {

	response := handler.authUsecase.RefreshAccessToken(r.Context(), "todo")

	response.WriteResponse(w)
}
