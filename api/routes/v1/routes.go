package routes

import "go.uber.org/fx"

var Module = fx.Options(
  fx.Provide(NewUserRoutes),
  fx.Provide(NewRoutes),
)

type Route interface {
  Setup()
}

type Routes []Route

func NewRoutes(
  userRoutes UserRoutes,
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
