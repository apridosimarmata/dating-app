package subscription

import (
	"dating-app/domain"
	"dating-app/domain/subscription"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type subscriptionHandler struct {
	subscriptionUsecase subscription.SubscriptionUsecase
}

func SetsubscriptionHandler(router *chi.Mux, usecases domain.Usecases) {
	subscriptionHandler := subscriptionHandler{
		subscriptionUsecase: usecases.SubscriptionUsecase,
	}

	router.Route("/api/v1/external/subscriptions", func(r chi.Router) {
		r.Post("/payment-callback", subscriptionHandler.HandlePaymentCallback)
	})

	router.Route("/api/v1/subscriptions", func(r chi.Router) {
		r.Use(usecases.AuthUsecase.AuthorizeRequestMiddleware)

		r.Post("/", subscriptionHandler.Subscribe)
	})
}

func (handler *subscriptionHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	userUid := r.Context().Value("userUid")

	request := subscription.CreateSubscriptionRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	request.UserUid = userUid.(string)
	response := handler.subscriptionUsecase.Subscribe(r.Context(), request)

	response.WriteResponse(w)
}

func (handler *subscriptionHandler) HandlePaymentCallback(w http.ResponseWriter, r *http.Request) {
	response := handler.subscriptionUsecase.HandleSubscriptionPaymentCallback(r.Context(), subscription.SubscriptionPaymentCallback{})

	response.WriteResponse(w)
}
