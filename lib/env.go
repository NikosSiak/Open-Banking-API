package lib

import (
	"fmt"

	"github.com/spf13/viper"
)

type ProviderCredentials struct {
	BaseUrl         string `mapstructure:"base_url"`
	BaseApiUrl      string `mapstructure:"base_api_url"`
	ClientId        string `mapstructure:"client_id"`
	ClientSecret    string `mapstructure:"client_secret"`
	SubscriptionKey string `mapstructure:"subscription_key"`
}

type RedisCredentials struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
}

type TwilioCredentials struct {
	AccountSID string `mapstructure:"account_sid"`
	AuthToken  string `mapstructure:"auth_token"`
	VerifySID  string `mapstructure:"verify_sid"`
}

type Env struct {
	AppUrl            string                          `mapstructure:"app_url"`
	ServerPort        string                          `mapstructure:"server_port"`
	DatabaseURI       string                          `mapstructure:"db_uri"`
	DatabaseName      string                          `mapstructure:"db_name"`
	RedisCredentials  *RedisCredentials               `mapstructure:"redis"`
	JWTSecret         string                          `mapstructure:"jwt_secret"`
	Environment       string                          `mapstructure:"environment"`
	Providers         map[string]*ProviderCredentials `mapstructure:"providers"`
	TwilioCredentials *TwilioCredentials              `mapstructure:"twilio"`
}

func GetEnv() Env {
	env := Env{
		ServerPort:  "1312",
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

func (e Env) IsProduction() bool {
	return e.Environment == "production"
}
