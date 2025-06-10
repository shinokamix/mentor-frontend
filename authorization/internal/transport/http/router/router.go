package router

import (
	"context"
	"log/slog"
	"mentorlink/internal/handlers/login"
	"mentorlink/internal/handlers/logout"
	"mentorlink/internal/handlers/profile"
	"mentorlink/internal/handlers/register"
	"mentorlink/internal/handlers/refresh"
	"mentorlink/internal/storage/db"
	"mentorlink/pkg/middleware/auth"
	"mentorlink/pkg/token"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRouter(
	ctx context.Context,
	log *slog.Logger,
	storage *db.Storage,
	tokenManager *token.TokenManager,
) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// Public routes
	r.Group(func(r chi.Router) {
		r.Post("/auth/register", register.Register(ctx, log, storage, storage))
		r.Post("/auth/login", login.Login(log, storage, tokenManager))
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware(tokenManager, log))

		r.Post("/auth/logout", logout.Logout(log, tokenManager))
		r.Post("/auth/refresh", refresh.Refresh(log, tokenManager, storage))
		r.Get("/auth/profile", profile.Get(log, storage))
	})

	return r
} 