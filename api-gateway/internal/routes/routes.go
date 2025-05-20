package routes

import (
	"api-gateway/internal/config"
	mwLogger "api-gateway/internal/middleware/logger"
	"log"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func newProxy(target string) http.HandlerFunc {
	parsedURL, err := url.Parse(target)
	if err != nil {
		log.Fatalf("invalid proxy target %q: %v", target, err)
	}

	proxy := httputil.NewSingleHostReverseProxy(parsedURL)

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("proxy error: %v", err)
		http.Error(w, "bad gateway", http.StatusBadGateway)
	}
	return proxy.ServeHTTP
}

func NewRouter(log *slog.Logger, cfg *config.Config) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)

	auth := cfg.Auth
	router.Route("/auth", func(r chi.Router) {
		r.Post("/register", newProxy(auth))
		r.Post("/login", newProxy(auth))
		r.Post("/refresh", newProxy(auth))
		r.Post("/logout", newProxy(auth))
	})

	reviewService := cfg.Review
	router.Route("/review", func(r chi.Router) {
		r.Post("/create", newProxy(reviewService))
		r.Put("/update", newProxy(reviewService))
		r.Delete("/delete/{id}", newProxy(reviewService))
		r.Get("/get", newProxy(reviewService))
	})

	mentorService := cfg.Mentor
	router.Route("/mentors", func(r chi.Router) {
		r.Get("/get", newProxy(mentorService))
	})

	return router
}
