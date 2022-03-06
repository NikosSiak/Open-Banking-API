package utils

import "github.com/gin-gonic/gin"

func NewError(ctx *gin.Context, status int, err error) {
	er := HTTPError{
		Message: err.Error(),
	}

	ctx.JSON(status, er)
}

type HTTPError struct {
	Message string `json:"message" example:"something went wrong"`
}
