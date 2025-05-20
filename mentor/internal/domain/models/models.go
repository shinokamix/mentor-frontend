package models

type MentorTable struct {
	MentorEmail   string  `json:"mentor_email" db:"mentor_email"`
	Contact       string  `json:"contact" db:"contact"`
	AverageRating float32 `json:"average_rating" db:"average_rating"`
}
