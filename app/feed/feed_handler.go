package feed

import (
	"dating-app/domain"
	"dating-app/domain/feed"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type feedHandler struct {
	feedUsecase feed.FeedUsecase
}

func SetFeedHandler(router *chi.Mux, usecases domain.Usecases) {
	feedHandler := feedHandler{
		feedUsecase: usecases.FeedUsecase,
	}

	router.Route("/api/v1/feeds", func(r chi.Router) {
		r.Use(usecases.AuthUsecase.AuthorizeRequestMiddleware)

		r.Get("/", feedHandler.GetUserProfileFeeds)
		r.Post("/like", feedHandler.LikeProfile)
		r.Post("/skip", feedHandler.SkipProfile)

	})
}

func (handler *feedHandler) GetUserProfileFeeds(w http.ResponseWriter, r *http.Request) {
	userUid := r.Context().Value("userUid")

	response := handler.feedUsecase.GetProfileFeeds(r.Context(), fmt.Sprintf("%v", userUid))

	response.WriteResponse(w)
}

func (handler *feedHandler) LikeProfile(w http.ResponseWriter, r *http.Request) {
	userUid := r.Context().Value("userUid")

	response := handler.feedUsecase.LikeProfile(r.Context(), fmt.Sprintf("%v", userUid), "todo")

	response.WriteResponse(w)
}

func (handler *feedHandler) SkipProfile(w http.ResponseWriter, r *http.Request) {
	userUid := r.Context().Value("userUid")

	response := handler.feedUsecase.SkipProfile(r.Context(), fmt.Sprintf("%v", userUid), "todo")

	response.WriteResponse(w)
}
