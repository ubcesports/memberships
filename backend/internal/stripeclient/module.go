package stripeclient

import "go.uber.org/fx"

var Module = fx.Module("stripe",
	fx.Provide(
		NewClient,
		func(client *Client) Gateway { return client },
	),
)
