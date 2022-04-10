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
	userId, err := primitive.ObjectIDFromHex(ctx.Query("state"))
	if err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	user := models.User{}

	err = e.db.FindOne(ctx.Request.Context(), &user, bson.M{"_id": userId}, bson.M{})
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	account, err := e.provider.GetUserTokens(code)
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	user.AddAccount(e.provider.Name(), &account)

	err = e.db.UpdateByID(ctx.Request.Context(), user.ID, &user)
	if err != nil {
		utils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
