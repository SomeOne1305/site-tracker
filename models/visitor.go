package models

import (
	"time"

	"github.com/google/uuid"
)
type Visitor struct {
    ID        uuid.UUID  `db:"id" json:"id"`
    IPAddress string     `db:"ip_address" json:"ip_address"`
    UserAgent *string    `db:"user_agent" json:"user_agent,omitempty"`
    VisitTime time.Time  `db:"visit_time" json:"visit_time"`
    Country   *string    `db:"country" json:"country,omitempty"`
    PathID    uuid.UUID  `db:"path_id" json:"path_id"`
}