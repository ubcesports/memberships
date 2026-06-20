package repository

import (
	"context"

	"github.com/ubcesports/memberships/internal/database/db"
)

type AdminUserRepository struct {
	store *db.Queries
}

func NewAdminUserRepository(store *db.Queries) *AdminUserRepository {
	return &AdminUserRepository{store: store}
}

func (r *AdminUserRepository) GetUsers(
	ctx context.Context,
	params db.GetUsersAdminParams) ([]db.GetUsersAdminRow, error) {
	return r.store.GetUsersAdmin(ctx, params)
}

func (r *AdminUserRepository) CountUsers(
	ctx context.Context,
	params db.CountUsersAdminParams,
) (int64, error) {
	return r.store.CountUsersAdmin(ctx, params)
}
