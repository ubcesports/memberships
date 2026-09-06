package membershippolicy

import "github.com/ubcesports/memberships/internal/dto"

func NewSSBMPolicy() Policy {
	return NewTransitionPolicy(
		"ssbm",
		map[Transition]dto.PurchaseType{
			{From: "", To: "ssbm_semesterly"}: dto.PurchaseNew,
			{From: "", To: "ssbm_yearly"}:     dto.PurchaseNew,

			{
				From: "ssbm_semesterly",
				To:   "ssbm_yearly",
			}: dto.PurchaseUpgrade,
		},
	)
}
