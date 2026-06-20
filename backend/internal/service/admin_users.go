package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ubcesports/memberships/internal/database/db"
	"github.com/ubcesports/memberships/internal/dto"
	"github.com/ubcesports/memberships/internal/repository"
)

type AdminUserFilters struct {
	FullName  string
	StudentID string
	Email     string
	Role      string
	IsStudent *bool
	Group     string
	Limit     int32
	Offset    int32
}

type AdminUserService struct {
	adminUserRepository *repository.AdminUserRepository
}

func NewAdminUserService(adminUserRepository *repository.AdminUserRepository) *AdminUserService {
	return &AdminUserService{adminUserRepository: adminUserRepository}
}

func (s *AdminUserService) GetUsers(ctx context.Context, filters AdminUserFilters) ([]dto.ProfileDTO, error) {
	if filters.Limit <= 0 {
		filters.Limit = 25
	}
	if filters.Limit > 100 {
		filters.Limit = 100
	}
	if filters.Offset < 0 {
		filters.Offset = 0
	}

	params := buildAdminUserQueryParams(filters)
	params.Limit = pgtype.Int4{Int32: filters.Limit, Valid: true}
	params.Offset = pgtype.Int4{Int32: filters.Offset, Valid: true}

	return s.getUsers(ctx, params)
}

func (s *AdminUserService) ExportUsers(ctx context.Context, filters AdminUserFilters) ([]dto.ProfileDTO, error) {
	return s.getUsers(ctx, buildAdminUserQueryParams(filters))
}

func buildAdminUserQueryParams(filters AdminUserFilters) db.GetUsersAdminParams {
	isStudent := pgtype.Bool{}
	if filters.IsStudent != nil {
		isStudent = pgtype.Bool{
			Bool:  *filters.IsStudent,
			Valid: true,
		}
	}

	return db.GetUsersAdminParams{
		FullName: pgtype.Text{
			String: filters.FullName,
			Valid:  filters.FullName != "",
		},
		StudentID: pgtype.Text{
			String: filters.StudentID,
			Valid:  filters.StudentID != "",
		},
		Email: pgtype.Text{
			String: filters.Email,
			Valid:  filters.Email != "",
		},
		Role: db.NullRoleType{
			RoleType: db.RoleType(filters.Role),
			Valid:    filters.Role != "",
		},
		IsStudent: isStudent,
		Group: db.NullGroupType{
			GroupType: db.GroupType(filters.Group),
			Valid:     filters.Group != "",
		},
	}
}

func (s *AdminUserService) getUsers(ctx context.Context, params db.GetUsersAdminParams) ([]dto.ProfileDTO, error) {
	rows, err := s.adminUserRepository.GetUsers(ctx, params)
	if err != nil {
		return nil, err
	}

	users := make([]dto.ProfileDTO, 0, len(rows))
	for _, row := range rows {
		groups := make([]dto.GroupType, 0, len(row.Groups))
		for _, group := range row.Groups {
			groups = append(groups, dto.GroupType(group))
		}

		users = append(users, dto.ProfileDTO{
			ID:                    row.ID.String(),
			Email:                 row.Email,
			StudentID:             textPointer(row.StudentID),
			Role:                  dto.RoleType(row.Role),
			CreatedAt:             row.CreatedAt.Time,
			UpdatedAt:             row.UpdatedAt.Time,
			FullName:              row.FullName,
			EmailVerifiedAt:       timestampPointer(row.EmailVerifiedAt),
			IsStudent:             row.IsStudent,
			OnboardingCompletedAt: timestampPointer(row.OnboardingCompletedAt),
			AvatarURL:             textPointer(row.AvatarUrl),
			Groups:                groups,
		})
	}

	return users, nil
}
