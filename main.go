package main

import (
	"context"
	"fmt"

	"github.com/NikosSiak/Open-Banking-API/api/controllers/v1"
	"github.com/NikosSiak/Open-Banking-API/api/middlewares/v1"
	"github.com/NikosSiak/Open-Banking-API/api/routes/v1"
	"github.com/NikosSiak/Open-Banking-API/lib"
	"github.com/NikosSiak/Open-Banking-API/services"
	"github.com/NikosSiak/Open-Banking-API/services/providers"
	"go.uber.org/fx"
)

// @title    Open Banking Demo
// @version  1.0

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
func main() {
	fx.New(
		lib.Module,
		routes.Module,
		middlewares.Module,
		controllers.Module,
		services.Module,
		providers.Module,
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

			go func() {
				routes.Setup()
				handler.Gin.Run(":" + env.ServerPort)
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("Stopping Server")

			db.Close(ctx)

			return nil
		},
	})
}
