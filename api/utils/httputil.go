package utils

import (
	"fmt"
	"runtime"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

func NewError(ctx *gin.Context, status int, err error) {
	er := HTTPError{
		Message: err.Error(),
	}

	if hub := sentrygin.GetHubFromContext(ctx); hub != nil {
		_, file, line, ok := runtime.Caller(1)

		hub.WithScope(func(scope *sentry.Scope) {
			if ok {
				scope.SetExtra("source", fmt.Sprintf("%s#%d", file, line))
			}
			hub.CaptureMessage(er.Message)
		})
	}

	ctx.JSON(status, er)
}

type HTTPError struct {
	Message string `json:"message" example:"something went wrong"`
}
