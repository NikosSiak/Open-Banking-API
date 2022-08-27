package providers

import (
	"crypto/tls"
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

type Eurobank struct {
	appUrl                  string
	baseAuthUrl, baseApiUrl string
	clientId, clientSecret  string
}

const EurobankName = "eurobank"

func NewEurobankProvider(appUrl string, provider lib.ProviderCredentials) Eurobank {
	return Eurobank{
		appUrl:       appUrl,
		baseAuthUrl:  provider.BaseAuthUrl,
		baseApiUrl:   provider.BaseApiUrl,
		clientId:     provider.ClientId,
		clientSecret: provider.ClientSecret,
	}
}

func (e Eurobank) Name() string {
	return EurobankName
}

func (e Eurobank) LoginUri(userID string) string {
	redirectUri := e.appUrl + "/v1/accounts/eurobank/webhook"

	loginUri := fmt.Sprintf(
		"%s/authorize?scope=accounts&response_type=code&client_id=%s&redirect_uri=%s&state=%s",
		e.baseAuthUrl,
		e.clientId,
		redirectUri,
		userID,
	)

	return loginUri
}

func (e Eurobank) GetUserTokens(code string) (models.AccountTokens, error) {
	account := models.AccountTokens{}

	url := e.baseAuthUrl + "/token"
	method := "POST"

	redirectUri := e.appUrl + "/v1/accounts/eurobank/webhook"

	payload := strings.NewReader(
		fmt.Sprintf("grant_type=authorization_code&scope=accounts&code=%s&redirect_uri=%s",
			code,
			redirectUri,
		),
	)

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return account, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	auth := e.clientId + ":" + e.clientSecret
	authEnc := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", "Basic "+authEnc)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			Renegotiation: tls.RenegotiationSupport(tls.RequestClientCert),
		},
	}

	client := &http.Client{Transport: tr}
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
		Error        string `json:"error_description"`
		HttpError    string `json:"moreInformation"`
	}

	resp := &response{}
	err = json.Unmarshal(body, resp)

	if err == nil && resp.HttpError != "" {
		err = errors.New(resp.HttpError)
	}

	if err == nil && resp.Error != "" {
		err = errors.New(resp.Error)
	}

	account.AccessToken = resp.AccessToken
	account.RefreshToken = resp.RefreshToken

	return account, err
}
