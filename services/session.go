package services

import (
	"context"
	"fmt"
	"net"
	"time"
	"visit-tracker/config"
	"visit-tracker/models"
	"visit-tracker/repository"

	"github.com/google/uuid"
	"github.com/ipinfo/go/v2/ipinfo"
)

type SessionService struct {
	repo *repository.SessionRepository
}

func NewSessionService(repo *repository.SessionRepository) *SessionService {
	return &SessionService{repo: repo}
}

func isLocalIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	return parsedIP.IsLoopback() || parsedIP.IsPrivate()
}

func (s *SessionService) CreateSession(ctx context.Context, userID uuid.UUID, refreshToken string, userAgent string, ipAddress string, expiresAt time.Time) (*models.Session, error) {
	token := config.LoadConfig().IPInfoToken
	client := ipinfo.NewClient(nil, nil, token)
	if isLocalIP(ipAddress) {
		location := "Local IP"
		return s.repo.CreateSession(ctx, userID, refreshToken, userAgent, ipAddress, location, expiresAt)
	}
	resp, err := client.GetIPInfo(net.ParseIP(ipAddress))
	if err != nil {
		fmt.Println("Error fetching IP info:", err)
		return nil, err
	}
	location := resp.City + ", " + resp.Region + ", " + resp.Country

	return s.repo.CreateSession(ctx, userID, refreshToken, userAgent, ipAddress, location, expiresAt)
}

func (s *SessionService) GetSessionById(ctx context.Context, id uuid.UUID) (*models.Session, error) {
	return s.repo.GetSessionById(ctx, id)
}
