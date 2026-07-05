package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
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

func (r *MembershipRepository) GetCurrentMembershipWithTransaction(ctx context.Context, userId string) (db.GetCurrentMembershipWithTransactionRow, error) {
	var pgUserId pgtype.UUID

	if err := pgUserId.Scan(userId); err != nil {
		return db.GetCurrentMembershipWithTransactionRow{}, err
	}

	return r.store.GetCurrentMembershipWithTransaction(ctx, pgUserId)
}

func (r *MembershipRepository) GetAllMembershipsWithTransactions(ctx context.Context, userId string) ([]db.GetAllMembershipsWithTransactionsRow, error) {
	var pgUserId pgtype.UUID

	if err := pgUserId.Scan(userId); err != nil {
		return []db.GetAllMembershipsWithTransactionsRow{}, err
	}

	return r.store.GetAllMembershipsWithTransactions(ctx, pgUserId)
}

func (r *MembershipRepository) GetEligibleTiersWithPrices(ctx context.Context, userId string) ([]db.GetEligibleTiersWithPricesRow, error) {
	var pgUserId pgtype.UUID

	if err := pgUserId.Scan(userId); err != nil {
		return []db.GetEligibleTiersWithPricesRow{}, err
	}

	return r.store.GetEligibleTiersWithPrices(ctx, pgUserId)
}

func (r *MembershipRepository) GetTierByTierId(ctx context.Context, tierId string) (db.GetTierByTierIdRow, error) {
	var pgTierId pgtype.UUID

	if err := pgTierId.Scan(tierId); err != nil {
		return db.GetTierByTierIdRow{}, err
	}

	return r.store.GetTierByTierId(ctx, pgTierId)
}
