package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ubcesports/memberships/internal/database/db"
)

type MembershipRepository struct {
	pool  *pgxpool.Pool
	store *db.Queries
}

func NewMembershipRepository(pool *pgxpool.Pool, store *db.Queries) *MembershipRepository {
	return &MembershipRepository{pool: pool, store: store}
}

func (r *MembershipRepository) GetPublicTiersAndPrices(ctx context.Context) ([]db.GetPublicTiersAndPricesRow, error) {
	return r.store.GetPublicTiersAndPrices(ctx)
}
