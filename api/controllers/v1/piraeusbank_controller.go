package controllers

import (
	"net/http"

	"github.com/NikosSiak/Open-Banking-API/api/utils"
	"github.com/NikosSiak/Open-Banking-API/lib"
	"github.com/NikosSiak/Open-Banking-API/models"
	"github.com/NikosSiak/Open-Banking-API/services/providers"
	"github.com/gin-gonic/gin"
)

type PiraeusbankController struct {
	provider providers.Piraeusbank
	db       lib.Database
}

func NewPiraeusbankController(env lib.Env, db lib.Database) PiraeusbankController {
	return PiraeusbankController{
		provider: providers.NewPiraeusbankProvider(env.AppUrl, *env.Providers[providers.PiraeusbankName]),
		db:       db,
	}
}

// Add Piraeusbank Account
// @Summary      Authorize the use of the user's Piraeusbank account
// @Description  Use the URI to open Piraeusbank's login page<br>You can find login credentials <a href="https://rapidlink.piraeusbank.gr/node/2059">here</a>
// @Tags         Accounts
// @Router       /v1/accounts/piraeusbank [post]
// @Security     BearerAuth
// @Success      200  {object}  responses.AddBankAccountResponse
// @Failure      401  {object}  responses.UnauthorizedError
// @Failure      500  {object}  utils.HTTPError
func (p PiraeusbankController) AddAccount(ctx *gin.Context) {
	user, _ := ctx.Get("user")

	ctx.JSON(http.StatusOK, gin.H{
		"uri": p.provider.LoginUri(user.(models.User).ID.Hex()),
	})
}

func (p PiraeusbankController) AuthorizationCodeHook(ctx *gin.Context) {
	code := ctx.Query("code")
	userId := ctx.Query("state")

	err := providers.AddAccountToUser(ctx.Request.Context(), p.db, p.provider, code, userId)
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
