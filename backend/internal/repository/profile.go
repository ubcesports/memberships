package repository

import (
	"context"

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

	return r.store.GetProfileByUserID(ctx, pgUserId)
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

	return r.store.OnboardUserByUserId(ctx, db.OnboardUserByUserIdParams{
		ID:        pgUserId,
		IsStudent: isStudent,
		StudentID: pgtype.Text{
			String: studentId,
			Valid:  studentId != "",
		},
	})
}

func (r *ProfileRepository) EnsureMemberGroupForUser(ctx context.Context, userId string) error {
	pgUserId, err := util.GetValidatedUUID(userId)
	if err != nil {
		return err
	}

	_, err = r.store.EnsureMemberGroupForUser(ctx, pgUserId)
	return err
}
