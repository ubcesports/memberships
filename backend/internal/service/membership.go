package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"slices"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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
			Price:             fmt.Sprintf("%.2f", float64(price.UnitAmount)/100), // Turn unit amount which is in cents, into readable format with 2 numbers after the decimal
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

func (s *MembershipService) GetEligibleTiersWithPrices(ctx context.Context, userId string) (*[]dto.EligibleMembershipTierDTO, error) {
	tiers, err := s.membershipRepo.GetEligibleTiersWithPrices(ctx, userId)
	if err != nil {
		return nil, err
	}

	returnTiers := make([]dto.EligibleMembershipTierDTO, 0, len(tiers))
	tierIndexById := make(map[string]int)

	// Get user info
	var pgUserId pgtype.UUID
	if err := pgUserId.Scan(userId); err != nil {
		return nil, err
	}
	user, err := s.profileService.GetProfileByUserID(ctx, pgUserId)
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

		if index, exists := tierIndexById[tierId]; exists {
			returnTiers[index].Prices = append(returnTiers[index].Prices, priceDto)
			return
		}

		tierIndexById[tierId] = len(returnTiers)
		returnTiers = append(returnTiers, dto.EligibleMembershipTierDTO{
			ID:           tier.ID.String(),
			Title:        tier.Title,
			Description:  tier.Description.String,
			Slug:         tier.Slug.String,
			PurchaseType: purchaseType,
			Prices:       []dto.MembershipTierPriceDTO{priceDto},
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

		switch purchaseType {
		case dto.PurchaseNew:
			// Get price from stripe price id
			price, err := s.stripeClient.GetPrice(ctx, tier.StripePriceID.String)
			if err != nil {
				return nil, err
			}

			// Set up price dto
			priceDto := dto.MembershipTierPriceDTO{
				Price:             fmt.Sprintf("%.2f", float64(price.UnitAmount)/100), // Turn unit amount which is in cents, into readable format with 2 numbers after the decimal
				IsStudentRequired: nil,                                                // Leave nil as this is not really required in this context.
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
				Price:             fmt.Sprintf("%.2f", float64(priceToPay)/100),
				IsStudentRequired: nil, // Leave nil as this is not really required in this context.
			}

			addPriceToTier(tier, purchaseType, priceDto)
		}
	}

	return &returnTiers, nil
}

func (s *MembershipService) getTierByTierId(ctx context.Context, tierId string) (*dto.MembershipTierDTO, error) {
	tier, err := s.membershipRepo.GetTierByTierId(ctx, tierId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &dto.MembershipTierDTO{
		ID:          tier.ID.String(),
		Title:       tier.Title,
		Description: tier.Description.String,
		Slug:        tier.Slug.String,
		Prices:      []dto.MembershipTierPriceDTO{},
	}, nil
}
