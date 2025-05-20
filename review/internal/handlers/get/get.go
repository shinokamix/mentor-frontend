package get

import (
	"log/slog"
	"net/http"
	"review/internal/domain/model"
	"review/internal/domain/response"
	requests "review/internal/domain/resuests"
	"review/internal/lib/logger/sl"
	"review/internal/lib/validate"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type GetReviews interface {
	GetReviewsByMentorEmail(mentorEmail string) ([]model.Review, error)
}

type RedisRepo interface {
	GetReviews(email string) ([]model.Review, error, bool)
	SaveReviews(email string, reviews []model.Review) error
}

func Get(log *slog.Logger, getReview GetReviews, redisRepo RedisRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.get.Get"
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req requests.EmailMentor

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request body"))
			return
		}
		if err := validate.IsValid(&req); err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request body"))
			return
		}

		reviews, err, reviewsExists := redisRepo.GetReviews(req.Email)
		if err != nil {
			log.Error("failed to get reviews from redis", sl.Err(err))
		}

		if !reviewsExists {
			reviews, err = getReview.GetReviewsByMentorEmail(req.Email)
			if err != nil {
				log.Error("falied to get reviews", sl.Err(err))
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, response.Error("server error"))
				return
			}
			err := redisRepo.SaveReviews(req.Email, reviews)
			if err != nil {
				log.Error("falied to save reviews from redis", sl.Err(err))
			}
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, map[string]any{
			"reviews": reviews,
		})

	}
}
