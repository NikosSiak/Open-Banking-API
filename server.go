package main

import (
  "context"
  "fmt"

  "github.com/NikosSiak/Open-Banking-API/api/controllers/v1"
  "github.com/NikosSiak/Open-Banking-API/api/routes/v1"
  "github.com/NikosSiak/Open-Banking-API/lib"
  "go.uber.org/fx"
)

func main() {
  fx.New(
    lib.Module,
    routes.Module,
    controllers.Module,
    fx.Invoke(bootstrap),
  ).Run()
}

func bootstrap(
  lifecycle fx.Lifecycle,
  handler lib.RequestHandler,
  env lib.Env,
  db lib.Database,
  routes routes.Routes,
) {
  lifecycle.Append(fx.Hook{
    OnStart: func(context.Context) error {
      fmt.Println("Starting Server")

      go func () {
        routes.Setup()
        handler.Gin.Run(":" + env.ServerPort)
      }()

      return nil
    },
    OnStop: func(context.Context) error {
      fmt.Println("Stopping Server")

      return nil
    },
  })
}
