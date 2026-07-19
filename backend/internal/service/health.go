package service

import (
	"context"

	"github.com/ubcesports/memberships/internal/repository"
)

type HealthService struct {
	healthRepo *repository.HealthRepository
}

func NewHealthService(healthRepo *repository.HealthRepository) *HealthService {
	return &HealthService{healthRepo: healthRepo}
}

func (s *HealthService) IsDatabaseHealthy(context context.Context) (bool, error) {
	return s.healthRepo.IsDatabaseHealthy(context)
}
