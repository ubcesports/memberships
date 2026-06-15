package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ubcesports/memberships/internal/database/db"
)

type MembershipRepository struct {
	store *db.Queries
	pool  *pgxpool.Pool
}

type CompletedPurchaseParams struct {
	UserID                  pgtype.UUID
	TierID                  pgtype.UUID
	GroupAtPurchase         db.GroupType
	StudentStatusAtPurchase db.StudentStatusType
	StartedAt               pgtype.Timestamptz
	ExpiresAt               pgtype.Timestamptz
	StripePaymentIntentID   pgtype.Text
	PriceAmount             pgtype.Numeric
}

type CompletedPurchaseResult struct {
	Membership  db.Membership
	Transaction db.Transaction
}

func NewMembershipRepository(store *db.Queries, pool *pgxpool.Pool) *MembershipRepository {
	return &MembershipRepository{store: store, pool: pool}
}

func (r *MembershipRepository) GetUserForMembershipPricing(ctx context.Context, userID pgtype.UUID) (db.GetUserForMembershipPricingRow, error) {
	return r.store.GetUserForMembershipPricing(ctx, userID)
}

func (r *MembershipRepository) ListUserGroups(ctx context.Context, userID pgtype.UUID) ([]db.GroupType, error) {
	return r.store.ListUserGroups(ctx, userID)
}

func (r *MembershipRepository) ListActiveMembershipTiersWithPrices(ctx context.Context, arg db.ListActiveMembershipTiersWithPricesParams) ([]db.ListActiveMembershipTiersWithPricesRow, error) {
	return r.store.ListActiveMembershipTiersWithPrices(ctx, arg)
}

func (r *MembershipRepository) GetTierPriceByCode(ctx context.Context, arg db.GetTierPriceByCodeParams) (db.GetTierPriceByCodeRow, error) {
	return r.store.GetTierPriceByCode(ctx, arg)
}

func (r *MembershipRepository) GetCurrentMembershipByUserID(ctx context.Context, userID pgtype.UUID) (db.GetCurrentMembershipByUserIDRow, error) {
	return r.store.GetCurrentMembershipByUserID(ctx, userID)
}

func (r *MembershipRepository) GetTransactionByStripePaymentIntentID(ctx context.Context, paymentIntentID pgtype.Text) (db.Transaction, error) {
	return r.store.GetTransactionByStripePaymentIntentID(ctx, paymentIntentID)
}

func (r *MembershipRepository) CreateCompletedPurchase(ctx context.Context, arg CompletedPurchaseParams) (CompletedPurchaseResult, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return CompletedPurchaseResult{}, err
	}
	defer tx.Rollback(ctx)

	qtx := r.store.WithTx(tx)

	membership, err := qtx.CreateMembership(ctx, db.CreateMembershipParams{
		UserID:                  arg.UserID,
		TierID:                  arg.TierID,
		GroupAtPurchase:         arg.GroupAtPurchase,
		StudentStatusAtPurchase: arg.StudentStatusAtPurchase,
		StartedAt:               arg.StartedAt,
		ExpiresAt:               arg.ExpiresAt,
	})
	if err != nil {
		return CompletedPurchaseResult{}, err
	}

	transaction, err := qtx.CreateTransaction(ctx, db.CreateTransactionParams{
		UserID:                arg.UserID,
		StripePaymentIntentID: arg.StripePaymentIntentID,
		PriceAmount:           arg.PriceAmount,
		Status:                db.TransactionStatusTypeCompleted,
	})
	if err != nil {
		return CompletedPurchaseResult{}, err
	}

	membership, err = qtx.SetMembershipTransactionID(ctx, db.SetMembershipTransactionIDParams{
		ID:            membership.ID,
		TransactionID: transaction.ID,
	})
	if err != nil {
		return CompletedPurchaseResult{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return CompletedPurchaseResult{}, err
	}

	return CompletedPurchaseResult{
		Membership:  membership,
		Transaction: transaction,
	}, nil
}
