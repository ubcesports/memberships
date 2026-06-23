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
	ErrMembershipActive     = errors.New("active membership does not allow this purchase")
	ErrTierUnavailable      = errors.New("membership tier is unavailable for this user")
	ErrPaymentProcessing    = errors.New("a checkout payment is processing")
	ErrUpgradeNotChargeable = errors.New("upgrade amount must be positive")
)

type MembershipService struct {
	repository *repository.MembershipRepository
	stripe     stripeclient.Gateway
	now        func() time.Time
}

type selectedTier struct {
	dto                 dto.MembershipTierDTO
	tierID              pgtype.UUID
	membershipID        pgtype.UUID
	group               db.GroupType
	priceID             string
	kind                db.TransactionKindType
	creditAmountMinor   int64
	checkoutAmountMinor int64
}

func NewMembershipService(repository *repository.MembershipRepository, stripeGateway stripeclient.Gateway) *MembershipService {
	return &MembershipService{repository: repository, stripe: stripeGateway, now: time.Now}
}

func (s *MembershipService) ListPublicTiers(ctx context.Context) ([]dto.PublicMembershipTierDTO, error) {
	mappings, err := s.repository.ListPublicTierPriceMappings(ctx)
	if err != nil {
		return nil, err
	}
	byTier := make(map[string]*dto.PublicMembershipTierDTO)
	order := make([]string, 0)
	for _, mapping := range mappings {
		price, err := s.validStripePrice(ctx, mapping.StripePriceID, mapping.StripeProductID.String)
		if err != nil {
			return nil, err
		}
		key := mapping.TierID.String()
		tier := byTier[key]
		if tier == nil {
			tier = &dto.PublicMembershipTierDTO{
				ID: key, Slug: mapping.Slug, Title: mapping.Title,
				Description: textPointer(mapping.Description),
				Prices:      make([]dto.MembershipTierPriceDTO, 0, 2),
				ExpiresAt:   mustMembershipExpiry(s.now()),
			}
			byTier[key] = tier
			order = append(order, key)
		}
		tier.Prices = append(tier.Prices, dto.MembershipTierPriceDTO{
			AmountMinor: price.UnitAmount, Currency: string(price.Currency), Group: dto.GroupType(mapping.Group),
		})
	}
	sort.Slice(order, func(i, j int) bool { return byTier[order[i]].Title < byTier[order[j]].Title })
	result := make([]dto.PublicMembershipTierDTO, 0, len(order))
	for _, key := range order {
		result = append(result, *byTier[key])
	}
	return result, nil
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
		ID: membership.ID.String(), TierID: membership.TierID.String(),
		TierSlug: membership.TierSlug, TierTitle: membership.TierTitle,
		TierDescription: textPointer(membership.TierDescription),
		GroupAtPurchase: dto.GroupType(membership.GroupAtPurchase),
		StartedAt:       membership.StartedAt.Time, ExpiresAt: membership.ExpiresAt.Time,
		CancelledAt: timestampPointer(membership.CancelledAt),
	}, nil
}

func (s *MembershipService) CreateCheckoutSession(ctx context.Context, userID pgtype.UUID, tierID string) (*dto.CheckoutSessionDTO, error) {
	requestedTierID, err := parseUUID(tierID)
	if err != nil {
		return nil, ErrTierUnavailable
	}
	for attempt := 0; attempt < 2; attempt++ {
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
			if _, activeErr := s.repository.GetActiveMembership(ctx, userID); activeErr == nil {
				return nil, ErrMembershipActive
			}
			return nil, ErrTierUnavailable
		}
		if isNonChargeableUpgrade(chosen.kind, chosen.checkoutAmountMinor) {
			return nil, ErrUpgradeNotChargeable
		}

		pending, err = s.repository.CreateOrGetPendingTransaction(ctx, repository.PendingTransactionInput{
			ID: uuid.New(), UserID: userID, MembershipID: chosen.membershipID,
			TierID: chosen.tierID, GroupAtPurchase: chosen.group,
			StripePriceID: chosen.priceID, AmountMinor: chosen.checkoutAmountMinor,
			CreditAmountMinor: chosen.creditAmountMinor,
			Currency:          chosen.dto.Price.Currency, Kind: chosen.kind,
		})
		if errors.Is(err, repository.ErrActiveMembership) || errors.Is(err, repository.ErrInvalidUpgrade) {
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
	price, err := s.stripe.GetPrice(ctx, transaction.StripePriceID.String)
	if err != nil {
		return nil, false, fmt.Errorf("retrieve target Stripe Price: %w", err)
	}
	if price.Product == nil {
		return nil, false, errors.New("target Stripe Price has no Product")
	}
	session, err := s.stripe.CreateCheckoutSession(ctx, stripeclient.CheckoutSessionRequest{
		TransactionID: transaction.ID.String(), UserID: transaction.UserID.String(),
		CustomerEmail: email, PriceID: transaction.StripePriceID.String,
		ProductID: price.Product.ID, AmountMinor: transaction.AmountMinor,
		Currency: transaction.Currency.String, IsUpgrade: transaction.Kind == db.TransactionKindTypeUpgrade,
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
	if !checkoutSessionReadyForFulfillment(session.PaymentStatus) {
		return nil
	}
	paymentIntentID := ""
	chargeID := ""
	if session.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
		if paidCheckoutMissingPaymentIntent(session) {
			return errors.New("paid Checkout Session is missing its PaymentIntent")
		}
		if session.PaymentIntent != nil && session.PaymentIntent.ID != "" {
			paymentIntent, err := s.stripe.GetPaymentIntent(ctx, session.PaymentIntent.ID)
			if err != nil {
				return fmt.Errorf("retrieve Stripe PaymentIntent: %w", err)
			}
			paymentIntentID = paymentIntent.ID
			if paymentIntent.LatestCharge != nil {
				chargeID = paymentIntent.LatestCharge.ID
			}
		}
	}
	expiresAt, err := membershipExpiry(occurredAt)
	if err != nil {
		return err
	}
	return s.repository.FulfillCheckout(ctx, session.ID, paymentIntentID, chargeID,
		session.AmountTotal, string(session.Currency), occurredAt, expiresAt)
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

func (s *MembershipService) selectEligibleTiers(ctx context.Context, userID pgtype.UUID) ([]selectedTier, error) {
	mappings, err := s.repository.ListEligibleTierPriceMappings(ctx, userID)
	if err != nil {
		return nil, err
	}
	active, activeErr := s.repository.GetActiveMembership(ctx, userID)
	hasActive := activeErr == nil
	if activeErr != nil && !errors.Is(activeErr, pgx.ErrNoRows) {
		return nil, activeErr
	}
	if hasActive && active.TierSlug != "regular" {
		return []selectedTier{}, nil
	}

	credit := int64(0)
	if hasActive {
		credit, err = s.repository.GetCompletedPaidAmount(ctx, active.ID)
		if err != nil {
			return nil, err
		}
	}
	result := make([]selectedTier, 0, len(mappings))
	for _, mapping := range mappings {
		if hasActive && mapping.Slug != "premium" {
			continue
		}
		price, err := s.validStripePrice(ctx, mapping.StripePriceID, mapping.StripeProductID.String)
		if err != nil {
			return nil, err
		}
		amountDue := price.UnitAmount
		kind := db.TransactionKindTypePurchase
		membershipID := pgtype.UUID{}
		if hasActive {
			kind = db.TransactionKindTypeUpgrade
			membershipID = active.ID
			amountDue = calculateUpgradeAmount(price.UnitAmount, credit)
		}
		result = append(result, selectedTier{
			tierID: mapping.TierID, membershipID: membershipID, group: mapping.Group,
			priceID: price.ID,
			kind:    kind, creditAmountMinor: credit, checkoutAmountMinor: amountDue,
			dto: dto.MembershipTierDTO{
				ID: mapping.TierID.String(), Slug: mapping.Slug, Title: mapping.Title,
				Description: textPointer(mapping.Description),
				Price:       dto.MembershipTierPriceDTO{AmountMinor: price.UnitAmount, Currency: string(price.Currency), Group: dto.GroupType(mapping.Group)},
				IsUpgrade:   hasActive, CreditAmountMinor: credit, AmountDueMinor: amountDue,
				ExpiresAt: mustMembershipExpiry(s.now()),
			},
		})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].dto.Title < result[j].dto.Title })
	return result, nil
}

func calculateUpgradeAmount(targetAmount, creditAmount int64) int64 {
	return targetAmount - creditAmount
}

func isNonChargeableUpgrade(kind db.TransactionKindType, amount int64) bool {
	return kind == db.TransactionKindTypeUpgrade && amount <= 0
}

func checkoutSessionReadyForFulfillment(status stripe.CheckoutSessionPaymentStatus) bool {
	return status == stripe.CheckoutSessionPaymentStatusPaid ||
		status == stripe.CheckoutSessionPaymentStatusNoPaymentRequired
}

func paidCheckoutMissingPaymentIntent(session *stripe.CheckoutSession) bool {
	return session.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid &&
		session.AmountTotal > 0 &&
		(session.PaymentIntent == nil || session.PaymentIntent.ID == "")
}

func (s *MembershipService) validStripePrice(ctx context.Context, priceID, productID string) (*stripe.Price, error) {
	price, err := s.stripe.GetPrice(ctx, priceID)
	if err != nil {
		return nil, fmt.Errorf("retrieve Stripe Price %s: %w", priceID, err)
	}
	if !price.Active || price.Type != stripe.PriceTypeOneTime || price.UnitAmount < 0 {
		return nil, fmt.Errorf("Stripe Price %s is not an active one-time price", priceID)
	}
	if price.Product == nil || price.Product.ID != productID {
		return nil, fmt.Errorf("Stripe Price %s does not belong to tier product %s", priceID, productID)
	}
	return price, nil
}

func membershipExpiry(purchasedAt time.Time) (time.Time, error) {
	location, err := time.LoadLocation("America/Vancouver")
	if err != nil {
		return time.Time{}, err
	}
	local := purchasedAt.In(location)
	year := local.Year()
	if !local.Before(time.Date(year, time.September, 1, 0, 0, 0, 0, location)) {
		year++
	}
	return time.Date(year, time.September, 1, 0, 0, 0, 0, location), nil
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
