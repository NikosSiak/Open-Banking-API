package main

import (
  controllers_v1 "github.com/one-twenty/GloryDays/controllers/v1"
  "github.com/one-twenty/GloryDays/routers/v1"
  "github.com/one-twenty/GloryDays/services/providers"

  "fmt"
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/spf13/viper"
)

var router *gin.Engine

func init() {
  viper.SetConfigFile("config.json")

  if err := viper.ReadInConfig(); err != nil {
    panic(fmt.Errorf("fatal error config file: %w", err))
  }

  router = gin.New()
  version1 := router.Group("/v1")

  providers := initProviders()
  balancesController := controllers_v1.NewBalancesController(providers)

  routers_v1.InitRoutes(version1, balancesController)
}

func main() {
  fmt.Println("Server Running on Port: ", 1312)
  http.ListenAndServe(":1312", router)
}

func initProviders() []providers.Bank {
  var bankProviders []providers.Bank

  type ProviderCredentials struct {
    ClientId string `mapstructure:"clientId"`
    ClientSecret string `mapstructure:"clientSecret"`
  }

  providersConfig := make(map[string]*ProviderCredentials)
  viper.UnmarshalKey("providers", &providersConfig)

  if len(providersConfig) == 0 {
    panic("No providers config")
  }

  if alphaCredentials, ok := providersConfig["alpha"]; ok {
    fmt.Println("Initializing Alpha bank provider")
    bankProviders = append(bankProviders,
                           providers.NewAlphaProvider(alphaCredentials.ClientId, alphaCredentials.ClientSecret))
  }

  return bankProviders
}
