package services

import (
	"context"
	"visit-tracker/models"
	"visit-tracker/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(ctx context.Context, userID string) (*models.User, error) {
	return s.repo.GetUser(ctx, userID)
}

func (s *UserService) UpdateUser(ctx context.Context, body models.UpdateUserRequest, userID string) (*models.User, error) {
	return s.repo.UpdateUser(ctx, body, userID)
}
