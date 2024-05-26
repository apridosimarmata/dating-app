package user

import (
	"dating-app/domain"
	"dating-app/domain/user"
	"encoding/json"
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
		r.Post("/register", userHandler.Register)
	})
}

func (handler *userHandler) Register(w http.ResponseWriter, r *http.Request) {
	request := user.UserRegistrationRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := handler.userUsecase.RegisterUser(r.Context(), request)

	response.WriteResponse(w)
}
