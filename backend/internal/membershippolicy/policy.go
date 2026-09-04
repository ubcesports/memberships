package membershippolicy

import "github.com/ubcesports/memberships/internal/dto"

type Policy interface {
	ProgramName() string

	Evaluate(
		profile *dto.ProfileDTO,
		current *dto.MembershipDTO,
		requested *dto.MembershipTierDTO,
	) (
		purchaseType dto.PurchaseType,
		allowed bool,
		err error,
	)
}
