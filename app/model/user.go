package model

type User struct {
	ID       int     `db:"id"`
	Username string  `db:"username"`
	Email    string  `db:"email"`
	Password string  `db:"password"`
	Bio      *string `db:"bio"`
	Image    *string `db:"image"`
}

func NewUser() *User { return &User{} }

type RegisterUser struct {
	Username string `json:"username" validate:"lte=255,gte=5"`
	Email    string `json:"email" validate:"lte=255,email"`
	Password string `json:"password" validate:"lte=100"`
}

type RegisterDTO struct {
	User RegisterUser `json:"user"`
}

type LoginUser struct {
	Email    string `json:"email" validate:"lte=255,email"`
	Password string `json:"password" validate:"lte=100"`
}

type LoginDTO struct {
	User LoginUser `json:"user"`
}

type UpdateUser struct {
	Username string `json:"username" validate:"lte=255,gte=5"`
	Password string `json:"password" validate:"lte=100"`
	Bio      string `json:"bio" validate:"lte=255"`
	Image    string `json:"image"`
}
