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

type Alpha struct {
	appUrl                                  string
	baseUrl, baseAPIUrl                     string
	clientId, clientSecret, subscriptionKey string
}

const AlphaName = "alpha"

func NewAlphaProvider(appUrl string, provider lib.ProviderCredentials) Alpha {
	return Alpha{
		appUrl:          appUrl,
		baseUrl:         provider.BaseUrl,
		baseAPIUrl:      provider.BaseApiUrl,
		clientId:        provider.ClientId,
		clientSecret:    provider.ClientSecret,
		subscriptionKey: provider.SubscriptionKey,
	}
}

func (a Alpha) Name() string {
	return AlphaName
}

func (a Alpha) LoginUri(userID string) (string, error) {
	accessToken, err := a.getClientAccessToken()
	if err != nil {
		return "", err
	}

	accountRequest, err := a.createAccountRequest(accessToken)
	if err != nil {
		return "", err
	}

	redirectUri := a.appUrl + "/v1/accounts/alpha/webhook/"
	loginUri := fmt.Sprintf(
		"%s/auth/authorize?client_id=%s&response_type=code&scope=account-info&redirect_uri=%s&request=%s&state=%s",
		a.baseUrl,
		a.clientId,
		redirectUri,
		accountRequest,
		userID,
	)

	return loginUri, nil
}

func (a Alpha) GetUserTokens(code string) (models.Account, error) {
	account := models.Account{}

	url := a.baseUrl + "/auth/token"
	method := "POST"

	redirectUri := a.appUrl + "/v1/accounts/alpha/webhook/"
	payload := strings.NewReader(
		fmt.Sprintf("grant_type=authorization_code&code=%s&client_id=%s&client_secret=%s&redirect_uri=%s",
			code,
			a.clientId,
			a.clientSecret,
			redirectUri,
		),
	)

	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return account, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencode")

	client := &http.Client{}
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
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
	}

	resp := &response{}
	err = json.Unmarshal(body, resp)
	if err == nil && resp.Error != "" {
		err = errors.New(resp.Error)
	}

	account.AccessToken = resp.AccessToken

	return account, err
}

func (a Alpha) getClientAccessToken() (string, error) {
	url := a.baseUrl + "/auth/token"
	method := "POST"

	payload := strings.NewReader("grant_type=client_credentials&scope=account-info-setup")

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return "", err
	}

	auth := a.clientId + ":" + a.clientSecret
	authEnc := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", "Basic "+authEnc)

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	type response struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
	}

	resp := &response{}
	err = json.Unmarshal(body, resp)
	if err == nil && resp.Error != "" {
		err = errors.New(resp.Error)
	}

	return resp.AccessToken, err
}

func (a Alpha) createAccountRequest(accessToken string) (string, error) {
	url := a.baseAPIUrl + "/accounts/v1/account-requests"
	method := "POST"

	payload := strings.NewReader("{\"Risk\": {},\"ProductIdentifiers\":null}")

	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Ocp-Apim-Subscription-Key", a.subscriptionKey)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	type response struct {
		AccountRequestId string `json:"AccountRequestId"`
		Error            string `json:"Description"`
	}

	resp := &response{}
	err = json.Unmarshal(body, resp)
	if err == nil && resp.Error != "" {
		err = errors.New(resp.Error)
	}

	return resp.AccountRequestId, err
}
