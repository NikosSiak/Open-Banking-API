package controllers

import (
	"net/http"

	"github.com/NikosSiak/Open-Banking-API/api/utils"
	"github.com/NikosSiak/Open-Banking-API/lib"
	"github.com/NikosSiak/Open-Banking-API/models"
	"github.com/NikosSiak/Open-Banking-API/services/providers"
	"github.com/gin-gonic/gin"
)

type AlphaController struct {
	provider providers.Alpha
	db       lib.Database
}

func NewAlphaController(env lib.Env, db lib.Database) AlphaController {
	return AlphaController{
		provider: providers.NewAlphaProvider(env.AppUrl, *env.Providers[providers.AlphaName]),
		db:       db,
	}
}

// Add AlphaBank Account
// @Summary      Authorize the use of the user's AlphaBank account
// @Description  Use the URI to open AlphaBank's login page
// @Tags         Accounts
// @Router       /v1/accounts/alpha [post]
// @Security     BearerAuth
// @Success      200  {object}  responses.AddBankAccountResponse
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
	userId := ctx.Query("state")

	err := providers.AddAccountToUser(ctx.Request.Context(), a.db, a.provider, code, userId)
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
