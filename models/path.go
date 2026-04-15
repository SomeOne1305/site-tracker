package models

import (
	"time"

	"github.com/google/uuid"
)
type Path struct {
    ID         uuid.UUID `db:"id" json:"id"`
    Path       string    `db:"path" json:"path"`
    VisitCount int64     `db:"visit_count" json:"visit_count"`
    ProjectID  uuid.UUID `db:"project_id" json:"project_id"`
    CreatedAt  time.Time `db:"created_at" json:"created_at"`
    UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}
