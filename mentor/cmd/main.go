package main

import (
	"context"
	"log/slog"
	"mentor/internal/config"
	"mentor/internal/lib/logger/sl"
	"mentor/internal/server"
	"mentor/internal/storage/cache"
	"mentor/internal/storage/db"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.LoadConfig()

	log := setupLogger(cfg.Env)

	log.Info(
		"starting mentors-service",
		slog.String("env", cfg.Env),
	)

	log.Debug("debug messages are enabled")

	storage, err := db.NewStorage(cfg.Config)
	if err != nil {
		log.Error("error created storage", sl.Err(err))
		os.Exit(1)
	}

	redisClient := cache.New(cfg.RedisConfig)

	redisRepository := cache.NewRedisRepository(redisClient)

	server, err := server.New(ctx, log, cfg, storage, redisRepository)
	if err != nil {
		log.Error("failed to create server", sl.Err(err))
		cancel()
	}

	doneChan := make(chan os.Signal, 1)
	signal.Notify(doneChan, syscall.SIGINT, syscall.SIGTERM)

	log.Info("starting server")
	go func() {
		if err := server.Start(ctx, log); err != nil {
			log.Error("failed to start server", sl.Err(err))
			os.Exit(1)
		}
	}()

	<-doneChan

	err = server.Stop(ctx)
	if err != nil {
		log.Error("failed to stop server", sl.Err(err))
		os.Exit(1)
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
