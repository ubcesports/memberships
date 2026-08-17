package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ubcesports/memberships/internal/database/db"
	"github.com/ubcesports/memberships/internal/dto"
	"github.com/ubcesports/memberships/internal/repository"
	"github.com/ubcesports/memberships/internal/util"
)

// Postgres error code for a unique constraint violation.
const pgUniqueViolationCode = "23505"

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

// UpdateUserRequest describes the edits an admin wants to apply to a user.
// Every field is optional; only the ones that are set are acted on.
type UpdateUserRequest struct {
	StudentID        *string
	IsStudent        *bool
	GroupsAdd        []db.GroupType
	GroupsRemove     []db.GroupType
	Role             *db.RoleType
	CancelMembership bool
}

// Audit log actions emitted by UpdateUser.
const (
	actionUserUpdated          = "user.updated"
	actionStudentIDUpdated     = "user.student_id.updated"
	actionStudentStatusUpdated = "user.student_status.updated"
	actionRoleUpdated          = "user.role.updated"
	actionGroupAdded           = "user.group.added"
	actionGroupRemoved         = "user.group.removed"
	actionMembershipCancelled  = "user.membership.cancelled"
)

// Number of times a random non-student ID is regenerated before giving up.
const maxNonStudentIDAttempts = 5

// pendingAuditLog is an audit entry that has been earned by a successful action
// but not yet written.
type pendingAuditLog struct {
	action      string
	description string
}

// auditableError attributes a failure to the action that caused it, so a
// "failed" audit entry can name it after the transaction has rolled back.
type auditableError struct {
	action      string
	description string
	err         error
}

func (e *auditableError) Error() string { return e.err.Error() }

func (e *auditableError) Unwrap() error { return e.err }

func auditable(action string, description string, err error) error {
	return &auditableError{action: action, description: description, err: err}
}

type AdminService struct {
	adminRepository repository.AdminStore
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

	auditErr := s.createAdminAuditLog(ctx, s.adminRepository, AdminAuditLogInput{
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

		return nil, auditErr
	}

	if exportErr != nil {
		return nil, exportErr
	}

	return users, nil
}

// GetUserByID returns a single user's profile for the admin detail view.
func (s *AdminService) GetUserByID(ctx context.Context, userId string) (*dto.ProfileDTO, error) {
	row, err := s.adminRepository.GetUserByID(ctx, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: user not found", ErrNotFound)
		}
		return nil, err
	}

	profile := buildAdminUserProfile(row)
	return &profile, nil
}

// GetUserMemberships returns every membership the user has held, newest first.
//
// A user with no memberships is a normal state — they exist from sign up but
// only gain a membership by completing a checkout — so this returns an empty
// slice rather than an error.
func (s *AdminService) GetUserMemberships(ctx context.Context, userId string) ([]dto.MembershipDTO, error) {
	// Distinguishes "user does not exist" from "user has never bought one".
	if _, err := s.adminRepository.GetUserByID(ctx, userId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: user not found", ErrNotFound)
		}
		return nil, err
	}

	rows, err := s.adminRepository.GetUserMemberships(ctx, userId)
	if err != nil {
		return nil, err
	}

	memberships := make([]dto.MembershipDTO, 0, len(rows))
	for _, row := range rows {
		memberships = append(memberships, dto.MembershipDTO{
			ID:          row.ID.String(),
			TierId:      row.TierID.String(),
			StartedAt:   row.StartedAt.Time,
			ExpiresAt:   row.ExpiresAt.Time,
			CancelledAt: util.TimestampPointer(row.CancelledAt),
			Transaction: dto.TransactionDTO{
				ID:              row.TransactionID.String(),
				AmountPaid:      fmt.Sprintf("%.2f", float64(row.AmountPaidCents.Int64)/100),
				Status:          dto.TransactionStatusType(row.Status),
				GroupAtPurchase: dto.GroupType(row.GroupAtPurchase.GroupType),
			},
		})
	}

	return memberships, nil
}

// UpdateUser applies every edit in req to the target user inside a single
// transaction and returns the user's resulting profile.
//
// Each individual change produces its own audit log entry. If any change fails
// the whole update is rolled back and a single "failed" entry naming the
// offending action is written instead.
func (s *AdminService) UpdateUser(
	ctx context.Context,
	actorId string,
	targetUserId string,
	requestId string,
	req UpdateUserRequest,
) (*dto.ProfileDTO, error) {
	var profile *dto.ProfileDTO

	updateErr := s.adminRepository.WithTx(ctx, func(store repository.AdminStore) error {
		user, err := store.GetUserByID(ctx, targetUserId)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return auditable(
					actionUserUpdated,
					"Failed to update user: user not found",
					fmt.Errorf("%w: user not found", ErrNotFound),
				)
			}
			return err
		}

		entries, err := s.applyUserUpdates(ctx, store, user, req)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			if err := s.createAdminAuditLog(ctx, store, AdminAuditLogInput{
				ActorUserID:  actorId,
				Action:       entry.action,
				TargetUserID: targetUserId,
				Outcome:      db.AdminAuditOutcomeTypeSuccess,
				RequestID:    requestId,
				Description:  entry.description,
			}); err != nil {
				return err
			}
		}

		updated, err := store.GetUserByID(ctx, targetUserId)
		if err != nil {
			return err
		}

		updatedProfile := buildAdminUserProfile(updated)
		profile = &updatedProfile
		return nil
	})

	if updateErr != nil {
		action := actionUserUpdated
		description := "Failed to update user"

		var auditErr *auditableError
		if errors.As(updateErr, &auditErr) {
			action = auditErr.action
			description = auditErr.description
		}

		// The transaction is gone, so the failure is recorded through the
		// pooled repository rather than the rolled back one.
		if logErr := s.createAdminAuditLog(ctx, s.adminRepository, AdminAuditLogInput{
			ActorUserID:  actorId,
			Action:       action,
			TargetUserID: targetUserId,
			Outcome:      db.AdminAuditOutcomeTypeFailed,
			RequestID:    requestId,
			Description:  description,
		}); logErr != nil {
			return nil, errors.Join(updateErr, logErr)
		}

		return nil, updateErr
	}

	return profile, nil
}

func (s *AdminService) GetAdminAuditLogs(ctx context.Context, filters AdminAuditLogFilters) ([]dto.AdminAuditLogResponse, int64, error) {
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
	total, err := s.adminRepository.CountAdminAuditLogs(ctx, pgtype.Text{
		String: actorName,
		Valid:  actorName != "",
	})

	rows, err := s.adminRepository.GetAdminAuditLogs(ctx, db.GetAdminAuditLogsParams{
		ActorName: pgtype.Text{
			String: actorName,
			Valid:  actorName != "",
		},
		Limit:  filters.Limit,
		Offset: filters.Offset,
	})
	if err != nil {
		return nil, 0, err
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

	return logs, total, nil
}

func (s *AdminService) ExportAuditLogs(
	ctx context.Context,
	filters AdminAuditLogFilters,
	actorId string,
	requestId string,
) ([]dto.AdminAuditLogResponse, error) {
	logs, exportErr := s.getAdminAuditLogs(ctx, buildAdminAuditLogParams(filters))

	outcome := db.AdminAuditOutcomeTypeSuccess
	description := fmt.Sprintf("Exported %d audit logs", len(logs))

	if exportErr != nil {
		outcome = db.AdminAuditOutcomeTypeFailed
		description = "Failed to export audit logs"
	}

	auditErr := s.createAdminAuditLog(ctx, s.adminRepository, AdminAuditLogInput{
		ActorUserID: actorId,
		Action:      "audit_logs.exported",
		Outcome:     outcome,
		RequestID:   requestId,
		Description: description,
	})

	if auditErr != nil {
		if exportErr != nil {
			return nil, errors.Join(exportErr, auditErr)
		}

		return nil, auditErr
	}

	if exportErr != nil {
		return nil, exportErr
	}

	return logs, nil
}

/*
	Private functions
*/

// applyUserUpdates performs every requested change and returns the audit
// entries the successful ones earned.
func (s *AdminService) applyUserUpdates(
	ctx context.Context,
	store repository.AdminStore,
	user db.GetAdminUserByIDRow,
	req UpdateUserRequest,
) ([]pendingAuditLog, error) {
	entries := make([]pendingAuditLog, 0)

	studentEntries, err := s.applyStudentUpdate(ctx, store, user, req)
	if err != nil {
		return nil, err
	}
	entries = append(entries, studentEntries...)

	roleEntries, err := s.applyRoleUpdate(ctx, store, user, req)
	if err != nil {
		return nil, err
	}
	entries = append(entries, roleEntries...)

	groupEntries, err := s.applyGroupUpdates(ctx, store, user, req)
	if err != nil {
		return nil, err
	}
	entries = append(entries, groupEntries...)

	membershipEntries, err := s.applyMembershipUpdates(ctx, store, user, req)
	if err != nil {
		return nil, err
	}
	entries = append(entries, membershipEntries...)

	return entries, nil
}

func (s *AdminService) applyStudentUpdate(
	ctx context.Context,
	store repository.AdminStore,
	user db.GetAdminUserByIDRow,
	req UpdateUserRequest,
) ([]pendingAuditLog, error) {
	update, err := planStudentUpdate(user, req)
	if err != nil {
		return nil, auditable(update.action, "Failed to update student details: "+err.Error(), err)
	}
	if !update.apply {
		return nil, nil
	}

	studentID := update.studentID
	if update.generate {
		studentID, err = generateUnusedNonStudentID(ctx, store)
		if err != nil {
			return nil, auditable(update.action, "Failed to update student status", err)
		}
	}

	if err := store.UpdateUserStudentInfo(ctx, user.ID.String(), update.isStudent, studentID); err != nil {
		if isUniqueViolation(err) {
			err = fmt.Errorf("%w: student ID %s is already in use", ErrConflict, studentID)
		}
		return nil, auditable(update.action, "Failed to update student details", err)
	}

	currentStudentID := textOrEmpty(user.StudentID)
	description := fmt.Sprintf("Updated student ID from %s to %s", displayValue(currentStudentID), studentID)
	if update.action == actionStudentStatusUpdated {
		description = fmt.Sprintf(
			"Updated student status from %s to %s (student ID %s to %s)",
			studentStatusLabel(user.IsStudent),
			studentStatusLabel(update.isStudent),
			displayValue(currentStudentID),
			studentID,
		)
	}

	return []pendingAuditLog{{action: update.action, description: description}}, nil
}

func (s *AdminService) applyRoleUpdate(
	ctx context.Context,
	store repository.AdminStore,
	user db.GetAdminUserByIDRow,
	req UpdateUserRequest,
) ([]pendingAuditLog, error) {
	role, apply, err := planRoleUpdate(user, req)
	if err != nil {
		return nil, auditable(actionRoleUpdated, "Failed to update role: "+err.Error(), err)
	}
	if !apply {
		return nil, nil
	}

	if err := store.UpdateUserRole(ctx, user.ID.String(), role); err != nil {
		return nil, auditable(actionRoleUpdated, "Failed to update role", err)
	}

	return []pendingAuditLog{{
		action:      actionRoleUpdated,
		description: fmt.Sprintf("Updated role from %s to %s", user.Role, role),
	}}, nil
}

func (s *AdminService) applyGroupUpdates(
	ctx context.Context,
	store repository.AdminStore,
	user db.GetAdminUserByIDRow,
	req UpdateUserRequest,
) ([]pendingAuditLog, error) {
	current := make(map[db.GroupType]struct{}, len(user.Groups))
	for _, group := range user.Groups {
		current[db.GroupType(group)] = struct{}{}
	}

	entries := make([]pendingAuditLog, 0, len(req.GroupsAdd)+len(req.GroupsRemove))

	for _, group := range req.GroupsAdd {
		if !isValidGroup(group) {
			err := fmt.Errorf("%w: invalid group %q", ErrValidation, group)
			return nil, auditable(actionGroupAdded, "Failed to add group: "+err.Error(), err)
		}
		if _, ok := current[group]; ok {
			continue
		}

		if err := store.AddUserGroup(ctx, user.ID.String(), group); err != nil {
			return nil, auditable(actionGroupAdded, fmt.Sprintf("Failed to add group %s", group), err)
		}

		current[group] = struct{}{}
		entries = append(entries, pendingAuditLog{
			action:      actionGroupAdded,
			description: fmt.Sprintf("Added group %s", group),
		})
	}

	for _, group := range req.GroupsRemove {
		// The member group is permanent, so removal requests for it are dropped.
		if group == db.GroupTypeMember {
			continue
		}
		if !isValidGroup(group) {
			err := fmt.Errorf("%w: invalid group %q", ErrValidation, group)
			return nil, auditable(actionGroupRemoved, "Failed to remove group: "+err.Error(), err)
		}
		if _, ok := current[group]; !ok {
			continue
		}

		if err := store.RemoveUserGroup(ctx, user.ID.String(), group); err != nil {
			return nil, auditable(actionGroupRemoved, fmt.Sprintf("Failed to remove group %s", group), err)
		}

		delete(current, group)
		entries = append(entries, pendingAuditLog{
			action:      actionGroupRemoved,
			description: fmt.Sprintf("Removed group %s", group),
		})
	}

	return entries, nil
}

func (s *AdminService) applyMembershipUpdates(
	ctx context.Context,
	store repository.AdminStore,
	user db.GetAdminUserByIDRow,
	req UpdateUserRequest,
) ([]pendingAuditLog, error) {
	userID := user.ID.String()

	if req.CancelMembership {
		hasActive, err := store.HasActiveMembership(ctx, userID)
		if err != nil {
			return nil, auditable(actionMembershipCancelled, "Failed to cancel membership", err)
		}
		if !hasActive {
			err := fmt.Errorf("%w: user has no active membership to cancel", ErrValidation)
			return nil, auditable(actionMembershipCancelled, "Failed to cancel membership: no active membership", err)
		}

		if err := store.CancelActiveMembershipsByUserId(ctx, userID, time.Now()); err != nil {
			return nil, auditable(actionMembershipCancelled, "Failed to cancel membership", err)
		}

		return []pendingAuditLog{{
			action:      actionMembershipCancelled,
			description: "Cancelled the user's active membership",
		}}, nil
	}

	return nil, nil
}

// studentUpdate is the resolved outcome of the student ID and student status
// rules for a single request.
type studentUpdate struct {
	apply     bool
	isStudent bool
	studentID string
	generate  bool // when true, studentID is a randomly generated non-student ID
	action    string
}

func planStudentUpdate(user db.GetAdminUserByIDRow, req UpdateUserRequest) (studentUpdate, error) {
	statusChanged := req.IsStudent != nil && *req.IsStudent != user.IsStudent

	// Becoming a student requires a real student ID to go with it.
	if statusChanged && *req.IsStudent {
		update := studentUpdate{action: actionStudentStatusUpdated}
		if req.StudentID == nil {
			return update, fmt.Errorf("%w: student ID is required when marking a user as a student", ErrValidation)
		}

		studentID := strings.TrimSpace(*req.StudentID)
		if !studentIDRegex.MatchString(studentID) {
			return update, fmt.Errorf("%w: student ID must be an 8 digit number", ErrValidation)
		}

		update.apply = true
		update.isStudent = true
		update.studentID = studentID
		return update, nil
	}

	// Dropping student status replaces the student ID with a generated one, so
	// supplying one alongside would be silently discarded.
	if statusChanged && !*req.IsStudent {
		update := studentUpdate{
			apply:     true,
			isStudent: false,
			generate:  true,
			action:    actionStudentStatusUpdated,
		}
		if req.StudentID != nil {
			return studentUpdate{action: actionStudentStatusUpdated},
				fmt.Errorf("%w: student ID cannot be set when marking a user as a non-student", ErrValidation)
		}
		return update, nil
	}

	if req.StudentID == nil {
		return studentUpdate{}, nil
	}

	update := studentUpdate{action: actionStudentIDUpdated}
	if !user.IsStudent {
		return update, fmt.Errorf("%w: student ID can only be edited for students", ErrValidation)
	}

	studentID := strings.TrimSpace(*req.StudentID)
	if !studentIDRegex.MatchString(studentID) {
		return update, fmt.Errorf("%w: student ID must be an 8 digit number", ErrValidation)
	}
	if studentID == textOrEmpty(user.StudentID) {
		return studentUpdate{}, nil
	}

	update.apply = true
	update.isStudent = true
	update.studentID = studentID
	return update, nil
}

// planRoleUpdate resolves the requested role. An absent role leaves the user
// untouched, while an empty one floors the user back to member.
func planRoleUpdate(user db.GetAdminUserByIDRow, req UpdateUserRequest) (db.RoleType, bool, error) {
	if req.Role == nil {
		return user.Role, false, nil
	}

	role := db.RoleType(strings.TrimSpace(string(*req.Role)))
	switch role {
	case "":
		role = db.RoleTypeMember
	case db.RoleTypeMember, db.RoleTypeAdmin:
	default:
		return user.Role, false, fmt.Errorf("%w: role must be either member or admin", ErrValidation)
	}

	if role == user.Role {
		return role, false, nil
	}
	return role, true, nil
}

func generateUnusedNonStudentID(ctx context.Context, store repository.AdminStore) (string, error) {
	for range maxNonStudentIDAttempts {
		candidate := util.GenerateNonStudentID()

		exists, err := store.StudentIDExists(ctx, candidate)
		if err != nil {
			return "", err
		}
		if !exists {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("unable to generate an unused non-student ID after %d attempts", maxNonStudentIDAttempts)
}

func isValidGroup(group db.GroupType) bool {
	switch group {
	case db.GroupTypeMember,
		db.GroupTypeCompetitiveTeam,
		db.GroupTypeExecutive,
		db.GroupTypeDirector,
		db.GroupTypeBoard:
		return true
	default:
		return false
	}
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolationCode
}

func studentStatusLabel(isStudent bool) string {
	if isStudent {
		return "student"
	}
	return "non-student"
}

func textOrEmpty(value pgtype.Text) string {
	if !value.Valid {
		return ""
	}
	return value.String
}

func displayValue(value string) string {
	if value == "" {
		return "none"
	}
	return value
}

func buildAdminUserProfile(row db.GetAdminUserByIDRow) dto.ProfileDTO {
	groups := make([]dto.GroupType, 0, len(row.Groups))
	for _, group := range row.Groups {
		groups = append(groups, dto.GroupType(group))
	}

	return dto.ProfileDTO{
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
	}
}

func (s *AdminService) createAdminAuditLog(
	ctx context.Context,
	store repository.AdminStore,
	input AdminAuditLogInput,
) error {
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
	return store.CreateAdminAuditLog(ctx, db.CreateAdminAuditLogParams{
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

func buildAdminAuditLogParams(filters AdminAuditLogFilters) db.GetAdminAuditLogsParams {
	actorName := strings.TrimSpace(filters.ActorName)
	return db.GetAdminAuditLogsParams{
		ActorName: pgtype.Text{
			String: actorName,
			Valid:  actorName != "",
		},
		Limit:  filters.Limit,
		Offset: filters.Offset,
	}
}

func (s *AdminService) getAdminAuditLogs(ctx context.Context, params db.GetAdminAuditLogsParams) ([]dto.AdminAuditLogResponse, error) {
	rows, err := s.adminRepository.GetAdminAuditLogs(ctx, params)

	if err != nil {
		return nil, err
	}

	logs := make([]dto.AdminAuditLogResponse, 0, len(rows))
	for _, row := range rows {
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
			TargetUser: &dto.AdminAuditLogActor{
				ActorUserId:    row.TargetID.String(),
				ActorFullName:  row.TargetName.String,
				ActorAvatarURL: row.TargetAvatarUrl.String,
			},
		})
	}

	return logs, nil
}
