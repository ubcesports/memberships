package repository

import (
	"context"

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
