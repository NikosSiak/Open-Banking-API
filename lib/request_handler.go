package lib

import (
  "github.com/gin-gonic/gin"
)

type RequestHandler struct {
  Gin *gin.Engine
}

func NewRequestHandler() RequestHandler {
  engine := gin.New()
  engine.Use(CORSMiddleware)
  return RequestHandler{Gin: engine}
}

func CORSMiddleware(c *gin.Context) {
  c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
  c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
  c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
  c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

  if c.Request.Method == "OPTIONS" {
      c.AbortWithStatus(204)
      return
  }

  c.Next()
}
