package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"regexp"
	"strings"

	"github.com/ubcesports/memberships/internal/database/db"
	"github.com/ubcesports/memberships/internal/dto"
	"github.com/ubcesports/memberships/internal/repository"
	"github.com/ubcesports/memberships/internal/util"
)

// Errors
var (
	ErrValidation = errors.New("Validation error")
	ErrConflict   = errors.New("Conflict")
)

var studentIDRegex = regexp.MustCompile(`^\d{8}$`) // Regex to ensure student id is an 8 digit number, which all ubc student ids are

type ProfileService struct {
	profileRepository *repository.ProfileRepository
}

func NewProfileService(profileRepository *repository.ProfileRepository) *ProfileService {
	return &ProfileService{profileRepository: profileRepository}
}

func (s *ProfileService) GetProfileByUserID(ctx context.Context, userID string) (*dto.ProfileDTO, error) {
	row, err := s.profileRepository.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return buildProfile(row), nil
}

func (s *ProfileService) OnboardUser(ctx context.Context, userId string, onboardUserRequest dto.OnboardUserRequest) error {
	var studentID string

	// Get user profile
	user, err := s.GetProfileByUserID(ctx, userId)
	if err != nil {
		return err
	}

	// Only allow non-onboarded users to onboard
	if user.OnboardingCompletedAt != nil {
		return fmt.Errorf("%w: current user is already onboarded!", ErrConflict)
	}

	if onboardUserRequest.IsStudent {

		// Students should provide their id
		if onboardUserRequest.StudentID == nil {
			return fmt.Errorf("%w: student ID is required for students", ErrValidation)
		}

		studentID = strings.TrimSpace(*onboardUserRequest.StudentID)
		if studentID == "" {
			return fmt.Errorf("%w: student ID is required for students", ErrValidation)
		}

		// ID should be 8 digit number
		if !studentIDRegex.MatchString(studentID) {
			return fmt.Errorf("%w: student ID must be an 8 digit number", ErrValidation)
		}
	} else {
		num := rand.IntN(10_000_000)          // 0 to 9999999 (7 digits)
		studentID = fmt.Sprintf("N%07d", num) // Produce a random student id in the format "N1234567"
	}

	return s.profileRepository.OnboardUserByUserId(
		ctx,
		userId,
		onboardUserRequest.IsStudent,
		studentID,
	)
}

/*
	Private functions
*/

func buildProfile(row db.GetProfileByUserIDRow) *dto.ProfileDTO {
	profile := &dto.ProfileDTO{
		ID:                    row.ID.String(),
		Email:                 row.Email,
		StudentID:             util.TextPointer(row.StudentID),
		Role:                  dto.RoleType(row.Role),
		CreatedAt:             row.CreatedAt.Time,
		UpdatedAt:             row.UpdatedAt.Time,
		FullName:              row.FullName,
		EmailVerifiedAt:       util.TimestampPointer(row.EmailVerifiedAt),
		IsStudent:             row.IsStudent,
		OnboardingCompletedAt: util.TimestampPointer(row.OnboardingCompletedAt),
		AvatarURL:             util.TextPointer(row.AvatarUrl),
		Groups:                make([]dto.GroupType, 0, len(row.Groups)),
	}

	for _, group := range row.Groups {
		profile.Groups = append(profile.Groups, dto.GroupType(group))
	}

	return profile
}
