package service

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ubcesports/memberships/internal/database/db"
	"github.com/ubcesports/memberships/internal/dto"
	"github.com/ubcesports/memberships/internal/repository"
)

type ProfileService struct {
	profileRepository *repository.ProfileRepository
}

func NewProfileService(profileRepository *repository.ProfileRepository) *ProfileService {
	return &ProfileService{profileRepository: profileRepository}
}

func (s *ProfileService) GetProfileByUserID(ctx context.Context, userID pgtype.UUID) (*dto.ProfileDTO, error) {
	row, err := s.profileRepository.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return buildProfile(row), nil
}

func buildProfile(row db.GetProfileByUserIDRow) *dto.ProfileDTO {
	profile := &dto.ProfileDTO{
		ID:                    row.ID.String(),
		Email:                 row.Email,
		StudentID:             textPointer(row.StudentID),
		Role:                  dto.RoleType(row.Role),
		CreatedAt:             row.CreatedAt.Time,
		UpdatedAt:             row.UpdatedAt.Time,
		FullName:              row.FullName,
		EmailVerifiedAt:       timestampPointer(row.EmailVerifiedAt),
		OnboardingCompletedAt: timestampPointer(row.OnboardingCompletedAt),
		AvatarURL:             textPointer(row.AvatarUrl),
		Groups:                make([]dto.GroupType, 0, len(row.Groups)),
	}

	for _, group := range row.Groups {
		profile.Groups = append(profile.Groups, dto.GroupType(group))
	}

	return profile
}

func textPointer(value pgtype.Text) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}

func timestampPointer(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}
	return &value.Time
}
