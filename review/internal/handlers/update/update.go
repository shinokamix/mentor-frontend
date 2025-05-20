package update

import (
	"log/slog"
	"net/http"
	"review/internal/domain/model"
	"review/internal/domain/response"
	"review/internal/lib/logger/sl"
	"review/internal/lib/validate"
	mwAuth "review/internal/middleware/auth"
	"review/pkg/token"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type ReviewUpdate interface {
	UpdateReview(review *model.Review) error
	GetReviewByID(id int64) (*model.Review, error)
}

type KafkaProducer interface {
	SendReviewEvent(review *model.ReviewEvent) error
}

func Update(log *slog.Logger, reviewUpdate ReviewUpdate, kafkaProducer KafkaProducer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.update.Update"
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		claims, ok := r.Context().Value(mwAuth.UserKey).(*token.Claims)
		if !ok || claims == nil {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, response.Error("unauthorized"))
			return
		}

		var req model.Review
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request body"))
			return
		}

		req.UserID = claims.UserID

		if err := validate.IsValid(req); err != nil {
			log.Error("validation error", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("server error"))
			return
		}

		rev, err := reviewUpdate.GetReviewByID(req.ID)
		if err != nil {
			log.Error("failed to get review by id", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("server error"))
			return
		}

		event := &model.ReviewEvent{
			Action: "deleted",
			ID:     rev.ID,
			Email:  rev.MentorEmail,
			Score:  rev.Rating,
		}

		if err := kafkaProducer.SendReviewEvent(event); err != nil {
			log.Error("failed to send kafka event", sl.Err(err))
		}

		if err := reviewUpdate.UpdateReview(&req); err != nil {
			log.Error("failed to update review", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("server error"))
			return
		}

		event = &model.ReviewEvent{
			Action: "updated",
			ID:     req.ID,
			Email:  req.MentorEmail,
			Score:  req.Rating,
		}

		if err := kafkaProducer.SendReviewEvent(event); err != nil {
			log.Error("failed to send kafka event", sl.Err(err))
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, map[string]any{
			"status": "updated",
		})
	}
}
