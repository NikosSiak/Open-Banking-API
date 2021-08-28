package controllers

import (
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/NikosSiak/Open-Banking-API/lib"
  "github.com/NikosSiak/Open-Banking-API/models"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
  db lib.Database
}

func NewUserController(db lib.Database) UserController {
  return UserController{ db: db }
}

func (u UserController) CreateUser(ctx *gin.Context) {
  user := models.User{}

  if err := ctx.ShouldBindJSON(&user); err != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{
      "error": err.Error(),
    })

    return
  }

  if err := user.HashPassword(); err != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{
      "error": err.Error(),
    })

    return
  }

  user.Accounts = map[string]*models.Account{}

  inserted, err := u.db.InsertOne(ctx.Request.Context(), user)
  if err != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{
      "error": err.Error(),
    })

    return
  }

  ctx.JSON(http.StatusOK, gin.H{
    "user_id": inserted.InsertedID,
  })
}
