package profile

import (
	"log/slog"
	"mentorlink/internal/domain/model"
	"mentorlink/internal/domain/response"
	"mentorlink/internal/lib/logger/sl"
	"mentorlink/pkg/token"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type UserGetter interface {
	GetByEmail(email string) (*model.User, error)
}

func Get(log *slog.Logger, userGetter UserGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.profile.Get"
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		claims, ok := r.Context().Value("user").(*token.Claims)
		if !ok || claims == nil {
			log.Error("failed to get user claims")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, response.Error("unauthorized"))
			return
		}

		user, err := userGetter.GetByEmail(claims.Email)
		if err != nil {
			log.Error("failed to get user", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to get user profile"))
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, map[string]any{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		})
	}
} 