package providers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/NikosSiak/Open-Banking-API/lib"
	"github.com/NikosSiak/Open-Banking-API/models"
)

type Piraeusbank struct {
	appUrl                  string
	baseAuthUrl, baseApiUrl string
	clientId, clientSecret  string
	scope                   string
}

const PiraeusbankName = "piraeusbank"

func NewPiraeusbankProvider(appUrl string, provider lib.ProviderCredentials) Piraeusbank {
	return Piraeusbank{
		appUrl:       appUrl,
		baseAuthUrl:  provider.BaseAuthUrl,
		baseApiUrl:   provider.BaseApiUrl,
		clientId:     provider.ClientId,
		clientSecret: provider.ClientSecret,
		scope:        provider.Scope,
	}
}

func (e Piraeusbank) Name() string {
	return PiraeusbankName
}

func (p Piraeusbank) LoginUri(userID string) string {
	redirectUri := p.appUrl + "/v1/accounts/piraeusbank/webhook"

	loginUri := fmt.Sprintf(
		"%s/authorize?scope=%s&response_type=code&client_id=%s&redirect_uri=%s&state=%s",
		p.baseAuthUrl,
		p.scope,
		p.clientId,
		redirectUri,
		userID,
	)

	return loginUri
}

func (p Piraeusbank) GetUserTokens(code string) (models.Account, error) {
	account := models.Account{}

	url := p.baseAuthUrl + "/token"
	method := "POST"

	payload := strings.NewReader(
		fmt.Sprintf("grant_type=authorization_code&scope=accounts&code=%s",
			code,
		),
	)

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return account, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	auth := p.clientId + ":" + p.clientSecret
	authEnc := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", "Basic "+authEnc)

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return account, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return account, err
	}

	type response struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		Error        string `json:"error"`
	}

	resp := &response{}
	err = json.Unmarshal(body, resp)
	if err == nil && resp.Error != "" {
		err = errors.New(resp.Error)
	}

	account.AccessToken = resp.AccessToken
	account.RefreshToken = resp.RefreshToken

	return account, err
}
