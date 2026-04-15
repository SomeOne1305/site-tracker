package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
    ID               uuid.UUID `db:"id" json:"id"`
    FirstName        string    `db:"first_name" json:"first_name"`
    LastName         string    `db:"last_name" json:"last_name"`
    Email            string    `db:"email" json:"email"`
    IsVerified       bool      `db:"is_verified" json:"is_verified"`
    Password         string    `db:"password" json:"-"`
    CreatedAt        time.Time `db:"created_at" json:"created_at"`
    UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

type CreateUserRequest struct {
    FirstName string `json:"first_name" validate:"required,min=2,max=100"`
    LastName  string `json:"last_name" validate:"required,min=2,max=100"`
    Email     string `json:"email" validate:"required,email"`
    Password  string `json:"password" validate:"required,min=8"`
}

type UpdateUserRequest struct {
    FirstName *string `json:"first_name,omitempty" validate:"omitempty,min=2,max=100"`
    LastName  *string `json:"last_name,omitempty" validate:"omitempty,min=2,max=100"`
}

type VerifyEmailRequest struct {
    Email            string `json:"email" validate:"required,email"`
    VerificationCode int    `json:"verification_code" validate:"required"`
}

type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}