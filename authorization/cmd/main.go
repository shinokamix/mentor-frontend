package main

import (
	"context"
	"fmt"
	"log/slog"
	"mentorlink/internal/config"
	grpcclient "mentorlink/internal/grpc/client"
	"mentorlink/internal/handlers/login"
	"mentorlink/internal/handlers/logout"
	"mentorlink/internal/handlers/refresh"
	"mentorlink/internal/handlers/register"
	mwLogger "mentorlink/internal/middleware/logger"
	"os/signal"
	"syscall"
	"time"

	"mentorlink/internal/lib/logger/sl"
	"mentorlink/internal/storage/cache"
	"mentorlink/internal/storage/db"
	"mentorlink/pkg/token"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	cfg := config.LoadConfig()

	log := setupLogger(cfg.Env)

	log.Info(
		"starting auth-service",
		slog.String("env", cfg.Env),
	)

	log.Debug("debug messages are enabled")

	redisClient := cache.New(cfg.RedisConfig)

	redisRepository := cache.NewRedisRepository(redisClient)

	storage, err := db.NewStorage(cfg.Config)
	if err != nil {
		log.Error("error creation storage", sl.Err(err))
		os.Exit(1)
	}

	tokemMn, err := token.NewTokenmanagerRSA(cfg.PrivateKeyPath, cfg.PublicKeyPath)
	if err != nil {
		log.Error("error with token manager", sl.Err(err))
		os.Exit(1)
	}

	client, err := grpcclient.NewMentorClient(fmt.Sprintf("mentor-server:%s", cfg.MentorServiceAddress))
	if err != nil {
		log.Error("error with new grpc client", sl.Err(err))
		os.Exit(1)
	}

	defer func() {
		if err := client.Close(); err != nil {
			log.Error("failed to close gRPC client", sl.Err(err))
		} else {
			log.Debug("gRPC client closed successfully")
		}
	}()

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.URLFormat)

	router.Post("/auth/register", register.Register(context.Background(), log, storage, client))
	router.Post("/auth/login", login.Login(log, storage, tokemMn))
	router.Post("/auth/logout", logout.Logout(log, redisRepository, tokemMn))
	router.Post("/auth/refresh", refresh.RefreshTokens(log, redisRepository, tokemMn))

	log.Info("starting server", slog.String("adsress", cfg.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("failed to start server", sl.Err(err))
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", sl.Err(err))

		return
	}

	log.Info("server stopped")

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
