package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stripe/stripe-go/v86"
	"github.com/ubcesports/memberships/internal/database/db"
	"github.com/ubcesports/memberships/internal/dto"
	"github.com/ubcesports/memberships/internal/membershippolicy"
	"github.com/ubcesports/memberships/internal/repository"
	"github.com/ubcesports/memberships/internal/stripeclient"
)

type MembershipService struct {
	membershipRepo     *repository.MembershipRepository
	stripeClient       *stripeclient.Client
	profileService     *ProfileService
	eligibilityService *membershippolicy.EligibilityService
}

func NewMembershipService(membershipRepo *repository.MembershipRepository, stripeClient *stripeclient.Client, profileService *ProfileService, eligibilityService *membershippolicy.EligibilityService) *MembershipService {
	return &MembershipService{membershipRepo: membershipRepo, stripeClient: stripeClient, profileService: profileService, eligibilityService: eligibilityService}
}

/*
	Consts
*/

var (
	ErrMembershipAlreadyExists    = errors.New("An active un-upgradeable membership already exists!")
	ErrTierNotEligible            = errors.New("Requested membership tier not eligible for current user. Please try another value.")
	ErrTierNotFound               = errors.New("Tier with given tier id not found.")
	ErrMembershipPurchaseClosed   = errors.New("Membership purchases are closed until the next membership period.")
	ErrPendingCheckoutAlreadyPaid = errors.New("A previous checkout payment is still being processed. Please wait a moment and refresh or contact an admin.")
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
		if !tier.PriceInCents.Valid {
			return nil, fmt.Errorf("membership tier price %s has no database price", tier.StripePriceID.String)
		}

		// Set up price dto
		var isStudentRequired *bool
		if tier.IsStudentRequired.Valid {
			isStudentRequired = &tier.IsStudentRequired.Bool
		} else {
			isStudentRequired = nil
		}

		priceDto := dto.MembershipTierPriceDTO{
			Price:             float64(tier.PriceInCents.Int64),
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
				Benefits:    tier.Benefits,
				Slug:        tier.Slug.String,
				ProductId:   tier.StripeProductID.String,
				Prices:      []dto.MembershipTierPriceDTO{priceDto},
				ProgramId:   tier.ProgramID.String(),
				ProgramName: tier.ProgramName,
			})
		}
	}

	return returnTiers, nil
}

func (s *MembershipService) GetCurrentMembershipsWithTransactions(ctx context.Context, userId string) ([]dto.MembershipDTO, error) {
	return s.membershipRepo.GetCurrentMembershipsWithTransactions(ctx, userId)
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
		var cancelledAt *time.Time
		if membership.CancelledAt.Valid {
			cancelledAt = &membership.CancelledAt.Time
		}

		membershipDto := dto.MembershipDTO{
			ID:          membership.ID.String(),
			TierId:      membership.TierID.String(),
			TierTitle:   membership.TierTitle,
			Slug:        membership.Slug.String,
			StartedAt:   membership.StartedAt.Time,
			ExpiresAt:   membership.ExpiresAt.Time,
			CancelledAt: cancelledAt,
			Transaction: dto.TransactionDTO{
				ID:              membership.TransactionID.String(),
				AmountPaid:      fmt.Sprintf("%.2f", float64(membership.AmountPaidCents.Int64)/100),
				Status:          dto.TransactionStatusType(membership.Status),
				GroupAtPurchase: dto.GroupType(membership.GroupAtPurchase.GroupType),
			},
			ProgramId:   membership.ProgramID.String(),
			ProgramName: membership.ProgramName,
		}
		returnMemberships = append(returnMemberships, membershipDto)
	}

	return &returnMemberships, nil
}

func (s *MembershipService) GetEligibleTiersWithPrices(ctx context.Context, userId string) ([]dto.EligibleMembershipTierDTO, error) {
	return s.eligibilityService.GetEligibleTiers(ctx, userId)
}

func (s *MembershipService) CreateCheckoutSession(ctx context.Context, userId string, req dto.CheckoutSessionRequest) (*dto.CheckoutSessionResponse, error) {
	// 1. Let the centralized eligibility service determine whether the
	// requested tier can be purchased, its purchase type, and final price.
	eligibleTiers, err := s.eligibilityService.GetEligibleTiers(
		ctx,
		userId,
	)
	if err != nil {
		return nil, err
	}

	// 2. Check if user is eligible for this tier
	var selectedTier *dto.EligibleMembershipTierDTO
	for _, t := range eligibleTiers {
		if req.TierId == t.ID {
			selectedTier = &t
			break
		}
	}
	if selectedTier == nil {
		return nil, ErrTierNotEligible
	}

	// 3. Check whether purchases are currently closed.
	isClosed, err := membershippolicy.IsPurchaseClosed(time.Now(), selectedTier.ExpirationType)
	if err != nil {
		return nil, err
	}
	if isClosed {
		return nil, ErrMembershipPurchaseClosed
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

			if session.Status == stripe.CheckoutSessionStatusComplete {
				return ErrPendingCheckoutAlreadyPaid
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

		// 2. Ensure membership purchase is not closed
		isClosed, err := membershippolicy.IsPurchaseClosed(occurredAt, dto.MembershipExpirationType(transaction.ExpirationType))
		if err != nil {
			return err
		}
		if isClosed {
			return mr.ExpirePendingTransactionById(ctx, transaction.ID.String())
		}

		// 3. Make the handler idempotent: Stripe can retry or duplicate webhook delivery.
		if transaction.Status == db.TransactionStatusTypeCompleted {
			return nil
		}
		if transaction.Status != db.TransactionStatusTypePending {
			return nil
		}

		// 4. Confirm Stripe says this checkout is paid before fulfilling the membership.
		if session.PaymentStatus != stripe.CheckoutSessionPaymentStatusPaid {
			return fmt.Errorf("checkout session %s is not paid", session.ID)
		}

		// 5. Cross-check Stripe metadata against the locked DB transaction.
		if session.Metadata["transaction_id"] != transaction.ID.String() {
			return fmt.Errorf("checkout session %s transaction metadata mismatch", session.ID)
		}
		if session.Metadata["user_id"] != transaction.UserID.String() {
			return fmt.Errorf("checkout session %s user metadata mismatch", session.ID)
		}

		// 6. Cancel any old active membership before creating the new fulfilled membership.
		if err := mr.CancelActiveMembershipsByUserIdAndProgramId(ctx, transaction.UserID.String(), transaction.ProgramID.String(), occurredAt); err != nil {
			return err
		}

		// 7. Create the fulfilled membership.
		expiresAt, err := membershippolicy.MembershipExpiresAt(occurredAt, dto.MembershipExpirationType(transaction.ExpirationType))
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

		// 8. Record payment details and mark the transaction completed.
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
		Benefits:    tier.Benefits,
		Slug:        tier.Slug.String,
		ProductId:   tier.StripeProductID.String,
		Prices:      []dto.MembershipTierPriceDTO{},
		ProgramId:   tier.ProgramID.String(),
		ProgramName: tier.ProgramName,
	}, nil
}
