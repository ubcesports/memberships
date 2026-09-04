package membershippolicy

import (
	"context"
	"fmt"
	"math"
	"slices"

	"github.com/ubcesports/memberships/internal/dto"
	"github.com/ubcesports/memberships/internal/repository"
)

type ProfileReader interface {
	GetProfileByUserID(
		ctx context.Context,
		userID string,
	) (*dto.ProfileDTO, error)
}

type EligibilityService struct {
	membershipRepo *repository.MembershipRepository
	profileReader  ProfileReader
	policies       *Registry
}

func NewEligibilityService(
	membershipRepo *repository.MembershipRepository,
	profileReader ProfileReader,
	policies *Registry,
) *EligibilityService {
	return &EligibilityService{membershipRepo: membershipRepo, profileReader: profileReader, policies: policies}
}

/*
	Public functions
*/

func (s *EligibilityService) GetEligibleTiers(ctx context.Context, userId string) ([]dto.EligibleMembershipTierDTO, error) {
	// Get user info
	user, err := s.profileReader.GetProfileByUserID(ctx, userId)
	if err != nil {
		return nil, err
	}

	tiers, err := s.membershipRepo.GetActiveTiersWithPrices(ctx)
	if err != nil {
		return nil, err
	}

	memberships, err := s.membershipRepo.GetCurrentMembershipsWithTransactions(ctx, userId)
	if err != nil {
		return nil, err
	}

	result := make([]dto.EligibleMembershipTierDTO, 0, len(tiers))
	for i := range tiers {
		tier := &tiers[i]

		// 1. Check if user group is eligible
		if !isUserGroupEligible(user.Groups, tier.RequiredGroup) {
			continue
		}

		// 2. Get user's personalized price
		price, err := getTierPriceForUser(user.IsStudent, *tier)
		if err != nil {
			return nil, err
		}
		if price == nil {
			continue
		}

		// 3. Check if user alr has a membership for the current program
		current, err := findCurrentMembershipForProgram(memberships, tier.ProgramId)
		if err != nil {
			return nil, err
		}

		// 4. Check policies for current program based on user's membership
		policy, err := s.policies.Get(tier.ProgramName)
		if err != nil {
			return nil, err
		}

		purchaseType, allowed, err := policy.Evaluate(user, current, tier)
		if err != nil {
			return nil, err
		}
		if !allowed {
			continue
		}

		finalPrice := price.Price

		// 5. Adjust price for upgrades
		if purchaseType == dto.PurchaseUpgrade {
			if current == nil {
				return nil, fmt.Errorf("upgrade not allowed without a membership in this program")
			}

			targetCents := int64(math.Round(price.Price * 100))
			amountDueCents := targetCents - current.Transaction.AmountPaidCents
			if amountDueCents < 0 {
				amountDueCents = 0
			}

			finalPrice = float64(amountDueCents) / 100
		}

		result = append(result, dto.EligibleMembershipTierDTO{
			ID:           tier.ID,
			Title:        tier.Title,
			Description:  tier.Description,
			Slug:         tier.Slug,
			PurchaseType: purchaseType,
			ProductId:    tier.ProductId,
			Benefits:     tier.Benefits,
			ProgramId:    tier.ProgramId,
			ProgramName:  tier.ProgramName,
			Price: dto.MembershipTierPriceDTO{
				Price:             finalPrice,
				PriceId:           price.PriceId,
				IsStudentRequired: nil,
			},
		})
	}

	return result, nil
}

/*
	Helper functions
*/

func isUserGroupEligible(
	userGroups []dto.GroupType,
	requiredGroup dto.GroupType,
) bool {
	if slices.Contains(userGroups, requiredGroup) {
		return true
	}

	if requiredGroup == dto.GroupExecutive {
		return slices.Contains(userGroups, dto.GroupDirector) || slices.Contains(userGroups, dto.GroupBoard)
	}

	return false
}

func getTierPriceForUser(
	studentStatus bool,
	reqTier dto.MembershipTierDTO,
) (*dto.MembershipTierPriceDTO, error) {
	var match *dto.MembershipTierPriceDTO

	for i := range reqTier.Prices {
		currPrice := &reqTier.Prices[i]

		eligible := currPrice.IsStudentRequired == nil || *currPrice.IsStudentRequired == studentStatus
		if !eligible {
			continue
		}

		match = currPrice
		break
	}

	return match, nil
}

func findCurrentMembershipForProgram(
	memberships []dto.MembershipDTO,
	programId string,
) (*dto.MembershipDTO, error) {
	var match *dto.MembershipDTO

	for i := range memberships {
		membership := &memberships[i]

		if membership.ProgramId != programId {
			continue
		}

		if match != nil {
			return nil, fmt.Errorf(
				"multiple active memberships found for program %q",
				membership.ProgramName,
			)
		}

		match = membership
	}

	return match, nil
}
