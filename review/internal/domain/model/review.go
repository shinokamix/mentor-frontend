package model

import "time"

type Review struct {
	ID          int64     `json:"id,omitempty" db:"id"`
	UserID      int64     `db:"user_id"`
	MentorEmail string    `json:"mentor_email" db:"mentor_email" validate:"required,email"`
	Rating      float32   `json:"rating" db:"rating"`
	Comment     string    `json:"comment" db:"comment"`
	UserContact string    `json:"user_contact" db:"user_contact"`
	CreatedAt   time.Time `db:"created_at"`
}

type ReviewEvent struct {
	Action string  `json:"action"` // created/updated/deleted
	ID     int64   `json:"id"`
	Email  string  `json:"email"`
	Score  float32 `json:"score"`
}
