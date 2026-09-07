package membershippolicy

import (
	"slices"

	"github.com/ubcesports/memberships/internal/dto"
)

type GeneralPolicy struct {
	standard *TransitionPolicy
}

func NewGeneralPolicy() Policy {
	return &GeneralPolicy{
		standard: NewTransitionPolicy(
			"general",
			map[Transition]dto.PurchaseType{
				{From: "", To: "day"}:    dto.PurchaseNew,
				{From: "", To: "basic"}:  dto.PurchaseNew,
				{From: "", To: "lounge"}: dto.PurchaseNew,

				{From: "day", To: "basic"}:  dto.PurchaseNew,
				{From: "day", To: "lounge"}: dto.PurchaseNew,

				{
					From: "basic",
					To:   "lounge",
				}: dto.PurchaseUpgrade,
			},
		),
	}
}

func (p *GeneralPolicy) ProgramName() string {
	return "general"
}

func (p *GeneralPolicy) Evaluate(
	profile *dto.ProfileDTO,
	current *dto.MembershipDTO,
	requested *dto.MembershipTierDTO,
) (dto.PurchaseType, bool, error) {
	// Exec groups have the highest priority inside general
	if slices.Contains(profile.Groups, dto.GroupExecutive) ||
		slices.Contains(profile.Groups, dto.GroupDirector) ||
		slices.Contains(profile.Groups, dto.GroupBoard) {
		return allowRestrictedGeneralTier(
			current,
			requested,
			"executive",
		)
	}

	// Comp players is the next priority
	if slices.Contains(profile.Groups, dto.GroupCompetitiveTeam) {
		return allowRestrictedGeneralTier(
			current,
			requested,
			"competitive_team",
		)
	}

	// Regular members follow day/basic/lounge transitions
	return p.standard.Evaluate(
		profile,
		current,
		requested,
	)
}

func allowRestrictedGeneralTier(
	current *dto.MembershipDTO,
	requested *dto.MembershipTierDTO,
	requiredSlug string,
) (dto.PurchaseType, bool, error) {
	// The user already has an active general membership
	if current != nil {
		return "", false, nil
	}

	if requested.Slug != requiredSlug {
		return "", false, nil
	}

	return dto.PurchaseNew, true, nil
}
