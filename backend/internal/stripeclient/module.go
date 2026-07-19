package stripeclient

import "go.uber.org/fx"

var Module = fx.Module("stripe",
	fx.Provide(NewClient),
)
