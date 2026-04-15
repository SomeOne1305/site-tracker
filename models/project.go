package models

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
    ID           uuid.UUID `db:"id" json:"id"`
    ProjectToken string    `db:"project_token" json:"project_token"`
    Name         string    `db:"name" json:"name"`
    Description  *string   `db:"description" json:"description,omitempty"`
    OwnerID      uuid.UUID `db:"owner_id" json:"owner_id"`
    CreatedAt    time.Time `db:"created_at" json:"created_at"`
    UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type CreateProjectRequest struct {
    Name        string  `json:"name" validate:"required,min=2,max=100"`
    Description *string `json:"description,omitempty" validate:"omitempty,max=500"`
}