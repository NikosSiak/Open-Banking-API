package controllers

import (
  "context"
  "net/http"
  "time"

  "github.com/gin-gonic/gin"
  "github.com/NikosSiak/Open-Banking-API/lib"
  "github.com/NikosSiak/Open-Banking-API/models"
  "github.com/NikosSiak/Open-Banking-API/services"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
  db lib.Database
  redis lib.Redis
  authService services.AuthService
}

func NewUserController(db lib.Database, redis lib.Redis, authService services.AuthService) UserController {
  return UserController{
    db: db,
    redis: redis,
    authService: authService,
  }
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

  inserted, err := u.db.InsertOne(ctx.Request.Context(), &user)
  if err != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{
      "error": err.Error(),
    })

    return
  }

  td, err := u.authService.CreateTokens()
  if err != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{
      "error": err.Error(),
    })

    return
  }

  userId := inserted.InsertedID.(primitive.ObjectID).Hex()
  if err = u.storeTokenDetails(ctx.Request.Context(), userId, td); err != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{
      "error": err.Error(),
    })

    return
  }

  ctx.JSON(http.StatusOK, gin.H{
    "access_token": td.AccessToken,
    "refresh_token": td.RefreshToken,
  })
}

func (u UserController) AuthenticateUser(ctx *gin.Context) {
  type UserLoginCredentials struct {
    Email string `json:"email"`
    Password string `json:"password"`
  }

  userCreds := UserLoginCredentials{}

  if err := ctx.ShouldBindJSON(&userCreds); err != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{
      "error": err.Error(),
    })

    return
  }

  user := models.User{}

  err := u.db.FindOne(
    ctx.Request.Context(),
    &user,
    bson.M{ "email": userCreds.Email },
    bson.M{},
  )
  if err != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{
      "error": err.Error(),
    })

    return
  }

  if !user.CheckPasswordHash(userCreds.Password) {
    ctx.JSON(http.StatusUnauthorized, gin.H{
      "error": "wrong email or password",
    })

    return
  }

  td, err := u.authService.CreateTokens()
  if err != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{
      "error": err.Error(),
    })

    return
  }

  userId := user.ID.Hex()
  if err = u.storeTokenDetails(ctx.Request.Context(), userId, td); err != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{
      "error": err.Error(),
    })

    return
  }

  ctx.JSON(http.StatusOK, gin.H{
    "access_token": td.AccessToken,
    "refresh_token": td.RefreshToken,
  })
}

func (u UserController) LogoutUser(ctx *gin.Context) {
}

func (u UserController) storeTokenDetails(c context.Context, userId string, td *services.TokenDetails) error {
  at := time.Unix(td.AtExpires, 0)
  rt := time.Unix(td.RtExpires, 0)
  now := time.Now()

  if err := u.redis.SetEX(c, td.AccessUuid, userId, at.Sub(now)).Err(); err != nil {
    return err
  }

  if err := u.redis.SetEX(c, td.RefreshUuid, userId, rt.Sub(now)).Err(); err != nil {
    return err
  }

  return nil
}
