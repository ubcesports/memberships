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

var (
	ErrActiveMembership    = errors.New("user already has an active membership")
	ErrInvalidUpgrade      = errors.New("membership is no longer eligible for upgrade")
	ErrInvalidReplacement  = errors.New("day pass is no longer eligible for replacement")
	ErrUpgradeWindowClosed = errors.New("membership expires too soon to upgrade")
)

type MembershipRepository struct {
	pool  *pgxpool.Pool
	store *db.Queries
}

type PendingTransactionInput struct {
	ID                uuid.UUID
	UserID            pgtype.UUID
	MembershipID      pgtype.UUID
	TierID            pgtype.UUID
	GroupAtPurchase   db.GroupType
	StripePriceID     string
	AmountMinor       int64
	CreditAmountMinor int64
	Currency          string
	Kind              db.TransactionKindType
	TargetTierSlug    string
	UpgradeWindow     time.Duration
}

func NewMembershipRepository(pool *pgxpool.Pool, store *db.Queries) *MembershipRepository {
	return &MembershipRepository{pool: pool, store: store}
}

func (r *MembershipRepository) ListPublicTierPriceMappings(ctx context.Context) ([]db.ListPublicTierPriceMappingsRow, error) {
	return r.store.ListPublicTierPriceMappings(ctx)
}

func (r *MembershipRepository) ListEligibleTierPriceMappings(ctx context.Context, userID pgtype.UUID) ([]db.ListEligibleTierPriceMappingsRow, error) {
	return r.store.ListEligibleTierPriceMappings(ctx, userID)
}

func (r *MembershipRepository) GetActiveMembership(ctx context.Context, userID pgtype.UUID) (db.GetActiveMembershipByUserIDRow, error) {
	return r.store.GetActiveMembershipByUserID(ctx, userID)
}

func (r *MembershipRepository) GetCompletedPaidAmount(ctx context.Context, membershipID pgtype.UUID) (int64, error) {
	return r.store.GetCompletedPaidAmountForMembership(ctx, membershipID)
}

func (r *MembershipRepository) GetPendingTransaction(ctx context.Context, userID pgtype.UUID) (db.Transaction, error) {
	return r.store.GetPendingTransactionByUserID(ctx, userID)
}

func (r *MembershipRepository) GetPendingTransactionBySession(ctx context.Context, sessionID string) (db.Transaction, error) {
	return r.store.GetTransactionBySessionIDForUpdate(ctx, pgtype.Text{String: sessionID, Valid: true})
}

func (r *MembershipRepository) GetTransactionTierSlugBySession(ctx context.Context, sessionID string) (string, error) {
	return r.store.GetTransactionTierSlugBySessionID(ctx, pgtype.Text{String: sessionID, Valid: true})
}

func (r *MembershipRepository) GetUserEmail(ctx context.Context, userID pgtype.UUID) (string, error) {
	return r.store.GetUserEmail(ctx, userID)
}

func (r *MembershipRepository) CreateOrGetPendingTransaction(ctx context.Context, input PendingTransactionInput) (db.Transaction, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return db.Transaction{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	queries := r.store.WithTx(tx)

	if _, err := queries.LockUserForCheckout(ctx, input.UserID); err != nil {
		return db.Transaction{}, err
	}
	active, activeErr := queries.GetActiveMembershipByUserID(ctx, input.UserID)
	switch input.Kind {
	case db.TransactionKindTypePurchase:
		if activeErr == nil {
			return db.Transaction{}, ErrActiveMembership
		}
	case db.TransactionKindTypeUpgrade:
		if activeErr != nil || active.TierSlug != "regular" || !input.MembershipID.Valid || active.ID != input.MembershipID {
			return db.Transaction{}, ErrInvalidUpgrade
		}
		if !active.ExpiresAt.Time.After(time.Now().Add(input.UpgradeWindow)) {
			return db.Transaction{}, ErrUpgradeWindowClosed
		}
	case db.TransactionKindTypeReplacement:
		if activeErr != nil || active.TierSlug != "day" || !input.MembershipID.Valid || active.ID != input.MembershipID {
			return db.Transaction{}, ErrInvalidReplacement
		}
		if input.TargetTierSlug != "regular" && input.TargetTierSlug != "premium" {
			return db.Transaction{}, ErrInvalidReplacement
		}
	default:
		return db.Transaction{}, errors.New("invalid transaction kind")
	}
	if activeErr != nil && !errors.Is(activeErr, pgx.ErrNoRows) {
		return db.Transaction{}, activeErr
	}
	if existing, err := queries.GetPendingTransactionByUserID(ctx, input.UserID); err == nil {
		if err := tx.Commit(ctx); err != nil {
			return db.Transaction{}, err
		}
		return existing, nil
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return db.Transaction{}, err
	}

	created, err := queries.CreatePendingTransaction(ctx, db.CreatePendingTransactionParams{
		ID:                uuidToPG(input.ID),
		UserID:            input.UserID,
		MembershipID:      input.MembershipID,
		TierID:            input.TierID,
		GroupAtPurchase:   db.NullGroupType{GroupType: input.GroupAtPurchase, Valid: true},
		StripePriceID:     pgtype.Text{String: input.StripePriceID, Valid: true},
		AmountMinor:       input.AmountMinor,
		CreditAmountMinor: input.CreditAmountMinor,
		Currency:          pgtype.Text{String: input.Currency, Valid: true},
		Kind:              input.Kind,
	})
	if err != nil {
		return db.Transaction{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return db.Transaction{}, err
	}
	return created, nil
}

func (r *MembershipRepository) AttachCheckoutSession(ctx context.Context, transactionID pgtype.UUID, sessionID string) error {
	return r.store.AttachCheckoutSession(ctx, db.AttachCheckoutSessionParams{
		ID: transactionID, StripeCheckoutSessionID: pgtype.Text{String: sessionID, Valid: true},
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
	if transaction.Status != db.TransactionStatusTypePending && transaction.Status != db.TransactionStatusTypeFailed {
		return tx.Commit(ctx)
	}
	if !transaction.TierID.Valid || !transaction.GroupAtPurchase.Valid {
		return errors.New("pending transaction is missing membership purchase data")
	}

	membershipID := transaction.MembershipID
	switch transaction.Kind {
	case db.TransactionKindTypeUpgrade:
		if !membershipID.Valid {
			return ErrInvalidUpgrade
		}
		previousExpiresAt, err := queries.CancelMembershipForUpgrade(ctx, db.CancelMembershipForUpgradeParams{
			ID:          membershipID,
			CancelledAt: pgtype.Timestamptz{Time: startedAt, Valid: true},
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrInvalidUpgrade
			}
			return err
		}

		newMembershipID := uuid.New()
		membershipID = uuidToPG(newMembershipID)
		if _, err := queries.CreateMembership(ctx, db.CreateMembershipParams{
			ID: membershipID, UserID: transaction.UserID, TierID: transaction.TierID,
			GroupAtPurchase: transaction.GroupAtPurchase.GroupType,
			StartedAt:       pgtype.Timestamptz{Time: startedAt, Valid: true},
			ExpiresAt:       previousExpiresAt,
		}); err != nil {
			return err
		}
	case db.TransactionKindTypeReplacement:
		if !membershipID.Valid {
			return ErrInvalidReplacement
		}
		cancelled, err := queries.CancelDayMembershipForReplacement(ctx, db.CancelDayMembershipForReplacementParams{
			ID:          membershipID,
			CancelledAt: pgtype.Timestamptz{Time: startedAt, Valid: true},
			UserID:      transaction.UserID,
		})
		if err != nil {
			return err
		}
		if cancelled == 0 {
			return ErrInvalidReplacement
		}

		newMembershipID := uuid.New()
		membershipID = uuidToPG(newMembershipID)
		if _, err := queries.CreateMembership(ctx, db.CreateMembershipParams{
			ID: membershipID, UserID: transaction.UserID, TierID: transaction.TierID,
			GroupAtPurchase: transaction.GroupAtPurchase.GroupType,
			StartedAt:       pgtype.Timestamptz{Time: startedAt, Valid: true},
			ExpiresAt:       pgtype.Timestamptz{Time: expiresAt, Valid: true},
		}); err != nil {
			return err
		}
	case db.TransactionKindTypePurchase:
		newMembershipID := uuid.New()
		membershipID = uuidToPG(newMembershipID)
		if _, err := queries.CreateMembership(ctx, db.CreateMembershipParams{
			ID: membershipID, UserID: transaction.UserID, TierID: transaction.TierID,
			GroupAtPurchase: transaction.GroupAtPurchase.GroupType,
			StartedAt:       pgtype.Timestamptz{Time: startedAt, Valid: true},
			ExpiresAt:       pgtype.Timestamptz{Time: expiresAt, Valid: true},
		}); err != nil {
			return err
		}
	default:
		return errors.New("invalid transaction kind")
	}

	if err := queries.CompleteTransaction(ctx, db.CompleteTransactionParams{
		ID: transaction.ID, MembershipID: membershipID,
		StripePaymentIntentID: pgtype.Text{String: paymentIntentID, Valid: paymentIntentID != ""},
		StripeChargeID:        pgtype.Text{String: chargeID, Valid: chargeID != ""},
		AmountMinor:           amount, Currency: pgtype.Text{String: currency, Valid: currency != ""},
	}); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func uuidToPG(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}
