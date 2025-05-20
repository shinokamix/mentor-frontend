package del

import (
	"log/slog"
	"net/http"
	"review/internal/domain/model"
	"review/internal/domain/response"
	"review/internal/lib/logger/sl"
	mwAuth "review/internal/middleware/auth"
	"review/pkg/token"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type DelReview interface {
	DeleteReview(userID, id int64) error
	GetReviewByID(id int64) (*model.Review, error)
}

type KafkaProducer interface {
	SendReviewEvent(review *model.ReviewEvent) error
}

func Delete(log *slog.Logger, delreview DelReview, kafkaProducer KafkaProducer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.delete.Delete"
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

		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Error("invalid ID format", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid ID"))
			return
		}

		rev, err := delreview.GetReviewByID(id)
		if err != nil {
			log.Error("failed to get review by id", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("server error"))
			return
		}

		if err := delreview.DeleteReview(claims.UserID, id); err != nil {
			log.Error("failed to delete review", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("server error"))
			return
		}

		event := &model.ReviewEvent{
			Action: "deleted",
			ID:     id,
			Email:  rev.MentorEmail,
			Score:  rev.Rating,
		}

		if err := kafkaProducer.SendReviewEvent(event); err != nil {
			log.Error("failed to send kafka event", sl.Err(err))
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, map[string]any{
			"status": "review deleted",
		})
	}
}
