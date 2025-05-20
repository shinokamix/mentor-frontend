package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"rating/internal/config"
	"rating/internal/kafka"
	grpccleint "rating/internal/transport/grpc/client"
	"syscall"
)

func main() {
	cfg := config.LoadConfig()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("starting rating server")

	mentorClient, err := grpccleint.NewMentorClient(cfg.MentorServiceAddress)
	if err != nil {
		logger.Error("failed to connect to mentor service", "error", err)
		os.Exit(1)
	}
	defer mentorClient.Close()

	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	consumer, err := kafka.NewConsumer([]string{cfg.KafkaBroker}, cfg.KafkaGroupID, mentorClient, logger)
	if err != nil {
		logger.Error("failed to create consumer", "error", err)
		os.Exit(1)
	}

	go consumer.Run(ctx, cfg.KafkaTopic)

	<-sigChan
	logger.Info("signal caught, shutting down gracefulle")

	cancel()

	closeErr := consumer.Close()
	if closeErr != nil {
		logger.Error("error while closing consumer", "error", closeErr)
	}

	logger.Info("consumer service stopped")
}
