package db

import (
	"context"
	"fmt"
	"mentor/internal/domain/models"
	"mentor/internal/domain/requests"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	UserName string `env:"POSTGRES_USER" env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	Host     string `env:"POSTGRES_HOST" env-required:"true"`
	Port     string `env:"POSTGRES_PORT" env-required:"true"`
	DBName   string `env:"POSTGRES_DB" env-required:"true"`
}

type Storage struct {
	db *sqlx.DB
}

func NewStorage(cfg Config) (*Storage, error) {
	dsn := fmt.Sprintf("host=%s port=%s  user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.UserName, cfg.Password, cfg.DBName)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error with connect to database: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) CreateMentor(ctx context.Context, mentor *requests.MentorRequest) error {
	const op = "storage.db.postgres.SaveMentor"
	queury := `INSERT INTO mentors (mentor_email, contact)
			   VALUES($1, $2)
			   RETURNING id`
	var newID int64
	err := s.db.QueryRow(queury, mentor.MentorEmail, mentor.Contact).Scan(&newID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) Get(ctx context.Context) ([]models.MentorTable, error) {
	const op = "storage.db.postgres.Get"
	query := `SELECT mentor_email, contact, average_rating
			  FROM mentors
			  ORDER BY average_rating DESC;`

	var mentors []models.MentorTable
	err := s.db.Select(&mentors, query)
	if err != nil {
		return []models.MentorTable{}, fmt.Errorf("%s, %w", op, err)
	}
	return mentors, nil
}

func (s *Storage) UpdateMentor(ctx context.Context, mentor *requests.RatingRequest) error {
	const op = "storage.db.postgres.UpdateMentor"
	query := `UPDATE mentors
			  SET count_reviews = count_reviews + 1, sum_rating = sum_rating + $1
			  WHERE mentor_email=$2`
	_, err := s.db.Exec(query, mentor.Rating, mentor.MentorEmail)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteReviewByMentor(ctx context.Context, mentor *requests.RatingRequest) error {
	const op = "storage.db.postgres.DeleteReviewByMentor"
	query := `UPDATE mentors
			  SET count_reviews = count_reviews - 1, sum_rating = sum_rating - $1
			  WHERE mentor_email=$2`
	_, err := s.db.Exec(query, mentor.Rating, mentor.MentorEmail)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) MentorExists(ctx context.Context, mentorEmail string) (bool, error) {
	const op = "storage.db.postgres.CheckMentorByEmail"
	query := `SELECT EXISTS(SELECT 1 FROM mentors WHERE mentor_email=$1)`
	var exists bool
	err := s.db.GetContext(ctx, &exists, query, mentorEmail)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return exists, nil
}
