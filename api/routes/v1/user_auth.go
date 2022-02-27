package routes

import (
  "github.com/NikosSiak/Open-Banking-API/api/controllers/v1"
  "github.com/NikosSiak/Open-Banking-API/api/middlewares/v1"
  "github.com/NikosSiak/Open-Banking-API/lib"
)

type UserAuthRoutes struct {
  handler lib.RequestHandler
  controller controllers.UserController
}

func NewUserAuthRoutes(
  handler lib.RequestHandler,
  userController controllers.UserController,
) UserAuthRoutes {
  return UserAuthRoutes{
    handler: handler,
    controller: userController,
  }
}

func (u UserAuthRoutes) Setup() {
  u.handler.Gin.POST("/register", u.controller.CreateUser)
  u.handler.Gin.POST("/login", u.controller.AuthenticateUser)
  u.handler.Gin.POST("/logout", u.controller.LogoutUser)
}
