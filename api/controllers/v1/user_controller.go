package controllers

import (
  "net/http"

  "github.com/gin-gonic/gin"
)

type UserController struct {
}

func NewUserController() UserController {
  return UserController{}
}

func (u UserController) CreateUser(ctx *gin.Context) {
  ctx.JSON(http.StatusOK, gin.H{
    "data": "user created",
  })
}
