package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ubcesports/memberships/internal/database/db"
)

type ProfileRepository struct {
	store *db.Queries
}

func NewProfileRepository(store *db.Queries) *ProfileRepository {
	return &ProfileRepository{store: store}
}

func (r *ProfileRepository) GetProfileByUserID(ctx context.Context, userID pgtype.UUID) (db.GetProfileByUserIDRow, error) {
	return r.store.GetProfileByUserID(ctx, userID)
}
