package repository

import (
	"context"
	"visit-tracker/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PathRepository struct {
	pool *pgxpool.Pool
}

func NewPathRepository(pool *pgxpool.Pool) *PathRepository {
	return &PathRepository{pool: pool}
}

func (r *PathRepository) TrackVisitor(ctx context.Context, pathID string, ipAddr string, country string, userAgent string, projectID string) (*models.Visitor, error) {
	return nil, nil
}
