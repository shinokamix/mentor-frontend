package model

const (
	RoleUser   = "user"
	RoleMentor = "mentor"
	RoleAdmin  = "admin"
)

type User struct {
	ID       int64  `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
	Role     string `db:"role"` // admin, mentor, user
}

type Mentor struct {
	MentorEmail string
	Contact     string
}
