package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
    ID               uuid.UUID  `db:"id" json:"id"`
    UserID           uuid.UUID  `db:"user_id" json:"user_id"`
    RefreshTokenHash string     `db:"refresh_token_hash" json:"-"`
    UserAgent        *string    `db:"user_agent" json:"user_agent,omitempty"`
    IPAddress        *string    `db:"ip_address" json:"ip_address,omitempty"`
    Location         *string    `db:"location" json:"location,omitempty"`
    ExpiresAt        time.Time  `db:"expires_at" json:"expires_at"`
    CreatedAt        time.Time  `db:"created_at" json:"created_at"`
    UpdatedAt        time.Time  `db:"updated_at" json:"updated_at"`
    RevokedAt        *time.Time `db:"revoked_at" json:"revoked_at,omitempty"`
}