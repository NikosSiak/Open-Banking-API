package controllers_v1

import (
  "github.com/one-twenty/GloryDays/services/providers"

  "github.com/gin-gonic/gin"
)


type BalancesController struct {
  providers []providers.Bank
}

func NewBalancesController(providers []providers.Bank) BalancesController {
  return BalancesController{ providers: providers }
}

func (controller BalancesController) Get(ctx *gin.Context) {
  var balance int64 = 0
  for _, provider := range controller.providers {
    balance += provider.GetBalance()
  }

  ctx.JSON(200, gin.H{
    "total_balance": balance,
  })
}
