package user

import (
    "time"
)

type User struct {
    ID        uint      `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    Password  string    `json:"-"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserInput struct {
    Username string `json:"username" validate:"required,min=3,max=32"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

type UpdateUserInput struct {
    Username *string `json:"username" validate:"omitempty,min=3,max=32"`
    Email    *string `json:"email" validate:"omitempty,email"`
}

type LoginUserInput struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
    User  User   `json:"user"`
    Token string `json:"token"`
}