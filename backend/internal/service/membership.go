package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stripe/stripe-go/v85"
	"github.com/ubcesports/memberships/internal/database/db"
	"github.com/ubcesports/memberships/internal/dto"
	"github.com/ubcesports/memberships/internal/repository"
	"github.com/ubcesports/memberships/internal/stripeclient"
)

var (
	ErrMembershipActive  = errors.New("user already has an active membership")
	ErrTierUnavailable   = errors.New("membership tier is unavailable for this user")
	ErrPaymentProcessing = errors.New("a checkout payment is processing")
)

type MembershipService struct {
	repository *repository.MembershipRepository
	stripe     stripeclient.Gateway
	now        func() time.Time
}

type selectedTier struct {
	dto     dto.MembershipTierDTO
	tierID  pgtype.UUID
	group   db.GroupType
	priceID string
}

func NewMembershipService(repository *repository.MembershipRepository, stripeGateway stripeclient.Gateway) *MembershipService {
	return &MembershipService{repository: repository, stripe: stripeGateway, now: time.Now}
}

func (s *MembershipService) ListEligibleTiers(ctx context.Context, userID pgtype.UUID) ([]dto.MembershipTierDTO, error) {
	selected, err := s.selectEligibleTiers(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]dto.MembershipTierDTO, 0, len(selected))
	for _, tier := range selected {
		result = append(result, tier.dto)
	}
	return result, nil
}

func (s *MembershipService) GetActiveMembership(ctx context.Context, userID pgtype.UUID) (*dto.MembershipDTO, error) {
	membership, err := s.repository.GetActiveMembership(ctx, userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &dto.MembershipDTO{
		ID:                  membership.ID.String(),
		TierID:              membership.TierID.String(),
		TierSlug:            membership.TierSlug,
		TierTitle:           membership.TierTitle,
		TierDescription:     textPointer(membership.TierDescription),
		GroupAtPurchase:     dto.GroupType(membership.GroupAtPurchase),
		IsStudentAtPurchase: membership.IsStudentAtPurchase,
		StartedAt:           membership.StartedAt.Time,
		ExpiresAt:           membership.ExpiresAt.Time,
		CancelledAt:         timestampPointer(membership.CancelledAt),
		Status:              "active",
	}, nil
}

func (s *MembershipService) CreateCheckoutSession(ctx context.Context, userID pgtype.UUID, tierID string) (*dto.CheckoutSessionDTO, error) {
	requestedTierID, err := parseUUID(tierID)
	if err != nil {
		return nil, ErrTierUnavailable
	}

	for attempt := 0; attempt < 2; attempt++ {
		if _, err := s.repository.GetActiveMembership(ctx, userID); err == nil {
			return nil, ErrMembershipActive
		} else if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}

		pending, err := s.repository.GetPendingTransaction(ctx, userID)
		if err == nil {
			result, retry, err := s.checkoutForPendingTransaction(ctx, pending)
			if err != nil || !retry {
				return result, err
			}
			continue
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}

		tiers, err := s.selectEligibleTiers(ctx, userID)
		if err != nil {
			return nil, err
		}
		var chosen *selectedTier
		for i := range tiers {
			if tiers[i].tierID == requestedTierID {
				chosen = &tiers[i]
				break
			}
		}
		if chosen == nil {
			return nil, ErrTierUnavailable
		}

		pending, _, err = s.repository.CreateOrGetPendingTransaction(ctx, repository.PendingTransactionInput{
			ID:                  uuid.New(),
			UserID:              userID,
			TierID:              chosen.tierID,
			GroupAtPurchase:     chosen.group,
			IsStudentAtPurchase: chosen.dto.Price.IsStudent,
			StripePriceID:       chosen.priceID,
			AmountMinor:         chosen.dto.Price.AmountMinor,
			Currency:            chosen.dto.Price.Currency,
		})
		if errors.Is(err, repository.ErrActiveMembership) {
			return nil, ErrMembershipActive
		}
		if err != nil {
			return nil, err
		}
		result, retry, err := s.checkoutForPendingTransaction(ctx, pending)
		if err != nil || !retry {
			return result, err
		}
	}
	return nil, errors.New("unable to establish a checkout session")
}

func (s *MembershipService) checkoutForPendingTransaction(ctx context.Context, transaction db.Transaction) (*dto.CheckoutSessionDTO, bool, error) {
	if transaction.StripeCheckoutSessionID.Valid {
		session, err := s.stripe.GetCheckoutSession(ctx, transaction.StripeCheckoutSessionID.String)
		if err != nil {
			return nil, false, fmt.Errorf("retrieve Stripe Checkout Session: %w", err)
		}
		switch session.Status {
		case stripe.CheckoutSessionStatusOpen:
			if _, err := s.stripe.ExpireCheckoutSession(ctx, session.ID); err != nil {
				return nil, false, fmt.Errorf("expire Stripe Checkout Session: %w", err)
			}
			if err := s.repository.ExpirePendingTransaction(ctx, transaction.ID); err != nil {
				return nil, false, err
			}
			return nil, true, nil
		case stripe.CheckoutSessionStatusComplete:
			return nil, false, ErrPaymentProcessing
		case stripe.CheckoutSessionStatusExpired:
			if err := s.repository.ExpirePendingTransaction(ctx, transaction.ID); err != nil {
				return nil, false, err
			}
			return nil, true, nil
		default:
			return nil, false, ErrPaymentProcessing
		}
	}

	if !transaction.StripePriceID.Valid {
		return nil, false, errors.New("pending transaction has no Stripe Price")
	}
	email, err := s.repository.GetUserEmail(ctx, transaction.UserID)
	if err != nil {
		return nil, false, err
	}
	session, err := s.stripe.CreateCheckoutSession(ctx, stripeclient.CheckoutSessionRequest{
		TransactionID: transaction.ID.String(),
		UserID:        transaction.UserID.String(),
		CustomerEmail: email,
		PriceID:       transaction.StripePriceID.String,
	})
	if err != nil {
		return nil, false, fmt.Errorf("create Stripe Checkout Session: %w", err)
	}
	if err := s.repository.AttachCheckoutSession(ctx, transaction.ID, session.ID); err != nil {
		return nil, false, err
	}
	return &dto.CheckoutSessionDTO{URL: session.URL}, false, nil
}

func (s *MembershipService) HandleCheckoutPaid(ctx context.Context, session *stripe.CheckoutSession, occurredAt time.Time) error {
	if session.PaymentStatus != stripe.CheckoutSessionPaymentStatusPaid {
		return nil
	}
	if session.PaymentIntent == nil || session.PaymentIntent.ID == "" {
		return errors.New("paid Checkout Session is missing its PaymentIntent")
	}
	paymentIntent, err := s.stripe.GetPaymentIntent(ctx, session.PaymentIntent.ID)
	if err != nil {
		return fmt.Errorf("retrieve Stripe PaymentIntent: %w", err)
	}
	chargeID := ""
	if paymentIntent.LatestCharge != nil {
		chargeID = paymentIntent.LatestCharge.ID
	}
	expiresAt, err := membershipExpiry(occurredAt)
	if err != nil {
		return err
	}
	return s.repository.FulfillCheckout(
		ctx,
		session.ID,
		paymentIntent.ID,
		chargeID,
		session.AmountTotal,
		string(session.Currency),
		occurredAt,
		expiresAt,
	)
}

func (s *MembershipService) HandleCheckoutFailed(ctx context.Context, sessionID string) error {
	transaction, err := s.repository.GetPendingTransactionBySession(ctx, sessionID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}
	return s.repository.FailPendingTransaction(ctx, transaction.ID)
}

func (s *MembershipService) HandleCheckoutExpired(ctx context.Context, sessionID string) error {
	transaction, err := s.repository.GetPendingTransactionBySession(ctx, sessionID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}
	return s.repository.ExpirePendingTransaction(ctx, transaction.ID)
}

func (s *MembershipService) HandleChargeRefunded(ctx context.Context, charge *stripe.Charge, occurredAt time.Time) error {
	if !isFullRefund(charge) {
		return nil
	}
	paymentIntentID := ""
	transactionID := charge.Metadata["transaction_id"]
	if charge.PaymentIntent != nil {
		paymentIntentID = charge.PaymentIntent.ID
		paymentIntent, err := s.stripe.GetPaymentIntent(ctx, paymentIntentID)
		if err != nil {
			return fmt.Errorf("retrieve refunded Stripe PaymentIntent: %w", err)
		}
		if transactionID == "" {
			transactionID = paymentIntent.Metadata["transaction_id"]
		}
	}
	return s.repository.ApplyFullRefund(ctx, paymentIntentID, transactionID, charge.ID, occurredAt)
}

func (s *MembershipService) selectEligibleTiers(ctx context.Context, userID pgtype.UUID) ([]selectedTier, error) {
	mappings, err := s.repository.ListEligibleTierPriceMappings(ctx, userID)
	if err != nil {
		return nil, err
	}
	chosen := make(map[string]selectedTier)
	for _, mapping := range mappings {
		price, err := s.stripe.GetPrice(ctx, mapping.StripePriceID)
		if err != nil {
			return nil, fmt.Errorf("retrieve Stripe Price %s: %w", mapping.StripePriceID, err)
		}
		if !price.Active || price.Type != stripe.PriceTypeOneTime || price.UnitAmount < 0 {
			continue
		}
		if price.Product == nil || price.Product.ID != mapping.StripeProductID.String {
			return nil, fmt.Errorf("Stripe Price %s does not belong to tier product %s", price.ID, mapping.StripeProductID.String)
		}
		key := mapping.TierID.String()
		current, exists := chosen[key]
		if exists && current.dto.Price.AmountMinor <= price.UnitAmount {
			continue
		}
		chosen[key] = selectedTier{
			tierID:  mapping.TierID,
			group:   mapping.Group,
			priceID: price.ID,
			dto: dto.MembershipTierDTO{
				ID:          key,
				Slug:        mapping.Slug,
				Title:       mapping.Title,
				Description: textPointer(mapping.Description),
				Price: dto.MembershipTierPriceDTO{
					AmountMinor: price.UnitAmount,
					Currency:    string(price.Currency),
					Group:       dto.GroupType(mapping.Group),
					IsStudent:   mapping.IsStudent,
				},
				ExpiresAt: mustMembershipExpiry(s.now()),
			},
		}
	}
	result := make([]selectedTier, 0, len(chosen))
	for _, tier := range chosen {
		result = append(result, tier)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].dto.Title < result[j].dto.Title })
	return result, nil
}

func isFullRefund(charge *stripe.Charge) bool {
	return charge != nil && charge.Refunded
}

func membershipExpiry(purchasedAt time.Time) (time.Time, error) {
	location, err := time.LoadLocation("America/Vancouver")
	if err != nil {
		return time.Time{}, err
	}
	local := purchasedAt.In(location)
	year := local.Year()
	mayFirst := time.Date(year, time.May, 1, 0, 0, 0, 0, location)
	if !local.Before(mayFirst) {
		year++
	}
	return time.Date(year, time.May, 1, 0, 0, 0, 0, location), nil
}

func mustMembershipExpiry(value time.Time) time.Time {
	expiresAt, err := membershipExpiry(value)
	if err != nil {
		panic(err)
	}
	return expiresAt
}

func parseUUID(value string) (pgtype.UUID, error) {
	parsed, err := uuid.Parse(value)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return pgtype.UUID{Bytes: parsed, Valid: true}, nil
}
