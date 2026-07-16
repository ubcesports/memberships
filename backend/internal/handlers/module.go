package handlers

import "go.uber.org/fx"

var Module = fx.Module("repository",
	fx.Provide(
		NewHealthHandler,
		NewProfileHandler,
		NewAdminHandler,
	),
)
