package requests

type RatingRequest struct {
	MentorEmail string  `json:"mentor_email" db:"mentor_email"`
	Rating      float32 `json:"rating" db:"rating"`
}

type MentorRequest struct {
	MentorEmail string `json:"mentor_email" db:"mentor_email"`
	Contact     string `json:"contact" db:"contact"`
}
