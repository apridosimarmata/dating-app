package auth

import (
	"dating-app/domain"
	"dating-app/domain/auth"
	"encoding/json"
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
		r.Post("/login", authHandler.AuthenticateUser)
		r.Post("/refresh", authHandler.RefreshToken)
	})
}

func (handler *authHandler) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	request := auth.LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := handler.authUsecase.AuthenticateUser(r.Context(), request)

	response.WriteResponse(w)
}

func (handler *authHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	request := auth.RefreshTokenRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := handler.authUsecase.RefreshAccessToken(r.Context(), *request.RefreshToken)

	response.WriteResponse(w)
}
