package lib

import (
  "fmt"

  "github.com/spf13/viper"
)

type ProviderCredentials struct {
  ClientId string `mapstructure:"client_id"`
  ClientSecret string `mapstructure:"client_secret"`
  SubscriptionKey string `mapstructure:"subscription_key"`
}

type Env struct {
  ServerPort string `mapstructure:"server_port"`
  Providers map[string]*ProviderCredentials `mapstructure:"providers"`
  DatabaseURI string `mapstructure:"db_uri"`
  DatabaseName string `mapstructure:"db_name"`
  Environment string `mapstructure:"environment"`
}

func GetEnv() Env {
  env := Env{
    ServerPort: "1312",
    Environment: "development",
  }

  viper.SetConfigFile("config.json")

  if err := viper.ReadInConfig(); err != nil {
    panic(fmt.Errorf("fatal error config file: %w", err))
  }

  err := viper.Unmarshal(&env)
  if err != nil {
    panic(fmt.Errorf("fatal error config file cant be loaded: %w", err))
  }

  return env
}
