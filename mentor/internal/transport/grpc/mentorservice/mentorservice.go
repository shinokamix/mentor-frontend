package mentorservice

import (
	"context"
	"fmt"
	"log/slog"
	"mentor/internal/domain/requests"
	client "mentor/pkg/api/proto"
)

const (
	ActionUpdate = "updated"
	ActionDelete = "deleted"
)

type PostgresRepository interface {
	UpdateMentor(ctx context.Context, mentor *requests.RatingRequest) error
	DeleteReviewByMentor(ctx context.Context, mentor *requests.RatingRequest) error
	CreateMentor(ctx context.Context, mentor *requests.MentorRequest) error
	MentorExists(ctx context.Context, mentorEmail string) (bool, error)
}

type MentorService struct {
	client.UnimplementedMentorServiceServer
	repo PostgresRepository
	log  *slog.Logger
}

func NewMentorService(log *slog.Logger, repo PostgresRepository) *MentorService {
	return &MentorService{log: log, repo: repo}
}

func (s *MentorService) MethodMentorRating(ctx context.Context, req *client.RatingRequest) (*client.Response, error) {
	s.log.Debug("processing rating reuqest",
		"mentor_email", req.MentorEmail,
		"rating", req.Rating,
		"action", req.Action,
	)
	request := &requests.RatingRequest{
		MentorEmail: req.MentorEmail,
		Rating:      req.Rating,
	}

	switch req.Action {
	case ActionDelete:
		s.log.Info("starting review deletion", "mentor_email", req.MentorEmail)
		err := s.repo.DeleteReviewByMentor(ctx, request)
		if err != nil {
			s.log.Error("review deletion failed",
				"error", err,
				"mentor_email", req.MentorEmail)
			return &client.Response{
					Success: false,
					Message: "error",
				},
				fmt.Errorf("failed to delete review: %w", err)
		}
		s.log.Info("mentor successfully deleted", "mentor_email", req.MentorEmail)
		return &client.Response{
			Success: true,
			Message: "ok",
		}, nil

	case ActionUpdate:
		s.log.Info("starting mentor update", "mentor_email", req.MentorEmail)
		err := s.repo.UpdateMentor(ctx, request)
		if err != nil {
			s.log.Error("mentor update failed",
				"error", err,
				"mentor_email", req.MentorEmail,
				"rating", req.Rating)
			return &client.Response{
					Success: false,
					Message: "error",
				},
				fmt.Errorf("failed to update review: %w", err)
		}
		s.log.Info("mentor successfully updated", "mentor_email", req.MentorEmail)
		return &client.Response{
			Success: true,
			Message: "ok",
		}, nil

	default:
		s.log.Warn("unknown action requested",
			"action", req.Action,
			"mentor_email", req.MentorEmail)
		return &client.Response{
				Success: false,
				Message: "error: action don't matched",
			},
			nil

	}
}

func (s *MentorService) NewMentor(ctx context.Context, req *client.MentorRequest) (*client.Response, error) {
	s.log.Debug("creating new mentor",
		"mentor_email", req.MentorEmail,
		"contact", req.Contact)

	request := &requests.MentorRequest{
		MentorEmail: req.MentorEmail,
		Contact:     req.Contact,
	}

	err := s.repo.CreateMentor(ctx, request)
	if err != nil {
		s.log.Error("mentor creation failed",
			"error", err,
			"mentor_email", req.MentorEmail,
			"contact", req.Contact)
		return &client.Response{
				Success: false,
				Message: "error",
			},
			fmt.Errorf("failed to create mentor: %w", err)
	}

	s.log.Info("mentor successfully created", "mentor_email", req.MentorEmail)
	return &client.Response{
		Success: true,
		Message: "ok",
	}, nil
}

func (s *MentorService) CheckMentor(ctx context.Context, req *client.CheckRequest) (*client.CheckResponse, error) {
	mentorEmail := req.MentorEmail

	s.log.Debug("checking mentor existence", "mentor_email", mentorEmail)

	exists, err := s.repo.MentorExists(ctx, mentorEmail)
	if err != nil {
		s.log.Error("failed to check mentor existence", "error", err, "mentor_email", mentorEmail)

		return &client.CheckResponse{
			Success: false,
			Exists:  false,
			Message: "error",
		}, fmt.Errorf("failed to check mentor: %w", err)
	}

	if exists {
		s.log.Info("mentor exists", "mentor_email", mentorEmail)
		return &client.CheckResponse{
			Success: true,
			Exists:  true,
			Message: "exists",
		}, nil
	}

	s.log.Info("mentor does not exist", "mentor_email", mentorEmail)
	return &client.CheckResponse{
		Success: true,
		Exists:  false,
		Message: "not exists",
	}, nil
}
