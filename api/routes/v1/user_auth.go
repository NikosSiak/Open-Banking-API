package routes

import (
  "github.com/NikosSiak/Open-Banking-API/api/controllers/v1"
  "github.com/NikosSiak/Open-Banking-API/api/middlewares/v1"
  "github.com/NikosSiak/Open-Banking-API/lib"
)

type UserAuthRoutes struct {
  handler lib.RequestHandler
  controller controllers.UserController
  jwtAuthMiddleware middlewares.JwtAuthMiddleware
}

func NewUserAuthRoutes(
  handler lib.RequestHandler,
  userController controllers.UserController,
  jwtAuthMiddleware middlewares.JwtAuthMiddleware,
) UserAuthRoutes {
  return UserAuthRoutes{
    handler: handler,
    controller: userController,
    jwtAuthMiddleware: jwtAuthMiddleware,
  }
}

func (u UserAuthRoutes) Setup() {
  u.handler.Gin.POST("/register", u.controller.CreateUser)
  u.handler.Gin.POST("/login", u.controller.AuthenticateUser)
  u.handler.Gin.POST("/logout", u.jwtAuthMiddleware.Handle, u.controller.LogoutUser)
}
