package providers

import "github.com/NikosSiak/Open-Banking-API/models"

type Bank interface {
	Name() string
	GetUserTokens(code string) (models.AccountTokens, error)
}
