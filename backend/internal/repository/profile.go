package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ubcesports/memberships/internal/database/db"
	"github.com/ubcesports/memberships/internal/util"
)

type ProfileRepository struct {
	store *db.Queries
}

func NewProfileRepository(store *db.Queries) *ProfileRepository {
	return &ProfileRepository{store: store}
}

func (r *ProfileRepository) GetProfileByUserID(ctx context.Context, userId string) (db.GetProfileByUserIDRow, error) {
	// Validate user id
	pgUserId, err := util.GetValidatedUUID(userId)
	if err != nil {
		return db.GetProfileByUserIDRow{}, err
	}

	row, err := r.store.GetProfileByUserID(ctx, pgUserId)
	if err != nil {
		return db.GetProfileByUserIDRow{}, fmt.Errorf("query profile by user ID: %w", err)
	}
	return row, nil
}

func (r *ProfileRepository) OnboardUserByUserId(
	ctx context.Context,
	userId string,
	isStudent bool,
	studentId string,
) error {
	// Validate user id
	pgUserId, err := util.GetValidatedUUID(userId)
	if err != nil {
		return err
	}

	err = r.store.OnboardUserByUserId(ctx, db.OnboardUserByUserIdParams{
		ID:        pgUserId,
		IsStudent: isStudent,
		StudentID: pgtype.Text{
			String: studentId,
			Valid:  studentId != "",
		},
	})
	if err != nil {
		return fmt.Errorf("update user onboarding status: %w", err)
	}
	return nil
}

func (r *ProfileRepository) EnsureMemberGroupForUser(ctx context.Context, userId string) error {
	pgUserId, err := util.GetValidatedUUID(userId)
	if err != nil {
		return err
	}

	_, err = r.store.EnsureMemberGroupForUser(ctx, pgUserId)
	return err
}
