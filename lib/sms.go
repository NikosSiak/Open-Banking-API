package lib

import (
	"github.com/twilio/twilio-go/client"
	VerifyV2 "github.com/twilio/twilio-go/rest/verify/v2"
)

const channel = "sms"

type SMSProvider struct {
	verifySID    string
	verifyClient *VerifyV2.ApiService
}

func NewSMSProvider(env Env) SMSProvider {
	defaultClient := &client.Client{
		Credentials: client.NewCredentials(env.TwilioCredentials.AccountSID, env.TwilioCredentials.AuthToken),
	}

	// defaultClient.SetAccountSid(env.TwilioCredentials.AccountSID)
	verifyClient := VerifyV2.NewApiServiceWithClient(defaultClient)
	return SMSProvider{
		verifySID:    env.TwilioCredentials.VerifySID,
		verifyClient: verifyClient,
	}
}

func (s SMSProvider) SendVerificationCode(to string) (string, error) {
	params := VerifyV2.CreateVerificationParams{}
	params.SetChannel(channel)
	params.SetTo(to)

	verification, err := s.verifyClient.CreateVerification(s.verifySID, &params)
	if err != nil {
		return "", err
	}

	return *verification.Sid, nil
}

func (s SMSProvider) VerifyCode(sid string, code string) (bool, error) {
	params := VerifyV2.CreateVerificationCheckParams{}
	params.SetCode(code)
	params.SetVerificationSid(sid)

	res, err := s.verifyClient.CreateVerificationCheck(s.verifySID, &params)

	return *res.Valid, err
}
