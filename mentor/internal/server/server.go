package server

import (
	"context"
	"fmt"
	"log/slog"
	"mentor/internal/config"
	"mentor/internal/domain/models"
	"mentor/internal/domain/requests"
	"mentor/internal/transport/grpc/mentorservice"
	get "mentor/internal/transport/http/handlers/getmentors"
	client "mentor/pkg/api/proto"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type PostgresRepository interface {
	UpdateMentor(ctx context.Context, mentor *requests.RatingRequest) error
	DeleteReviewByMentor(ctx context.Context, mentor *requests.RatingRequest) error
	CreateMentor(ctx context.Context, mentor *requests.MentorRequest) error
	Get(ctx context.Context) ([]models.MentorTable, error)
	MentorExists(ctx context.Context, mentorEmail string) (bool, error)
}

type RedisRepository interface {
	GetMentors(ctx context.Context) ([]models.MentorTable, error, bool)
	SaveMentors(ctx context.Context, mentor []models.MentorTable) error
}

type Server struct {
	grpcServer   *grpc.Server
	httpServer   *http.Server
	grpcListener net.Listener
}

func New(ctx context.Context, log *slog.Logger, cfg *config.Config, postgresRepository PostgresRepository, redisRepository RedisRepository) (*Server, error) {
	gRPCaddr := fmt.Sprintf(":%d", cfg.GRPCPort)
	grpcListener, err := net.Listen("tcp", gRPCaddr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on %s: %w", gRPCaddr, err)
	}

	opts := []grpc.ServerOption{}
	grpcSrv := grpc.NewServer(opts...)
	mentorSrv := mentorservice.NewMentorService(log, postgresRepository)
	client.RegisterMentorServiceServer(grpcSrv, mentorSrv)

	router := chi.NewRouter()
	router.Get("/mentors/get", get.Get(ctx, log, postgresRepository, redisRepository))

	httpSrv := &http.Server{
		Addr:         cfg.AddressServerHTTP,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &Server{
		grpcServer:   grpcSrv,
		httpServer:   httpSrv,
		grpcListener: grpcListener,
	}, nil

}

func (s *Server) Start(ctx context.Context, log *slog.Logger) error {
	eg := errgroup.Group{}

	log.Info("starting server")

	eg.Go(func() error {
		log.Debug("gRPPC server starting", "addr", s.grpcListener.Addr())
		if err := s.grpcServer.Serve(s.grpcListener); err != nil {
			log.Error("gRPC server failed", "error", err)
			return fmt.Errorf("gRPC server error: %w", err)
		}
		return nil
	})

	eg.Go(func() error {
		log.Debug("http server starting", "addr", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("http server failed", "error", err)
			return fmt.Errorf("HTTP server error: %w", err)
		}
		return nil
	})

	log.Info("server started successfully")
	return eg.Wait()
}

func (s *Server) Stop(ctx context.Context) error {
	s.grpcServer.GracefulStop()
	return s.httpServer.Shutdown(ctx)
}
