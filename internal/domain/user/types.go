package user

import (
	"ecommerce-app/internal/pkg/response"
	"time"
)

type User struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role 	string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserWithToken struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

type UsersWithMeta struct {
	Users []User     `json:"users"`
	Meta  response.Meta `json:"meta"`
}