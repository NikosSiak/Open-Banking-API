package routers_v1

import (
	"github.com/gin-gonic/gin"
	controllers_v1 "github.com/one-twenty/GloryDays/controllers/v1"
)

func InitRoutes(group *gin.RouterGroup, balancesController controllers_v1.BalancesController) {
  SetBalanceRoutes(group, balancesController)
}
