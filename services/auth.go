package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"
	"visit-tracker/mailer"
	"visit-tracker/models"
	"visit-tracker/repository"
	"visit-tracker/utils"

	"github.com/redis/go-redis/v9"
)

type AuthService struct {
    repo *repository.AuthRepository
		redis *redis.Client
		mailer *mailer.Mailer
}

func NewAuthService(repo *repository.AuthRepository, redis *redis.Client, mailer *mailer.Mailer) *AuthService {
    return &AuthService{repo: repo, redis: redis, mailer: mailer}
}

var ErrInvalidCredentials = errors.New("invalid credentials")

func (s *AuthService) Register(ctx context.Context, req models.CreateUserRequest) (*models.User, error) {
    u, err := s.repo.GetUserByEmail(ctx, req.Email)
    if err == nil || u != nil {
        return nil, errors.New("email already exists")
    }
    hashedPassword, err := utils.HashPassword(req.Password)
    if err != nil {
        return nil, err
    }
	code, _ := utils.GenerateSixDigitOTP()
    req.Password = hashedPassword
    duration := 15 * time.Minute
	s.redis.Set(ctx, req.Email, code, duration)
	s.mailer.SendMail(req.Email, `Your verification code`, fmt.Sprintf("Your verification code is: %s", code))
    
    return s.repo.CreateUser(ctx, req)
}

func (s *AuthService) VerifyEmail(ctx context.Context, req models.VerifyEmailRequest) (*models.User, error) {
    code, err := s.redis.Get(ctx, req.Email).Result()
    if err != nil {
        return nil, errors.New("invalid verification code")
    }
    if code != strconv.Itoa(req.VerificationCode) {
        return nil, errors.New("invalid verification code")
    }
    return s.repo.VerifyUser(ctx, true, req.Email)
}

func (s *AuthService) Login(ctx context.Context, req models.LoginRequest) (*models.User, error) {
    user, err := s.repo.GetUserByEmail(ctx, req.Email)
    if err != nil {
        return nil, ErrInvalidCredentials
    }

    if !utils.CheckPasswordHash(req.Password, user.Password) {
        return nil, ErrInvalidCredentials
    }

    return user, nil
}

