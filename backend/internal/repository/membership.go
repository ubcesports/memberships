package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ubcesports/memberships/internal/database/db"
	"github.com/ubcesports/memberships/internal/dto"
)

type CreatePendingTransactionParams struct {
	UserId            string
	GroupAtPurchase   dto.GroupType
	StudentAtPurchase bool
	PurchaseType      dto.PurchaseType
}

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

func (r *MembershipRepository) GetPendingTransactionForUpdate(ctx context.Context, userId string) (db.GetPendingTransactionForUpdateRow, error) {
	var pgUserId pgtype.UUID

	if err := pgUserId.Scan(userId); err != nil {
		return db.GetPendingTransactionForUpdateRow{}, err
	}

	return r.store.GetPendingTransactionForUpdate(ctx, pgUserId)
}

func (r *MembershipRepository) ExpirePendingTransactionById(ctx context.Context, transactionId string) error {
	var pgTransactionId pgtype.UUID

	if err := pgTransactionId.Scan(transactionId); err != nil {
		return err
	}

	return r.store.ExpirePendingTransactionById(ctx, pgTransactionId)
}

func (r *MembershipRepository) PutStripeCheckoutSessionId(ctx context.Context, transactionId string, stripeCheckoutSessionId string) error {
	var pgTransactionId pgtype.UUID

	if err := pgTransactionId.Scan(transactionId); err != nil {
		return err
	}

	return r.store.PutStripeCheckoutSessionId(ctx, db.PutStripeCheckoutSessionIdParams{
		ID: pgTransactionId,
		StripeCheckoutSessionID: pgtype.Text{
			String: stripeCheckoutSessionId,
			Valid:  true,
		},
	})
}

func (r *MembershipRepository) UpdateTransactionStatusById(ctx context.Context, transactionId string, status dto.TransactionStatusType) error {
	var pgTransactionId pgtype.UUID

	if err := pgTransactionId.Scan(transactionId); err != nil {
		return err
	}

	return r.store.UpdateTransactionStatusById(ctx, db.UpdateTransactionStatusByIdParams{
		ID:     pgTransactionId,
		Status: db.TransactionStatusType(status),
	})
}

func (r *MembershipRepository) CreatePendingTransaction(ctx context.Context, params CreatePendingTransactionParams) (string, error) {
	var userID pgtype.UUID
	if err := userID.Scan(params.UserId); err != nil {
		return "", err
	}

	dbParams := db.CreatePendingTransactionParams{
		UserID: userID,
		GroupAtPurchase: db.NullGroupType{
			GroupType: db.GroupType(params.GroupAtPurchase),
			Valid:     true,
		},
		StudentAtPurchase: pgtype.Bool{
			Bool:  params.StudentAtPurchase,
			Valid: true,
		},
		PurchaseType: db.NullPurchaseType{
			PurchaseType: db.PurchaseType(params.PurchaseType),
			Valid:        true,
		},
	}

	id, err := r.store.CreatePendingTransaction(ctx, dbParams)
	if err != nil {
		return "", err
	}

	return uuid.UUID(id.Bytes).String(), nil
}

func (r *MembershipRepository) WithTx(ctx context.Context, fn func(*MembershipRepository) error) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	txRepo := &MembershipRepository{
		pool:  r.pool,
		store: r.store.WithTx(tx),
	}

	if err := fn(txRepo); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
