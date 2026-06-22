package service

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/ubcesports/memberships/internal/dto"
	"github.com/ubcesports/memberships/internal/repository"
)

type OnboardingService struct {
	onboardingRepo *repository.OnboardingRepository
}

func NewOnboardingService(onboardingRepo *repository.OnboardingRepository) *OnboardingService {
	return &OnboardingService{onboardingRepo: onboardingRepo}
}

// CompleteOnboarding marks a user as onboarded with either a student ID or a generated ID
func (s *OnboardingService) CompleteOnboarding(ctx context.Context, userID string, req *dto.OnboardUserDTO) error {
	studentID := ""

	if req.IsStudent && req.StudentID != nil {
		// User is a student - use provided student ID
		studentID = *req.StudentID
	} else if !req.IsStudent {
		// User is not a student - generate random ID (N#######)
		studentID = generateRandomID()
	}

	return s.onboardingRepo.CompleteOnboarding(ctx, userID, req.IsStudent, studentID)
}

// generateRandomID generates a random ID in format N#######
func generateRandomID() string {
	// Generate 7 random digits
	randomNum := rand.Intn(10000000) // 0-9999999
	return fmt.Sprintf("N%07d", randomNum)
}
