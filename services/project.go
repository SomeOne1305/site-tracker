package services

import (
	"context"
	"visit-tracker/models"
	"visit-tracker/repository"
	"visit-tracker/utils"
)

type ProjectService struct {
	repo repository.ProjectRepository
}

func NewProjectRepository(repo repository.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) CreateProject(ctx context.Context, userID string, body models.CreateProjectRequest) (*models.Project, error) {
	token, _ := utils.GenerateRefreshToken()
	return s.repo.CreateProject(ctx, body, token, userID)
}

func (s *ProjectService) GetProjects(ctx context.Context, userID string) ([]models.Project, error) {
	return s.repo.GetProjects(ctx, userID)
}
