package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/NikosSiak/Open-Banking-API/lib"
	"github.com/NikosSiak/Open-Banking-API/models"
	"github.com/NikosSiak/Open-Banking-API/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JwtAuthMiddleware struct {
	db          lib.Database
	redis       lib.Redis
	authService services.AuthService
}

func NewJwtAuthMiddleware(
	db lib.Database,
	redis lib.Redis,
	authService services.AuthService,
) JwtAuthMiddleware {
	return JwtAuthMiddleware{db: db, redis: redis, authService: authService}
}

func (j JwtAuthMiddleware) Handle(ctx *gin.Context) {
	authHeader := ctx.Request.Header.Get("Authorization")
	t := strings.Split(authHeader, " ")

	if len(t) != 2 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "missing bearer token",
		})

		ctx.Abort()
		return
	}

	accessUuid, err := j.authService.GetAccessUuid(t[1])
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})

		ctx.Abort()
		return
	}

	user, err := j.getUser(ctx.Request.Context(), accessUuid)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})

		ctx.Abort()
		return
	}

	ctx.Set("user", *user)
	ctx.Next()
}

func (j JwtAuthMiddleware) getUser(ctx context.Context, accessUuid string) (*models.User, error) {
	_userId, err := j.redis.Get(ctx, accessUuid).Result()
	if err != nil {
		return nil, err
	}

	userId, err := primitive.ObjectIDFromHex(_userId)
	if err != nil {
		return nil, err
	}

	user := &models.User{}

	err = j.db.FindOne(
		ctx,
		user,
		bson.M{"_id": userId},
		bson.M{"hashed_password": 0},
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
