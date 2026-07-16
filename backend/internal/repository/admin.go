package repository

import (
	"context"
	"fmt"

	"github.com/ubcesports/memberships/internal/database/db"
)

type AdminRepository struct {
	store *db.Queries
}

func NewAdminRepository(store *db.Queries) *AdminRepository {
	return &AdminRepository{store: store}
}

func (r *AdminRepository) GetUsers(
	ctx context.Context,
	params db.GetUsersAdminParams) ([]db.GetUsersAdminRow, error) {
	rows, err := r.store.GetUsersAdmin(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("query admin users: %w", err)
	}
	return rows, nil
}

func (r *AdminRepository) CountUsers(
	ctx context.Context,
	params db.CountUsersAdminParams,
) (int64, error) {
	count, err := r.store.CountUsersAdmin(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("count admin users: %w", err)
	}
	return count, nil
}

func (r *AdminRepository) CreateAdminAuditLog(ctx context.Context, params db.CreateAdminAuditLogParams) error {
	err := r.store.CreateAdminAuditLog(ctx, params)
	if err != nil {
		return fmt.Errorf("create admin audit log: %w", err)
	}
	return nil
}

func (r *AdminRepository) GetAdminAuditLogs(ctx context.Context, params db.GetAdminAuditLogsParams) ([]db.GetAdminAuditLogsRow, error) {
	logs, err := r.store.GetAdminAuditLogs(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("query admin audit logs: %w", err)
	}
	return logs, nil
}
