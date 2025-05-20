package requests

type EmailMentor struct {
	Email string `json:"mentor_email" db:"email" validate:"required,email"`
}
