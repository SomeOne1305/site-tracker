package repository

import (
	"context"
	"errors"
	"fmt"
	"visit-tracker/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("record not found")

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) GetUser(ctx context.Context, userID string) (*models.User, error) {
	query := `SELECT id, first_name, last_name, email, password, is_verified, created_at, updated_at FROM users WHERE id = $1`
	var retrievedUser models.User
	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&retrievedUser.ID,
		&retrievedUser.FirstName,
		&retrievedUser.LastName,
		&retrievedUser.Email,
		&retrievedUser.Password,
		&retrievedUser.IsVerified,
		&retrievedUser.CreatedAt,
		&retrievedUser.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &retrievedUser, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, body models.UpdateUserRequest, userID string) (*models.User, error) {
	query := `UPDATE users SET first_name = COALESCE($1, first_name), last_name = COALESCE($2, last_name), updated_at = NOW() WHERE id = $3 RETURNING id, first_name, last_name, email, password, is_verified, created_at, updated_at`
	var user models.User
	err := r.pool.QueryRow(ctx, query, body.FirstName, body.LastName, userID).Scan(
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
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context) (bool, error) {
	query := `DELETE FROM users WHERE id = $1 RETURNING id, first_name, last_name, email, password, is_verified, created_at, updated_at`
	err := r.pool.QueryRow(ctx, query)
	if err != nil {
		return false, fmt.Errorf("failed to delete user: %w", err)
	}
	return true, nil
}
