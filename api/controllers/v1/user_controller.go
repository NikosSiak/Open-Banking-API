package controllers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/NikosSiak/Open-Banking-API/lib"
	"github.com/NikosSiak/Open-Banking-API/models"
	"github.com/NikosSiak/Open-Banking-API/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	db          lib.Database
	redis       lib.Redis
	authService services.AuthService
}

func NewUserController(db lib.Database, redis lib.Redis, authService services.AuthService) UserController {
	return UserController{
		db:          db,
		redis:       redis,
		authService: authService,
	}
}

// Register User
// @Summary  Register a new User
// @Tags     User
// @Router   /register [post]
// @Param    email     body      string  true  "User email"
// @Param    password  body      string  true  "User password"
// @Success  200       {object}  responses.TokenResponse
// @Failure  500       {object}  utils.HTTPError
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
		"access_token":  td.AccessToken,
		"refresh_token": td.RefreshToken,
	})
}

// Login
// @Summary  Get access and refresh tokens for user
// @Tags     User
// @Router   /login [post]
// @Param    data  body      models.UserLoginCredentials  true  "User credentials"
// @Success  200   {object}  responses.TokenResponse
// @Failure  401   {object}  utils.HTTPError
// @Failure  500   {object}  utils.HTTPError
func (u UserController) AuthenticateUser(ctx *gin.Context) {
	userCreds := models.UserLoginCredentials{}

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
		bson.M{"email": userCreds.Email},
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
		"access_token":  td.AccessToken,
		"refresh_token": td.RefreshToken,
	})
}

// Logout
// @Summary   Invalidate User tokens
// @Tags      User
// @Router    /logout [post]
// @Security  BearerAuth
// @Success   200  {object}  responses.SuccessResponse
// @Failure   401  {object}  responses.UnauthorizedError
// @Failure   500  {object}  utils.HTTPError
func (u UserController) LogoutUser(ctx *gin.Context) {
	authHeader := ctx.Request.Header.Get("Authorization")
	t := strings.Split(authHeader, " ")

	if len(t) != 2 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "missing bearer token",
		})

		ctx.Abort()
		return
	}

	accessUuid, _ := u.authService.GetAccessUuid(t[1])
	if err := u.deleteToken(ctx.Request.Context(), accessUuid); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
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

func (u UserController) deleteToken(c context.Context, tokenUuid string) error {
	_, err := u.redis.Del(c, tokenUuid).Result()
	if err != nil {
		return err
	}

	return nil
}
