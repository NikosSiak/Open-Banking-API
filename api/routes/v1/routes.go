package routes

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewUserAuthRoutes),
	fx.Provide(NewAccountRoutes),
	fx.Provide(NewRoutes),
)

type Route interface {
	Setup()
}

type Routes []Route

func NewRoutes(
	userRoutes UserAuthRoutes,
	accountRoutes AccountRoutes,
) Routes {
	return Routes{
		userRoutes,
		accountRoutes,
	}
}

func (r Routes) Setup() {
	for _, route := range r {
		route.Setup()
	}
}
