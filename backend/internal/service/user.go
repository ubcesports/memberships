package service

import (
	"context"

	"github.com/ubcesports/memberships/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) DeleteUser(ctx context.Context, userID any) error {
	return s.userRepo.DeleteUser(ctx, userID)
}
