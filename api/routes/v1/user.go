package routes

import (
  "github.com/NikosSiak/Open-Banking-API/api/controllers/v1"
  "github.com/NikosSiak/Open-Banking-API/lib"
)

type UserRoutes struct {
  handler lib.RequestHandler
  controller controllers.UserController
}

func NewUserRoutes(
  handler lib.RequestHandler,
  userController controllers.UserController,
) UserRoutes {
  return UserRoutes{
    handler: handler,
    controller: userController,
  }
}

func (u UserRoutes) Setup() {
  api := u.handler.Gin.Group("/v1")

  api.POST("/user", u.controller.CreateUser)
}
