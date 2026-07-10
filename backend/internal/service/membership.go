package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"slices"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stripe/stripe-go/v86"
	"github.com/ubcesports/memberships/internal/database/db"
	"github.com/ubcesports/memberships/internal/dto"
	"github.com/ubcesports/memberships/internal/repository"
	"github.com/ubcesports/memberships/internal/stripeclient"
)

type MembershipService struct {
	membershipRepo *repository.MembershipRepository
	stripeClient   *stripeclient.Client
	profileService *ProfileService
}

func NewMembershipService(membershipRepo *repository.MembershipRepository, stripeClient *stripeclient.Client, profileService *ProfileService) *MembershipService {
	return &MembershipService{membershipRepo: membershipRepo, stripeClient: stripeClient, profileService: profileService}
}

/*
	Consts
*/

var (
	ErrMembershipAlreadyExists  = errors.New("An active un-upgradeable membership already exists! Can't create new checkout session.")
	ErrTierNotEligible          = errors.New("Requested membership tier not eligible for current user. Please try another value.")
	ErrTierNotFound             = errors.New("Tier with given tier id not found.")
	ErrMembershipPurchaseClosed = errors.New("Membership purchases are closed until the next membership period.")
)

/*
	Public functions
*/

func (s *MembershipService) GetPublicTiersAndPrices(ctx context.Context) ([]dto.MembershipTierDTO, error) {
	tiers, err := s.membershipRepo.GetPublicTiersAndPrices(ctx)
	if err != nil {
		return nil, err
	}

	returnTiers := make([]dto.MembershipTierDTO, 0, len(tiers))
	tierIndexById := make(map[string]int)

	for _, tier := range tiers {
		tierId := tier.ID.String()

		// Get price from stripe price id
		price, err := s.stripeClient.GetPrice(ctx, tier.StripePriceID.String)
		if err != nil {
			return nil, err
		}

		// Set up price dto
		var isStudentRequired *bool
		if tier.IsStudentRequired.Valid {
			isStudentRequired = &tier.IsStudentRequired.Bool
		} else {
			isStudentRequired = nil
		}

		priceDto := dto.MembershipTierPriceDTO{
			Price:             float64(price.UnitAmount) / 100, // Turn unit amount which is in cents, into readable format with 2 numbers after the decimal
			PriceId:           tier.StripePriceID.String,
			IsStudentRequired: isStudentRequired,
		}

		if index, exists := tierIndexById[tierId]; exists {
			// If tier already exists in return body, append price to it
			returnTiers[index].Prices = append(returnTiers[index].Prices, priceDto)
		} else {
			// If tier doesn't exist in return body, add a new tier with price
			tierIndexById[tierId] = len(returnTiers)
			returnTiers = append(returnTiers, dto.MembershipTierDTO{
				ID:          tier.ID.String(),
				Title:       tier.Title,
				Description: tier.Description.String,
				Slug:        tier.Slug.String,
				ProductId:   tier.StripeProductID.String,
				Prices:      []dto.MembershipTierPriceDTO{priceDto},
			})
		}
	}

	return returnTiers, nil
}

func (s *MembershipService) GetCurrentMembershipWithTransaction(ctx context.Context, userId string) (*dto.MembershipDTO, error) {
	membership, err := s.membershipRepo.GetCurrentMembershipWithTransaction(ctx, userId)
	if err != nil {
		// If user has no current membership, return nil
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &dto.MembershipDTO{
		ID:          membership.ID.String(),
		TierId:      membership.TierID.String(),
		StartedAt:   membership.StartedAt.Time,
		ExpiresAt:   membership.ExpiresAt.Time,
		CancelledAt: &membership.CancelledAt.Time,
		Transaction: dto.TransactionDTO{
			ID:              membership.TransactionID.String(),
			AmountPaid:      fmt.Sprintf("%.2f", float64(membership.AmountPaidCents.Int64)/100),
			Status:          dto.TransactionStatusType(membership.Status),
			GroupAtPurchase: dto.GroupType(membership.GroupAtPurchase.GroupType),
		},
	}, nil
}

func (s *MembershipService) GetAllMembershipsWithTransactions(ctx context.Context, userId string) (*[]dto.MembershipDTO, error) {
	memberships, err := s.membershipRepo.GetAllMembershipsWithTransactions(ctx, userId)
	if err != nil {
		// If user has no current membership, return nil
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	returnMemberships := make([]dto.MembershipDTO, 0, len(memberships))
	for _, membership := range memberships {
		membershipDto := dto.MembershipDTO{
			ID:          membership.ID.String(),
			TierId:      membership.TierID.String(),
			StartedAt:   membership.StartedAt.Time,
			ExpiresAt:   membership.ExpiresAt.Time,
			CancelledAt: &membership.CancelledAt.Time,
			Transaction: dto.TransactionDTO{
				ID:              membership.TransactionID.String(),
				AmountPaid:      fmt.Sprintf("%.2f", float64(membership.AmountPaidCents.Int64)/100),
				Status:          dto.TransactionStatusType(membership.Status),
				GroupAtPurchase: dto.GroupType(membership.GroupAtPurchase.GroupType),
			},
		}
		returnMemberships = append(returnMemberships, membershipDto)
	}

	return &returnMemberships, nil
}

func (s *MembershipService) GetEligibleTiersWithPrices(ctx context.Context, userId string) (*[]dto.EligibleMembershipTierDTO, error) {
	tiers, err := s.membershipRepo.GetEligibleTiersWithPrices(ctx, userId)
	if err != nil {
		return nil, err
	}

	returnTiers := make([]dto.EligibleMembershipTierDTO, 0, len(tiers))
	tierIndexById := make(map[string]int)

	// Get user info
	user, err := s.profileService.GetProfileByUserID(ctx, userId)
	if err != nil {
		return nil, err
	}

	// Get current membership info
	currMembership, err := s.GetCurrentMembershipWithTransaction(ctx, userId)
	if err != nil {
		return nil, err
	}
	var currTierSlug string
	if currMembership != nil {
		tier, err := s.getTierByTierId(ctx, currMembership.TierId)
		if err != nil {
			return nil, err
		}

		currTierSlug = tier.Slug
	}

	// Store slugs for eligible tiers
	eligibleSlugs := make(map[string]dto.PurchaseType)

	addPriceToTier := func(tier db.GetEligibleTiersWithPricesRow, purchaseType dto.PurchaseType, priceDto dto.MembershipTierPriceDTO) {
		tierId := tier.ID.String()

		if _, exists := tierIndexById[tierId]; exists {
			return
		}

		tierIndexById[tierId] = len(returnTiers)
		returnTiers = append(returnTiers, dto.EligibleMembershipTierDTO{
			ID:           tier.ID.String(),
			Title:        tier.Title,
			Description:  tier.Description.String,
			Slug:         tier.Slug.String,
			PurchaseType: purchaseType,
			ProductId:    tier.StripeProductID.String,
			Price:        priceDto,
		})
	}

	// 1. Users in executive, director, and board groups should not be able to see any other tier other than executive tier
	// 	  If they already have an active executive membership, they should not be able to see any eligible tiers
	//    If a user is in an exec and a competitive player, prioritize exec membership over comp membership
	if slices.Contains(user.Groups, dto.GroupExecutive) ||
		slices.Contains(user.Groups, dto.GroupDirector) ||
		slices.Contains(user.Groups, dto.GroupBoard) {

		if currMembership != nil {
			// If user has a current membership and its tier slug is "exec" return no eligible memberships
			if currTierSlug == "executive" {
				return nil, nil
			}

			return nil, fmt.Errorf("exec/director/board user has unexpected active membership tier: %s", currTierSlug)
		} else {
			eligibleSlugs["executive"] = dto.PurchaseNew
		}
	} else if slices.Contains(user.Groups, dto.GroupCompetitiveTeam) {
		// 2. Users in competitive_team group should not be able to see any other tier other than competitive team tier
		//	  If they already have an active competitive team membership, they should not be able to see any eligible tiers

		if currMembership != nil {
			// If user has a current membership and its tier slug is "comp" return no eligible memberships
			if currTierSlug == "competitive_team" {
				return nil, nil
			}

			return nil, fmt.Errorf("competitive team user has unexpected active membership tier: %s", currTierSlug)
		} else {
			eligibleSlugs["competitive_team"] = dto.PurchaseNew
		}
	} else {
		// 3. All other users should see the day pass, regular and premium memberships
		// CASES
		// If user has day pass, return regular/premium pass which they have to pay full price for
		// If user has regular pass, return premium pass which they pay the difference in price for
		// If user has premium pass, return no eligible tiers

		if currMembership != nil {
			switch currTierSlug {
			// If user has a current membership and its tier slug is "day" return regular/premium memberships, which would be a replacement
			case "day":
				eligibleSlugs["regular"] = dto.PurchaseNew
				eligibleSlugs["premium"] = dto.PurchaseNew
			case "regular":
				// If user has a current membership and its tier slug is "regular" return premium membership, which would be an upgrade
				eligibleSlugs["premium"] = dto.PurchaseUpgrade
			case "premium":
				// If user has a current membership and its tier slug is "premium" return no eligible memberships
				return nil, nil
			default:
				return nil, fmt.Errorf("unsupported current membership tier slug: %s", currTierSlug)
			}
		} else {
			eligibleSlugs["day"] = dto.PurchaseNew
			eligibleSlugs["regular"] = dto.PurchaseNew
			eligibleSlugs["premium"] = dto.PurchaseNew
		}
	}

	for _, tier := range tiers {
		purchaseType, ok := eligibleSlugs[tier.Slug.String]
		if !ok {
			continue
		}

		if tier.IsStudentRequired.Valid && tier.IsStudentRequired.Bool != user.IsStudent {
			continue
		}

		switch purchaseType {
		case dto.PurchaseNew:
			// Get price from stripe price id
			price, err := s.stripeClient.GetPrice(ctx, tier.StripePriceID.String)
			if err != nil {
				return nil, err
			}

			// Set up price dto
			priceDto := dto.MembershipTierPriceDTO{
				Price:             float64(price.UnitAmount) / 100, // Turn unit amount which is in cents, into readable format with 2 numbers after the decimal
				PriceId:           tier.StripePriceID.String,
				IsStudentRequired: nil, // Leave nil as this is not really required in this context.
			}

			addPriceToTier(tier, purchaseType, priceDto)
		case dto.PurchaseUpgrade:
			if currMembership == nil {
				return nil, fmt.Errorf("cannot calculate upgrade price without current membership")
			}

			// Get price from stripe price id
			tierPrice, err := s.stripeClient.GetPrice(ctx, tier.StripePriceID.String)
			if err != nil {
				return nil, err
			}

			// Calculate upgrade price to pay
			amountPaidFloat, err := strconv.ParseFloat(currMembership.Transaction.AmountPaid, 64)
			if err != nil {
				return nil, err
			}

			amountPaidInCents := int64(math.Round(amountPaidFloat * 100))
			priceToPay := tierPrice.UnitAmount - amountPaidInCents
			if priceToPay < 0 {
				priceToPay = 0 // Ensure negative price cannot be paid
			}

			// Set up price dto
			priceDto := dto.MembershipTierPriceDTO{
				Price:             float64(priceToPay) / 100,
				PriceId:           tier.StripePriceID.String,
				IsStudentRequired: nil, // Leave nil as this is not really required in this context.
			}

			addPriceToTier(tier, purchaseType, priceDto)
		}
	}

	return &returnTiers, nil
}

func (s *MembershipService) CreateCheckoutSession(ctx context.Context, userId string, req dto.CheckoutSessionRequest) (*dto.CheckoutSessionResponse, error) {
	// 1. Ensure user isn't trying to purchase a meaningless membership too late in the membership period
	isClosed, err := membershipPurchaseClosedAt(time.Now())
	if err != nil {
		return nil, err
	}
	if isClosed {
		return nil, ErrMembershipPurchaseClosed
	}

	reqTier, err := s.getTierByTierId(ctx, req.TierId)
	if err != nil {
		return nil, err
	}

	// 2. Check if an active membership already exists. If if does, return an error.
	membership, err := s.GetCurrentMembershipWithTransaction(ctx, userId)
	if err != nil {
		return nil, err
	}

	// If the user already has an active membership, only allow checkout sessions
	// for supported upgrade paths. All other purchases are rejected to prevent
	// multiple active memberships.
	if membership != nil {
		currTier, err := s.getTierByTierId(ctx, membership.TierId)
		if err != nil {
			return nil, err
		}

		// Users can create an "upgrade" checkout session if they have a regular membership and want to buy premium
		// OR they have a day pass and want to buy regular/premium
		// Return ErrMembershipAlreadyExists if they don't meet the above criteria
		if !((currTier.Slug == "regular" && reqTier.Slug == "premium") ||
			(currTier.Slug == "day" && (reqTier.Slug == "regular" || reqTier.Slug == "premium"))) {
			return nil, ErrMembershipAlreadyExists
		}
	}

	// 3. Check if the requested tier is eligible for the user. If not, return an error.
	eligibleTiers, err := s.GetEligibleTiersWithPrices(ctx, userId)
	if err != nil {
		return nil, err
	}

	var selectedTier *dto.EligibleMembershipTierDTO
	if eligibleTiers != nil {
		for i := range *eligibleTiers {
			tier := &(*eligibleTiers)[i]
			if req.TierId == tier.ID {
				selectedTier = tier
				break
			}
		}
	}

	if selectedTier == nil {
		return nil, ErrTierNotEligible
	}

	// 4. If there is a pending transaction, then expire it and its stripe checkout session
	err = s.membershipRepo.WithTx(ctx, func(mr *repository.MembershipRepository) error {
		pending, err := mr.GetPendingTransactionForUpdate(ctx, userId)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		if err != nil {
			return err
		}

		if !pending.StripeCheckoutSessionID.Valid || pending.StripeCheckoutSessionID.String == "" {
			return mr.ExpirePendingTransactionById(ctx, pending.ID.String())
		}

		_, err = s.stripeClient.ExpireCheckoutSession(ctx, pending.StripeCheckoutSessionID.String)
		if err != nil {
			session, getErr := s.stripeClient.GetCheckoutSession(ctx, pending.StripeCheckoutSessionID.String)
			if getErr != nil {
				return err
			}
			if session.Status != stripe.CheckoutSessionStatusExpired {
				return err
			}
		}

		return mr.ExpirePendingTransactionById(ctx, pending.ID.String())
	})
	if err != nil {
		return nil, err
	}

	// 5. Create new pending transaction and checkout session

	// Get user profile to create checkout session with their email
	profile, err := s.profileService.GetProfileByUserID(ctx, userId)
	if err != nil {
		return nil, err
	}

	// Create new pending transaction
	transactionId, err := s.membershipRepo.CreatePendingTransaction(ctx, repository.CreatePendingTransactionParams{
		UserId:            userId,
		TierId:            selectedTier.ID,
		GroupAtPurchase:   getGroupAtPurchase(profile.Groups),
		StudentAtPurchase: profile.IsStudent,
		PurchaseType:      selectedTier.PurchaseType,
	})
	if err != nil {
		return nil, err
	}

	// Create stripe checkout session
	session, err := s.stripeClient.CreateCheckoutSession(ctx, stripeclient.CheckoutSessionRequest{
		TransactionID: transactionId,
		UserID:        userId,
		CustomerEmail: profile.Email,
		PriceID:       selectedTier.Price.PriceId,
		ProductID:     selectedTier.ProductId,
		AmountInCents: int64(math.Round(selectedTier.Price.Price * 100)),
		Currency:      "cad", // Always in canadian dollars
		IsUpgrade:     selectedTier.PurchaseType == dto.PurchaseUpgrade,
	})
	if err != nil {
		markFailedErr := s.membershipRepo.UpdateTransactionStatusById(ctx, transactionId, dto.TransactionFailed)
		if markFailedErr != nil {
			return nil, fmt.Errorf("create checkout session failed: %w; also failed to mark transaction failed: %v", err, markFailedErr)
		}
		return nil, err
	}

	// Put strile checkout session id into pending transaction
	err = s.membershipRepo.PutStripeCheckoutSessionId(ctx, transactionId, session.ID)
	if err != nil {
		_, expireErr := s.stripeClient.ExpireCheckoutSession(ctx, session.ID)
		markFailedErr := s.membershipRepo.UpdateTransactionStatusById(ctx, transactionId, dto.TransactionFailed)
		if expireErr != nil || markFailedErr != nil {
			return nil, fmt.Errorf("save checkout session id failed: %w; expire checkout session failed: %v; mark transaction failed failed: %v", err, expireErr, markFailedErr)
		}
		return nil, err
	}

	return &dto.CheckoutSessionResponse{Url: session.URL}, nil
}

/*
	Stripe webhook callback functions
*/

func (s *MembershipService) HandleCheckoutPaid(ctx context.Context, session *stripe.CheckoutSession, occurredAt time.Time) error {
	return s.membershipRepo.WithTx(ctx, func(mr *repository.MembershipRepository) error {
		// 1. Lock the transaction for this checkout session so duplicate webhooks cannot fulfill it twice.
		transaction, err := mr.GetTransactionByCheckoutSessionIdForUpdate(ctx, session.ID)
		if err != nil {
			return err
		}

		// 2. Make the handler idempotent: Stripe can retry or duplicate webhook delivery.
		if transaction.Status == db.TransactionStatusTypeCompleted {
			return nil
		}
		if transaction.Status != db.TransactionStatusTypePending {
			return nil
		}

		// 3. Confirm Stripe says this checkout is paid before fulfilling the membership.
		if session.PaymentStatus != stripe.CheckoutSessionPaymentStatusPaid {
			return fmt.Errorf("checkout session %s is not paid", session.ID)
		}

		// 4. Cross-check Stripe metadata against the locked DB transaction.
		if session.Metadata["transaction_id"] != transaction.ID.String() {
			return fmt.Errorf("checkout session %s transaction metadata mismatch", session.ID)
		}
		if session.Metadata["user_id"] != transaction.UserID.String() {
			return fmt.Errorf("checkout session %s user metadata mismatch", session.ID)
		}

		// 5. Cancel any old active membership before creating the new fulfilled membership.
		if err := mr.CancelActiveMembershipsByUserId(ctx, transaction.UserID.String(), occurredAt); err != nil {
			return err
		}

		// 6. Create the fulfilled membership. Memberships expire after April 30 in Vancouver time.
		expiresAt, err := membershipExpiresAt(occurredAt)
		if err != nil {
			return err
		}
		membershipId, err := mr.CreateMembership(ctx, repository.CreateMembershipParams{
			UserId:    transaction.UserID.String(),
			TierId:    transaction.TierID.String(),
			StartedAt: occurredAt,
			ExpiresAt: expiresAt,
		})
		if err != nil {
			return err
		}

		// 7. Record payment details and mark the transaction completed.
		var paymentIntentId string
		if session.PaymentIntent != nil {
			paymentIntentId = session.PaymentIntent.ID
		}
		return mr.CompleteTransaction(ctx, repository.CompleteTransactionParams{
			TransactionId:         transaction.ID.String(),
			MembershipId:          membershipId,
			StripePaymentIntentId: paymentIntentId,
			AmountPaidCents:       session.AmountTotal,
		})
	})
}

func (s *MembershipService) HandleCheckoutExpired(ctx context.Context, sessionId string) error {
	return s.membershipRepo.UpdatePendingTransactionStatusByCheckoutId(ctx, sessionId, dto.TransactionExpired)
}

func (s *MembershipService) HandleCheckoutFailed(ctx context.Context, sessionId string) error {
	return s.membershipRepo.UpdatePendingTransactionStatusByCheckoutId(ctx, sessionId, dto.TransactionFailed)
}

/*
	Private functions
*/

// memberships follow the UBC Esports membership year, which runs from
// May 1 00:00:00 to April 30 23:59:59 (America/Vancouver).
//
// Expiry rules:
//   - Purchases made from January 1 through April 30 expire on
//     April 30 23:59:59 of the same calendar year.
//   - Purchases made on or after May 1 expire on
//     April 30 23:59:59 of the following calendar year.
//
// All calculations are performed in the America/Vancouver time zone,
// regardless of the purchaser's local time zone.
func membershipExpiresAt(purchasedAt time.Time) (time.Time, error) {
	location, err := time.LoadLocation("America/Vancouver")
	if err != nil {
		return time.Time{}, err
	}

	localPurchasedAt := purchasedAt.In(location)

	expiryYear := localPurchasedAt.Year()

	// Membership year rolls over at May 1 00:00:00.
	cutoff := time.Date(expiryYear, time.May, 1, 0, 0, 0, 0, location)

	if !localPurchasedAt.Before(cutoff) {
		expiryYear++
	}

	return time.Date(
		expiryYear,
		time.April,
		30,
		23, 59, 59,
		0,
		location,
	), nil
}

func membershipPurchaseClosedAt(now time.Time) (bool, error) {
	location, err := time.LoadLocation("America/Vancouver")
	if err != nil {
		return false, err
	}

	localNow := now.In(location)

	closedFrom := time.Date(localNow.Year(), time.April, 25, 0, 0, 0, 0, location)
	closedUntil := time.Date(localNow.Year(), time.May, 1, 0, 0, 0, 0, location)

	return !localNow.Before(closedFrom) && localNow.Before(closedUntil), nil
}

// returns the highest priority group a user belongs to at the time of membership purchase.
//
// Group priority (highest to lowest):
//   - Board
//   - Director
//   - Executive
//   - Competitive Team
//   - Member (default)
//
// If a user belongs to multiple groups, the highest priority group is used.
func getGroupAtPurchase(groups []dto.GroupType) dto.GroupType {
	if slices.Contains(groups, dto.GroupBoard) {
		return dto.GroupBoard
	}
	if slices.Contains(groups, dto.GroupDirector) {
		return dto.GroupDirector
	}
	if slices.Contains(groups, dto.GroupExecutive) {
		return dto.GroupExecutive
	}
	if slices.Contains(groups, dto.GroupCompetitiveTeam) {
		return dto.GroupCompetitiveTeam
	}
	return dto.GroupMember
}

func (s *MembershipService) getTierByTierId(ctx context.Context, tierId string) (*dto.MembershipTierDTO, error) {
	tier, err := s.membershipRepo.GetTierByTierId(ctx, tierId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTierNotFound
		}

		return nil, err
	}

	return &dto.MembershipTierDTO{
		ID:          tier.ID.String(),
		Title:       tier.Title,
		Description: tier.Description.String,
		Slug:        tier.Slug.String,
		ProductId:   tier.StripeProductID.String,
		Prices:      []dto.MembershipTierPriceDTO{},
	}, nil
}
