package db

import (
	"database/sql"
	"errors"
	"fmt"
	"review/internal/domain/model"

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
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.UserName, cfg.Password, cfg.DBName)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error with connect to database: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) IfExist(userID int64, mentorEmail string) (bool, error) {
	const op = "storage.db.ifExist"
	query := `SELECT id FROM reviews WHERE user_id=$1 and mentor_email=$2`
	var id int64
	err := s.db.Get(&id, query, userID, mentorEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("%s, %w", op, err)
	}

	return true, nil
}

func (s *Storage) CreateReview(review *model.Review) (int64, error) {
	const op = "storage.db.CreateReview"
	query := `INSERT INTO reviews (user_id, mentor_email, rating, comment, user_contact, created_at)
			  VALUES ($1, $2, $3, $4, $5, $6)
			  RETURNING id;`

	var newID int64
	err := s.db.QueryRow(query, review.UserID, review.MentorEmail, review.Rating, review.Comment, review.UserContact, review.CreatedAt).Scan(&newID)
	if err != nil {
		return -1, fmt.Errorf("%s, %w", op, err)
	}
	review.ID = newID
	return newID, nil
}

func (s *Storage) UpdateReview(review *model.Review) error {
	const op = "storage.db.UpdateReview"
	query := `UPDATE reviews
			  SET mentor_email=$1, rating=$2, comment=$3, user_contact=$4
			  WHERE id=$5 and user_id=$6;`
	result, err := s.db.Exec(query, review.MentorEmail, review.Rating, review.Comment, review.UserContact, review.ID, review.UserID)
	if err != nil {
		return fmt.Errorf("%s, %w", op, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to check rows affected: %w", op, err)
	}

	if rows == 0 {
		return fmt.Errorf("%s: no review found with id=%d and user_id=%d", op, review.ID, review.UserID)
	}
	return nil
}

func (s *Storage) GetReviewsByMentorEmail(mentorEmail string) ([]model.Review, error) {
	const op = "storage.db.GetReviewsBeMentorEmail"
	query := `SELECT id, mentor_email, rating, comment, user_contact, created_at
			  FROM reviews 
			  WHERE mentor_email=$1 
			  ORDER BY created_at DESC;`
	var reviews []model.Review
	err := s.db.Select(&reviews, query, mentorEmail)
	if err != nil {
		return []model.Review{}, fmt.Errorf("%s, %w", op, err)
	}

	return reviews, nil
}

func (s *Storage) DeleteReview(userID, id int64) error {
	const op = "storage.db.DeleteReview"
	query := `DELETE FROM reviews
			  WHERE id=$1 and user_id=$2;`
	result, err := s.db.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("%s, %w", op, err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("%s: review with %d id not found", op, id)
	}
	return nil
}

func (s *Storage) GetReviewByID(id int64) (*model.Review, error) {
	const op = "storage.db.GetReviewByID"
	query := `SELECT id, mentor_email, rating, comment, user_contact, created_at
			  FROM reviews 
			  WHERE id=$1`

	var review model.Review
	err := s.db.Get(&review, query, id)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	return &review, nil
}
