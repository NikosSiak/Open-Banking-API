package controllers

import (
	"net/http"

	"github.com/NikosSiak/Open-Banking-API/api/utils"
	"github.com/NikosSiak/Open-Banking-API/lib"
	"github.com/NikosSiak/Open-Banking-API/models"
	"github.com/NikosSiak/Open-Banking-API/services/providers"
	"github.com/gin-gonic/gin"
)

type EurobankController struct {
	provider providers.Eurobank
	db       lib.Database
}

func NewEurobankController(env lib.Env, db lib.Database) EurobankController {
	return EurobankController{
		provider: providers.NewEurobankProvider(env.AppUrl, *env.Providers[providers.EurobankName]),
		db:       db,
	}
}

// Add Eurobank Account
// @Summary      Authorize the use of the user's Eurobank account
// @Description  Use the URI to open Eurobank's login page<br>You can find login credentials <a href="https://developer.eurobank.gr/eurobank/apis/support">here</a> under "Which users can be used to log-in to sandbox environment?"
// @Tags         Accounts
// @Router       /v1/accounts/eurobank [post]
// @Security     BearerAuth
// @Success      200  {object}  responses.AddBankAccountResponse
// @Failure      401  {object}  responses.UnauthorizedError
// @Failure      500  {object}  utils.HTTPError
func (e EurobankController) AddAccount(ctx *gin.Context) {
	user, _ := ctx.Get("user")

	ctx.JSON(http.StatusOK, gin.H{
		"uri": e.provider.LoginUri(user.(models.User).ID.Hex()),
	})
}

func (e EurobankController) AuthorizationCodeHook(ctx *gin.Context) {
	code := ctx.Query("code")
	userId := ctx.Query("state")

	err := providers.AddAccountToUser(ctx.Request.Context(), e.db, e.provider, code, userId)
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
