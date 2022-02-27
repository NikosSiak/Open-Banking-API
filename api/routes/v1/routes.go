package routes

import "go.uber.org/fx"

var Module = fx.Options(
  fx.Provide(NewUserAuthRoutes),
  fx.Provide(NewRoutes),
)

type Route interface {
  Setup()
}

type Routes []Route

func NewRoutes(
  userRoutes UserAuthRoutes,
) Routes {
  return Routes{
    userRoutes,
  }
}

func (r Routes) Setup() {
  for _, route := range r {
    route.Setup()
  }
}
