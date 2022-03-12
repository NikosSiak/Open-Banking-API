package lib

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewEnv),
	fx.Provide(NewDB),
	fx.Provide(NewRedis),
	fx.Provide(NewRequestHandler),
)
