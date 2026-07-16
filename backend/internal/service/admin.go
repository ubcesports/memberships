package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ubcesports/memberships/internal/database/db"
	"github.com/ubcesports/memberships/internal/dto"
	"github.com/ubcesports/memberships/internal/repository"
	"github.com/ubcesports/memberships/internal/util"
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

type AdminAuditLogFilters struct {
	ActorName string
	Limit     int32
	Offset    int32
}

type AdminAuditLogInput struct {
	ActorUserID  string
	Action       string
	TargetUserID string
	Outcome      db.AdminAuditOutcomeType
	RequestID    string
	Description  string
}

type AdminService struct {
	adminRepository *repository.AdminRepository
}

/*
	Public functions
*/

func NewAdminService(adminRepository *repository.AdminRepository) *AdminService {
	return &AdminService{adminRepository: adminRepository}
}

func (s *AdminService) GetUsers(ctx context.Context, filters AdminUserFilters) ([]dto.ProfileDTO, int64, error) {
	if filters.Limit <= 0 {
		filters.Limit = 25
	}
	if filters.Limit > 100 {
		filters.Limit = 100
	}
	if filters.Offset < 0 {
		filters.Offset = 0
	}

	params := buildAdminQueryParams(filters)
	params.Limit = pgtype.Int4{Int32: filters.Limit, Valid: true}
	params.Offset = pgtype.Int4{Int32: filters.Offset, Valid: true}

	total, err := s.adminRepository.CountUsers(ctx, db.CountUsersAdminParams{
		FullName:  params.FullName,
		StudentID: params.StudentID,
		Email:     params.Email,
		Role:      params.Role,
		IsStudent: params.IsStudent,
		Group:     params.Group,
	})
	if err != nil {
		return nil, 0, err
	}

	users, err := s.getUsers(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (s *AdminService) ExportUsers(
	ctx context.Context,
	filters AdminUserFilters,
	actorId string,
	requestId string,
) ([]dto.ProfileDTO, error) {
	users, exportErr := s.getUsers(ctx, buildAdminQueryParams(filters))

	outcome := db.AdminAuditOutcomeTypeSuccess
	description := fmt.Sprintf("Exported %d users", len(users))

	if exportErr != nil {
		outcome = db.AdminAuditOutcomeTypeFailed
		description = "Failed to export users"
	}

	auditErr := s.createAdminAuditLog(ctx, AdminAuditLogInput{
		ActorUserID: actorId,
		Action:      "users.exported",
		Outcome:     outcome,
		RequestID:   requestId,
		Description: description,
	})

	if auditErr != nil {
		if exportErr != nil {
			return nil, errors.Join(exportErr, auditErr)
		}

		// Fail closed: don't return sensitive export data when it could not
		// be recorded in the audit trail.
		return nil, auditErr
	}

	if exportErr != nil {
		return nil, exportErr
	}

	return users, nil
}

func (s *AdminService) GetAdminAuditLogs(ctx context.Context, filters AdminAuditLogFilters) ([]dto.AdminAuditLogResponse, error) {
	// Ensure limit is a proper number. Shouldn't return too many items at once
	if filters.Limit <= 0 {
		filters.Limit = 25
	}
	if filters.Limit > 100 {
		filters.Limit = 100
	}
	if filters.Offset < 0 {
		filters.Offset = 0
	}

	actorName := strings.TrimSpace(filters.ActorName)
	rows, err := s.adminRepository.GetAdminAuditLogs(ctx, db.GetAdminAuditLogsParams{
		ActorName: pgtype.Text{
			String: actorName,
			Valid:  actorName != "",
		},
		Limit:  filters.Limit,
		Offset: filters.Offset,
	})
	if err != nil {
		return nil, err
	}

	logs := make([]dto.AdminAuditLogResponse, 0, len(rows))
	for _, row := range rows {
		var targetUser *dto.AdminAuditLogActor
		if row.TargetID.Valid {
			targetUser = &dto.AdminAuditLogActor{
				ActorUserId:    row.TargetID.String(),
				ActorFullName:  row.TargetName.String,
				ActorAvatarURL: row.TargetAvatarUrl.String,
			}
		}

		logs = append(logs, dto.AdminAuditLogResponse{
			Actor: dto.AdminAuditLogActor{
				ActorUserId:    row.ActorID.String(),
				ActorFullName:  row.ActorName,
				ActorAvatarURL: row.ActorAvatarUrl.String,
			},
			OccuredAt:   row.OccurredAt.Time,
			Action:      row.Action,
			Description: util.TextPointer(row.Description),
			Outcome:     dto.AdminAuditLogOutcomeType(row.Outcome),
			RequestId:   row.RequestID,
			TargetUser:  targetUser,
		})
	}

	return logs, nil
}

/*
	Private functions
*/

func (s *AdminService) createAdminAuditLog(ctx context.Context, input AdminAuditLogInput) error {
	actorID, err := util.GetValidatedUUID(input.ActorUserID)
	if err != nil {
		return fmt.Errorf("invalid audit actor user ID: %w", err)
	}

	targetID := pgtype.UUID{}
	if input.TargetUserID != "" {
		targetID, err = util.GetValidatedUUID(input.TargetUserID)
		if err != nil {
			return fmt.Errorf("invalid audit target user ID: %w", err)
		}
	}

	action := strings.TrimSpace(input.Action)
	if action == "" {
		return fmt.Errorf("audit action is required")
	}
	requestID := strings.TrimSpace(input.RequestID)
	if requestID == "" {
		return fmt.Errorf("audit request ID is required")
	}

	switch input.Outcome {
	case db.AdminAuditOutcomeTypeSuccess,
		db.AdminAuditOutcomeTypeFailed,
		db.AdminAuditOutcomeTypeDenied:
	default:
		return fmt.Errorf("invalid audit outcome: %q", input.Outcome)
	}

	description := strings.TrimSpace(input.Description)
	return s.adminRepository.CreateAdminAuditLog(ctx, db.CreateAdminAuditLogParams{
		ActorUserID:  actorID,
		Action:       action,
		TargetUserID: targetID,
		Outcome:      input.Outcome,
		RequestID:    requestID,
		Description: pgtype.Text{
			String: description,
			Valid:  description != "",
		},
	})
}

func buildAdminQueryParams(filters AdminUserFilters) db.GetUsersAdminParams {
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

func (s *AdminService) getUsers(ctx context.Context, params db.GetUsersAdminParams) ([]dto.ProfileDTO, error) {
	rows, err := s.adminRepository.GetUsers(ctx, params)
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
			StudentID:             util.TextPointer(row.StudentID),
			Role:                  dto.RoleType(row.Role),
			CreatedAt:             row.CreatedAt.Time,
			UpdatedAt:             row.UpdatedAt.Time,
			FullName:              row.FullName,
			EmailVerifiedAt:       util.TimestampPointer(row.EmailVerifiedAt),
			IsStudent:             row.IsStudent,
			OnboardingCompletedAt: util.TimestampPointer(row.OnboardingCompletedAt),
			AvatarURL:             util.TextPointer(row.AvatarUrl),
			Groups:                groups,
		})
	}

	return users, nil
}
