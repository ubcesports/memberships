package service

import "go.uber.org/fx"

var Module = fx.Module("repository",
	fx.Provide(
		NewHealthService,
		NewOnboardingService,
		NewProfileService,
		NewAdminUserService,
	),
)
