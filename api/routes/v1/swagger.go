package routes

import (
	"net/url"

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
	if !s.env.IsProduction() {
		url, err := url.Parse(s.env.AppUrl)
		if err != nil {
			return
		}

		docs.SwaggerInfo.Host = url.Host
		docs.SwaggerInfo.Schemes = []string{url.Scheme}
		s.handler.Gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
}
