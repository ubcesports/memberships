package membershippolicy

import (
	"go.uber.org/fx"
)

var Module = fx.Module("membership_policy",
	fx.Provide(
		NewRegistry,
		NewEligibilityService,
	),
)
