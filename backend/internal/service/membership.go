package service

import (
	"context"
	"fmt"

	"github.com/ubcesports/memberships/internal/dto"
	"github.com/ubcesports/memberships/internal/repository"
	"github.com/ubcesports/memberships/internal/stripeclient"
)

type MembershipService struct {
	membershipRepo *repository.MembershipRepository
	stripeClient   *stripeclient.Client
}

func NewMembershipService(membershipRepo *repository.MembershipRepository, stripeClient *stripeclient.Client) *MembershipService {
	return &MembershipService{membershipRepo: membershipRepo, stripeClient: stripeClient}
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
				Prices:      []dto.MembershipTierPriceDTO{priceDto},
			})
		}
	}

	return returnTiers, nil
}
