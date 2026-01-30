package model

type User struct {
	ID       int64  `db:"id"`
	Username string `db:"username" validate:"required,min=3,max=32,alphanum"`
	Name     string `db:"name" validate:"required,min=1,max=64"`
	Password string `db:"password" validate:"required,min=6"`
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=32,alphanum"`
	Name     string `json:"name" validate:"required,min=1,max=64"`
	Password string `json:"password" validate:"required,min=6"`
}
