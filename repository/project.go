package repository

import (
	"context"
	"fmt"
	"visit-tracker/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectRepository struct {
	pool *pgxpool.Pool
}

func NewProjectRepository(pool *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{pool: pool}
}

func (r *ProjectRepository) CreateProject(ctx context.Context, body models.CreateProjectRequest, token string, ownerID string) (*models.Project, error) {
	query := `INSERT INTO projects (project_token, name, description, owner_id, created_at, updated_at) VALUES ($1,  $2, $3, $4, NOW(), NOW()) RETURNING id, project_token, name, description,owner_id, created_at, updated_at`
	fmt.Println(ownerID)
	var project models.Project
	err := r.pool.QueryRow(ctx, query, token, body.Name, body.Description, ownerID).Scan(
		&project.ID,
		&project.ProjectToken,
		&project.Name,
		&project.Description,
		&project.OwnerID,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	return &project, nil
}
func (r *ProjectRepository) GetProjects(ctx context.Context, userID string) ([]models.Project, error) {
	query := `SELECT id, project_token, name, description, owner_id, created_at, updated_at FROM projects WHERE owner_id = $1`

	projects := []models.Project{}

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var project models.Project

		err := rows.Scan(
			&project.ID,
			&project.ProjectToken,
			&project.Name,
			&project.Description,
			&project.OwnerID,
			&project.CreatedAt,
			&project.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		projects = append(projects, project)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return projects, nil
}
