package routes

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewUserAuthRoutes),
	fx.Provide(NewAccountRoutes),
	fx.Provide(NewSwaggerRoutes),
	fx.Provide(NewRoutes),
)

type Route interface {
	Setup()
}

type Routes []Route

func NewRoutes(
	userRoutes UserAuthRoutes,
	accountRoutes AccountRoutes,
	swaggerRoutes SwaggerRoutes,
) Routes {
	return Routes{
		userRoutes,
		accountRoutes,
		swaggerRoutes,
	}
}

func (r Routes) Setup() {
	for _, route := range r {
		route.Setup()
	}
}
