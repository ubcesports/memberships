package mailer

import (
	"go.uber.org/fx"
)

var Module = fx.Module("mailer",
	fx.Invoke(invokeMailer),
)

func invokeMailer() error {
	return Init()
}
