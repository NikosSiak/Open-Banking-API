package routes

import (
	"github.com/NikosSiak/Open-Banking-API/docs"
	"github.com/NikosSiak/Open-Banking-API/lib"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type SwaggerRoutes struct {
	handler lib.RequestHandler
	env     lib.Env
}

func NewSwaggerRoutes(handler lib.RequestHandler, env lib.Env) SwaggerRoutes {
	return SwaggerRoutes{handler: handler, env: env}
}

func (s SwaggerRoutes) Setup() {
	docs.SwaggerInfo.Host = s.env.AppUrl
	s.handler.Gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
