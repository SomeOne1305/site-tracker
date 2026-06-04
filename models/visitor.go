package models

import (
	"time"

	"github.com/google/uuid"
)

type Visitor struct {
	ID         uuid.UUID `db:"id" json:"id"`
	PathID     uuid.UUID `db:"path_id" json:"path_id"`
	SessionID  string    `db:"session_id" json:"session_id"`
	IPAddress  string    `db:"ip_address" json:"ip_address"`
	UserAgent  *string   `db:"user_agent" json:"user_agent,omitempty"`
	Country    *string   `db:"country" json:"country,omitempty"`
	FirstVisit time.Time `db:"first_visit" json:"first_visit"`
	LastVisit  time.Time `db:"last_visit" json:"last_visit"`
	VisitCount int32     `db:"visit_count" json:"visit_count"`
}

type DailyStats struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	ProjectID      uuid.UUID  `db:"project_id" json:"project_id"`
	PathID         *uuid.UUID `db:"path_id" json:"path_id,omitempty"`
	Date           time.Time  `db:"date" json:"date"`
	Pageviews      int64      `db:"pageviews" json:"pageviews"`
	UniqueVisitors int64      `db:"unique_visitors" json:"unique_visitors"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
}

type ProjectStats struct {
	TotalViews     int64 `json:"total_views"`
	UniqueVisitors int64 `json:"unique_visitors"`
}

type DailyData struct {
	Date           time.Time `json:"date"`
	Pageviews      int64     `json:"pageviews"`
	UniqueVisitors int64     `json:"unique_visitors"`
}

type PageStats struct {
	Path           string `json:"path"`
	Pageviews      int64  `json:"pageviews"`
	UniqueVisitors int64  `json:"unique_visitors"`
}
