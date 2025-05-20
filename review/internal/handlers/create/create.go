package create

import (
	"context"
	"log/slog"
	"net/http"
	"review/internal/domain/model"
	"review/internal/domain/response"
	"review/internal/lib/logger/sl"
	"review/internal/lib/validate"
	mwAuth "review/internal/middleware/auth"
	"review/pkg/token"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type ReviewCreater interface {
	CreateReview(review *model.Review) (int64, error)
	IfExist(userID int64, mentorEmail string) (bool, error)
}

type KafkaProducer interface {
	SendReviewEvent(review *model.ReviewEvent) error
}

type CheckMentor interface {
	CheckMentor(ctx context.Context, mentorEmail string) (bool, error)
}

func Create(ctx context.Context, log *slog.Logger, reviewCreater ReviewCreater, kafkaProducer KafkaProducer, checkMentor CheckMentor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.create.Create"
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		claims := r.Context().Value(mwAuth.UserKey).(*token.Claims)

		var req model.Review
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
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

		existsMentor, err := checkMentor.CheckMentor(ctx, req.MentorEmail)
		if err != nil {
			log.Error("failed to check mentor in mentor-service", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("server error"))
			return
		}

		if !existsMentor {
			log.Error("mentor doesn't exists", "mentor_email", req.MentorEmail)
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, response.Error("mentor doesn't exists"))
			return
		}

		existsReview, err := reviewCreater.IfExist(req.UserID, req.MentorEmail)
		if err != nil {
			log.Error("falied to find review", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("server error"))
			return
		}

		if existsReview {
			log.Warn("review already exist")
			render.Status(r, http.StatusConflict)
			render.JSON(w, r, map[string]any{
				"status": "review already exist",
			})
			return
		}

		req.CreatedAt = time.Now()
		id, err := reviewCreater.CreateReview(&req)
		if err != nil {
			log.Error("falied to create review", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("server error"))
			return
		}

		event := &model.ReviewEvent{
			Action: "updated",
			ID:     id,
			Email:  req.MentorEmail,
			Score:  req.Rating,
		}

		if err := kafkaProducer.SendReviewEvent(event); err != nil {
			log.Error("failed to send kafka even", sl.Err(err))
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, map[string]any{
			"id": id,
		})
	}
}
