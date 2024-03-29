package controllers

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/NikosSiak/Open-Banking-API/api/utils"
	"github.com/NikosSiak/Open-Banking-API/lib"
	"github.com/NikosSiak/Open-Banking-API/lib/sms"
	"github.com/NikosSiak/Open-Banking-API/models"
	"github.com/NikosSiak/Open-Banking-API/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	db          lib.Database
	redis       lib.Redis
	authService services.AuthService
	sms         sms.SMSProvider
}

func NewUserController(
	db lib.Database,
	redis lib.Redis,
	authService services.AuthService,
	sms sms.SMSProvider,
) UserController {
	return UserController{
		db:          db,
		redis:       redis,
		authService: authService,
		sms:         sms,
	}
}

// Register User
// @Summary  Register a new User
// @Tags     User
// @Router   /register [post]
// @Param    data  body      models.User  true  "User info"
// @Success  200   {object}  responses.LoginResponse
// @Failure  500   {object}  utils.HTTPError
func (u UserController) CreateUser(ctx *gin.Context) {
	user := models.User{}

	if err := ctx.ShouldBindJSON(&user); err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	if err := validator.New().Struct(user); err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	user.HasTwoFa = true

	if err := user.HashPassword(); err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	user.Accounts = map[string]*models.AccountTokens{}

	inserted, err := u.db.InsertOne(ctx.Request.Context(), &user)
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	userId := inserted.InsertedID.(primitive.ObjectID).Hex()

	if user.HasTwoFa {
		verificationId, err := u.sms.SendVerificationCode(user.PhoneNumber)
		if err != nil {
			utils.NewError(ctx, http.StatusInternalServerError, err)
			return
		}

		if err = u.redis.Set(ctx.Request.Context(), verificationId, userId, 0).Err(); err != nil {
			utils.NewError(ctx, http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"verification_id": verificationId,
		})

		return
	}

	td, err := u.authService.CreateTokens()
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	if err = u.storeTokenDetails(ctx.Request.Context(), userId, td); err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  td.AccessToken,
		"refresh_token": td.RefreshToken,
	})
}

// Login
// @Summary      Get access and refresh tokens for user
// @Tags         User
// @Description  If the user has enabled TwoFa the result will have a verification ID for the verify route, else the access and refresh tokens are returned
// @Router       /login [post]
// @Param        data  body      models.UserLoginCredentials  true  "User credentials"
// @Success      200   {object}  responses.LoginResponse
// @Failure      401   {object}  utils.HTTPError
// @Failure      500   {object}  utils.HTTPError
func (u UserController) AuthenticateUser(ctx *gin.Context) {
	userCreds := models.UserLoginCredentials{}

	if err := ctx.ShouldBindJSON(&userCreds); err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	if err := validator.New().Struct(userCreds); err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
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
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	if !user.CheckPasswordHash(userCreds.Password) {
		utils.NewError(ctx, http.StatusUnauthorized, errors.New("wrong email or password"))
		return
	}

	if user.HasTwoFa {
		verificationId, err := u.sms.SendVerificationCode(user.PhoneNumber)
		if err != nil {
			utils.NewError(ctx, http.StatusInternalServerError, err)
			return
		}

		if err = u.redis.Set(ctx.Request.Context(), verificationId, user.ID.Hex(), 0).Err(); err != nil {
			utils.NewError(ctx, http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"verification_id": verificationId,
		})

		return
	}

	td, err := u.authService.CreateTokens()
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	userId := user.ID.Hex()
	if err = u.storeTokenDetails(ctx.Request.Context(), userId, td); err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  td.AccessToken,
		"refresh_token": td.RefreshToken,
	})
}

// Two Factor Authentication
// @Summary  Verify a user login with a twofa code
// @Tags     User
// @Router   /verify [post]
// @Param    verification_id  query     string  true  "Verification ID provided by login"
// @Param    code             query     string  true  "TwoFactor authentication code"
// @Success  200              {object}  responses.TokenResponse
// @Failure  401              {object}  utils.HTTPError
// @Failure  500              {object}  utils.HTTPError
func (u UserController) ValidateCode(ctx *gin.Context) {
	verificationId := ctx.Query("verification_id")
	code := ctx.Query("code")

	valid, err := u.sms.VerifyCode(verificationId, code)
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	if !valid {
		utils.NewError(ctx, http.StatusUnauthorized, errors.New("invalid code"))
		return
	}

	td, err := u.authService.CreateTokens()
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	userId, err := u.redis.Get(ctx.Request.Context(), verificationId).Result()
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	if err = u.storeTokenDetails(ctx.Request.Context(), userId, td); err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	u.redis.Del(ctx.Request.Context(), verificationId)

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

	accessUuid, _ := u.authService.GetAccessUuid(t[1])
	if err := u.deleteToken(ctx.Request.Context(), accessUuid); err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

// Refresh Tokens
// @Summary  Get a new access and refresh tokens
// @Tags     User
// @Router   /refresh [post]
// @Param    refresh_token  query     string  true  "Refresh token provided by login"
// @Success  200            {object}  responses.TokenResponse
// @Failure  500            {object}  utils.HTTPError
func (u UserController) RefreshTokens(ctx *gin.Context) {
	refreshToken := ctx.Query("refresh_token")

	userId, err := u.redis.Get(ctx, refreshToken).Result()
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	td, err := u.authService.CreateTokens()
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	u.deleteToken(ctx.Request.Context(), refreshToken)

	if err = u.storeTokenDetails(ctx.Request.Context(), userId, td); err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  td.AccessToken,
		"refresh_token": td.RefreshToken,
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
