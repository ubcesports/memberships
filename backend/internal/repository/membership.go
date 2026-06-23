package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ubcesports/memberships/internal/database/db"
)

var ErrActiveMembership = errors.New("user already has an active membership")

type MembershipRepository struct {
	pool  *pgxpool.Pool
	store *db.Queries
}

type PendingTransactionInput struct {
	ID                  uuid.UUID
	UserID              pgtype.UUID
	TierID              pgtype.UUID
	GroupAtPurchase     db.GroupType
	IsStudentAtPurchase bool
	StripePriceID       string
	AmountMinor         int64
	Currency            string
}

func NewMembershipRepository(pool *pgxpool.Pool, store *db.Queries) *MembershipRepository {
	return &MembershipRepository{pool: pool, store: store}
}

func (r *MembershipRepository) ListEligibleTierPriceMappings(ctx context.Context, userID pgtype.UUID) ([]db.ListEligibleTierPriceMappingsRow, error) {
	return r.store.ListEligibleTierPriceMappings(ctx, userID)
}

func (r *MembershipRepository) GetActiveMembership(ctx context.Context, userID pgtype.UUID) (db.GetActiveMembershipByUserIDRow, error) {
	return r.store.GetActiveMembershipByUserID(ctx, userID)
}

func (r *MembershipRepository) GetPendingTransaction(ctx context.Context, userID pgtype.UUID) (db.Transaction, error) {
	return r.store.GetPendingTransactionByUserID(ctx, userID)
}

func (r *MembershipRepository) GetPendingTransactionBySession(ctx context.Context, sessionID string) (db.Transaction, error) {
	return r.store.GetTransactionBySessionIDForUpdate(ctx, pgtype.Text{String: sessionID, Valid: true})
}

func (r *MembershipRepository) GetUserEmail(ctx context.Context, userID pgtype.UUID) (string, error) {
	return r.store.GetUserEmail(ctx, userID)
}

func (r *MembershipRepository) CreateOrGetPendingTransaction(ctx context.Context, input PendingTransactionInput) (db.Transaction, bool, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return db.Transaction{}, false, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	queries := r.store.WithTx(tx)

	if _, err := queries.LockUserForCheckout(ctx, input.UserID); err != nil {
		return db.Transaction{}, false, err
	}
	if _, err := queries.GetActiveMembershipByUserID(ctx, input.UserID); err == nil {
		return db.Transaction{}, false, ErrActiveMembership
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return db.Transaction{}, false, err
	}
	if existing, err := queries.GetPendingTransactionByUserID(ctx, input.UserID); err == nil {
		if err := tx.Commit(ctx); err != nil {
			return db.Transaction{}, false, err
		}
		return existing, true, nil
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return db.Transaction{}, false, err
	}

	created, err := queries.CreatePendingTransaction(ctx, db.CreatePendingTransactionParams{
		ID:                  uuidToPG(input.ID),
		UserID:              input.UserID,
		TierID:              input.TierID,
		GroupAtPurchase:     db.NullGroupType{GroupType: input.GroupAtPurchase, Valid: true},
		IsStudentAtPurchase: input.IsStudentAtPurchase,
		StripePriceID:       pgtype.Text{String: input.StripePriceID, Valid: true},
		AmountMinor:         input.AmountMinor,
		Currency:            pgtype.Text{String: input.Currency, Valid: true},
	})
	if err != nil {
		return db.Transaction{}, false, err
	}
	if err := tx.Commit(ctx); err != nil {
		return db.Transaction{}, false, err
	}
	return created, false, nil
}

func (r *MembershipRepository) AttachCheckoutSession(ctx context.Context, transactionID pgtype.UUID, sessionID string) error {
	return r.store.AttachCheckoutSession(ctx, db.AttachCheckoutSessionParams{
		ID:                      transactionID,
		StripeCheckoutSessionID: pgtype.Text{String: sessionID, Valid: true},
	})
}

func (r *MembershipRepository) FailPendingTransaction(ctx context.Context, transactionID pgtype.UUID) error {
	_, err := r.store.MarkPendingTransactionFailed(ctx, transactionID)
	return err
}

func (r *MembershipRepository) ExpirePendingTransaction(ctx context.Context, transactionID pgtype.UUID) error {
	_, err := r.store.MarkPendingTransactionExpired(ctx, transactionID)
	return err
}

func (r *MembershipRepository) FulfillCheckout(ctx context.Context, sessionID, paymentIntentID, chargeID string, amount int64, currency string, startedAt, expiresAt time.Time) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	queries := r.store.WithTx(tx)

	transaction, err := queries.GetTransactionBySessionIDForUpdate(ctx, pgtype.Text{String: sessionID, Valid: true})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}
	if transaction.Status == db.TransactionStatusTypeCompleted || transaction.Status == db.TransactionStatusTypeRefunded {
		return tx.Commit(ctx)
	}
	if !transaction.TierID.Valid || !transaction.GroupAtPurchase.Valid {
		return errors.New("pending transaction is missing membership purchase data")
	}

	membershipID := uuid.New()
	if _, err := queries.CreateMembership(ctx, db.CreateMembershipParams{
		ID:                  uuidToPG(membershipID),
		UserID:              transaction.UserID,
		TierID:              transaction.TierID,
		GroupAtPurchase:     transaction.GroupAtPurchase.GroupType,
		IsStudentAtPurchase: transaction.IsStudentAtPurchase,
		StartedAt:           pgtype.Timestamptz{Time: startedAt, Valid: true},
		ExpiresAt:           pgtype.Timestamptz{Time: expiresAt, Valid: true},
	}); err != nil {
		return err
	}
	if err := queries.CompleteTransaction(ctx, db.CompleteTransactionParams{
		ID:                    transaction.ID,
		MembershipID:          uuidToPG(membershipID),
		StripePaymentIntentID: pgtype.Text{String: paymentIntentID, Valid: paymentIntentID != ""},
		StripeChargeID:        pgtype.Text{String: chargeID, Valid: chargeID != ""},
		AmountMinor:           amount,
		Currency:              pgtype.Text{String: currency, Valid: currency != ""},
	}); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *MembershipRepository) ApplyFullRefund(ctx context.Context, paymentIntentID, transactionID, chargeID string, cancelledAt time.Time) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	queries := r.store.WithTx(tx)

	var transaction db.Transaction
	found := false
	if paymentIntentID != "" {
		transaction, err = queries.GetTransactionByPaymentIntentForUpdate(ctx, pgtype.Text{String: paymentIntentID, Valid: true})
		found = err == nil
	}
	if (paymentIntentID == "" || errors.Is(err, pgx.ErrNoRows)) && transactionID != "" {
		parsed, parseErr := uuid.Parse(transactionID)
		if parseErr == nil {
			transaction, err = queries.GetTransactionByIDForUpdate(ctx, uuidToPG(parsed))
			found = err == nil
		}
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	if !found {
		return tx.Commit(ctx)
	}
	if transaction.Status == db.TransactionStatusTypeRefunded {
		return tx.Commit(ctx)
	}

	if err := queries.MarkTransactionRefunded(ctx, db.MarkTransactionRefundedParams{
		ID:      transaction.ID,
		Column2: paymentIntentID,
		Column3: chargeID,
	}); err != nil {
		return err
	}
	if transaction.MembershipID.Valid {
		if err := queries.CancelMembership(ctx, db.CancelMembershipParams{
			ID:          transaction.MembershipID,
			CancelledAt: pgtype.Timestamptz{Time: cancelledAt, Valid: true},
		}); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func uuidToPG(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}
