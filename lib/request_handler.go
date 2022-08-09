package lib

import (
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

type RequestHandler struct {
	Gin *gin.Engine
}

func NewRequestHandler(env Env) RequestHandler {
	engine := gin.New()
	engine.Use(CORSMiddleware)

	if env.IsProduction() {
		err := sentry.Init(sentry.ClientOptions{
			Dsn: env.SentryDsn,
			// Set TracesSampleRate to 1.0 to capture 100%
			// of transactions for performance monitoring.
			// We recommend adjusting this value in production,
			TracesSampleRate: 1.0,
		})

		if err == nil {
			engine.Use(sentrygin.New(sentrygin.Options{
				Repanic: true,
			}))
		}
	}

	return RequestHandler{Gin: engine}
}

func CORSMiddleware(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}
