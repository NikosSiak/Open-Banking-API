package routes

import (
	"github.com/NikosSiak/Open-Banking-API/api/controllers/v1"
	"github.com/NikosSiak/Open-Banking-API/api/middlewares/v1"
	"github.com/NikosSiak/Open-Banking-API/lib"
)

type AccountRoutes struct {
	handler           lib.RequestHandler
	alphaController   controllers.AlphaController
	jwtAuthMiddleware middlewares.JwtAuthMiddleware
}

func NewAccountRoutes(
	handler lib.RequestHandler,
	alphaController controllers.AlphaController,
	jwtAuthMiddleware middlewares.JwtAuthMiddleware,
) AccountRoutes {
	return AccountRoutes{
		handler:           handler,
		alphaController:   alphaController,
		jwtAuthMiddleware: jwtAuthMiddleware,
	}
}

func (a AccountRoutes) Setup() {
	accountsGroup := a.handler.Gin.Group("v1/accounts")
	{
		accountsGroup.POST("/alpha", a.jwtAuthMiddleware.Handle, a.alphaController.AddAccount)
		accountsGroup.GET("/alpha/webhook", a.alphaController.AuthorizationCodeHook)
	}
}
