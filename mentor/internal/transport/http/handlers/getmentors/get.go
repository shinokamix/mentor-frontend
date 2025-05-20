package get

import (
	"context"
	"log/slog"
	"mentor/internal/domain/models"
	"mentor/internal/domain/response"
	"mentor/internal/lib/logger/sl"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type GetMentors interface {
	Get(ctx context.Context) ([]models.MentorTable, error)
}

type RedisRepository interface {
	GetMentors(ctx context.Context) ([]models.MentorTable, error, bool)
	SaveMentors(ctx context.Context, mentor []models.MentorTable) error
}

func Get(ctx context.Context, log *slog.Logger, getMentors GetMentors, redisRepo RedisRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.getmentors.get.Get"
		log := slog.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		mentors, err, ifExists := redisRepo.GetMentors(ctx)
		if err != nil {
			log.Error("failed to get mentors mrom redis", sl.Err(err))
		}

		if !ifExists {
			mentors, err = getMentors.Get(ctx)
			if err != nil {
				log.Error("failed to get mentors", sl.Err(err))
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, response.Error("server error"))
				return
			}
			err := redisRepo.SaveMentors(ctx, mentors)
			if err != nil {
				log.Error("failed to save mentors from redis", sl.Err(err))
			}
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, map[string]any{
			"mentors": mentors,
		})
	}
}
