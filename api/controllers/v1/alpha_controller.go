package controllers

import (
	"net/http"

	"github.com/NikosSiak/Open-Banking-API/api/utils"
	"github.com/NikosSiak/Open-Banking-API/lib"
	"github.com/NikosSiak/Open-Banking-API/models"
	"github.com/NikosSiak/Open-Banking-API/services/providers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AlphaController struct {
	provider providers.Alpha
	db       lib.Database
}

func NewAlphaController(provider providers.Alpha, db lib.Database) AlphaController {
	return AlphaController{
		provider: provider,
		db:       db,
	}
}

// Add AlphaBank Account
// @Summary      Authorize the use of the user's AlphaBank account
// @Description  Use the URI to open AlphaBank's login page
// @Tags         Accounts
// @Router       /v1/accounts/alpha [post]
// @Security     BearerAuth
// @Success      200  {object}  responses.AddAlphaBankAccountResponse
// @Failure      401  {object}  responses.UnauthorizedError
// @Failure      500  {object}  utils.HTTPError
func (a AlphaController) AddAccount(ctx *gin.Context) {
	user, _ := ctx.Get("user")

	loginUri, err := a.provider.LoginUri(user.(models.User).ID.Hex())
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"uri": loginUri,
	})
}

func (a AlphaController) AuthorizationCodeHook(ctx *gin.Context) {
	code := ctx.Query("code")
	userId, err := primitive.ObjectIDFromHex(ctx.Query("state"))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"message": err,
		})
		return
	}

	user := models.User{}

	err = a.db.FindOne(ctx.Request.Context(), &user, bson.M{"_id": userId}, bson.M{})
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"message": err,
		})
		return
	}

	accessToken, err := a.provider.GetUserAccessToken(code)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"message": err,
		})
		return
	}

	user.AddAccount(a.provider.Name(), &models.Account{AccessToken: accessToken})

	err = a.db.UpdateByID(ctx.Request.Context(), user.ID, &user)

	ctx.JSON(http.StatusOK, gin.H{
		"message": err,
	})
}
