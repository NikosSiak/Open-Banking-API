package routers_v1

import (
  "github.com/gin-gonic/gin"
  controllers_v1 "github.com/one-twenty/GloryDays/controllers/v1"
)

func SetBalanceRoutes(group *gin.RouterGroup, controller controllers_v1.BalancesController) {
  group.GET("/balances", controller.Get)
}
