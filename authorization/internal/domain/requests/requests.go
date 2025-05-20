package requests

type Register struct {
	Email          string `json:"email" validate:"required,email"`
	Password       string `json:"password" validate:"required,min=6"`
	RepeatPassword string `json:"repeat_password" validate:"required,eqfield=Password"`
	Role           string `json:"role" validate:"required,oneof=user mentor admin"`
	Contact        string `json:"contact"`
}

type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RFToken struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
