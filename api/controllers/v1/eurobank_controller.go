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
