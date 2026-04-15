package repository

import (
	"context"
	"fmt"
	"time"
	"visit-tracker/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepository struct {
	pool *pgxpool.Pool
}

func NewSessionRepository(pool *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{pool: pool}
}

func (r *SessionRepository) CreateSession(ctx context.Context, userID uuid.UUID, refreshToken string, userAgent string, ipAddress string, location string, expiresAt time.Time) (*models.Session, error) {
	query := `INSERT INTO sessions (user_id, refresh_token_hash, user_agent, ip_address, location, expires_at,  created_at) VALUES ($1, $2, $3, $4, $5, $6, NOW()) RETURNING id, user_id, refresh_token_hash, created_at`

	var session models.Session
	err := r.pool.QueryRow(ctx, query, userID, refreshToken, userAgent, ipAddress, location, expiresAt).Scan(
		&session.ID,
		&session.UserID,
		&session.RefreshTokenHash,
		&session.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}
	return &session, nil
}

func (r *SessionRepository) GetSessionById(ctx context.Context, id uuid.UUID) (*models.Session, error) {
	query := `SELECT id, user_id, refresh_token_hash, created_at FROM sessions WHERE id = $1`

	var session models.Session
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&session.ID,
		&session.UserID,
		&session.RefreshTokenHash,
		&session.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to load session: %w", err)
	}
	return &session, nil
}

func (r *SessionRepository) DeleteSession(ctx context.Context, refreshToken string) error {
	query := `DELETE FROM sessions WHERE refresh_token_hash = $1`
	_, err := r.pool.Exec(ctx, query, refreshToken)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}
