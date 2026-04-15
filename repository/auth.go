package repository

import (
	"context"
	"fmt"
	"visit-tracker/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	pool *pgxpool.Pool
}

func NewAuthRepository(pool *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{pool: pool}
}

func (r *AuthRepository) CreateUser(ctx context.Context, body models.CreateUserRequest) (*models.User, error) {
	query := `INSERT INTO users (first_name, last_name, email, password, is_verified, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) RETURNING id, first_name, last_name, email, is_verified, created_at, updated_at`

	var user models.User
	err := r.pool.QueryRow(ctx, query, body.FirstName, body.LastName, body.Email, body.Password, false).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.IsVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return &user, nil
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, first_name, last_name, email, password, is_verified, created_at, updated_at FROM users WHERE email = $1`

	var user models.User
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.IsVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load user: %w", err)
	}
	return &user, nil
}

func (r *AuthRepository) VerifyUser(ctx context.Context, status bool, email string) (*models.User, error) {
	query := `UPDATE users SET is_verified = $1, updated_at = NOW() WHERE email = $2 RETURNING id, first_name, last_name, email, password, is_verified, created_at, updated_at`
	var user models.User
	err := r.pool.QueryRow(ctx, query, status, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.IsVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to verify user: %w", err)
	}
	return &user, nil
}