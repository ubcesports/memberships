package membershippolicy

import (
	"fmt"

	"github.com/ubcesports/memberships/internal/dto"
)

type Transition struct {
	From string
	To   string
}

type TransitionPolicy struct {
	programName string
	transitions map[Transition]dto.PurchaseType
}

func NewTransitionPolicy(
	programName string,
	transitions map[Transition]dto.PurchaseType,
) *TransitionPolicy {
	return &TransitionPolicy{
		programName: programName,
		transitions: transitions,
	}
}

func (p *TransitionPolicy) ProgramName() string {
	return p.programName
}

func (p *TransitionPolicy) Evaluate(
	_ *dto.ProfileDTO,
	current *dto.MembershipDTO,
	requested *dto.MembershipTierDTO,
) (dto.PurchaseType, bool, error) {
	if requested.ProgramName != p.programName {
		return "", false, fmt.Errorf(
			"policy %q cannot evaluate program %q",
			p.programName,
			requested.ProgramName,
		)
	}

	from := ""
	if current != nil {
		from = current.Slug
	}

	purchaseType, allowed := p.transitions[Transition{
		From: from,
		To:   requested.Slug,
	}]

	return purchaseType, allowed, nil
}
