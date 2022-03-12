package sms

import (
	"github.com/NikosSiak/Open-Banking-API/lib"
	"github.com/twilio/twilio-go/client"
	VerifyV2 "github.com/twilio/twilio-go/rest/verify/v2"
)

const channel = "sms"

type twilioProvider struct {
	verifySID    string
	verifyClient *VerifyV2.ApiService
}

func newTwilioProvider(credentials *lib.TwilioCredentials) twilioProvider {
	defaultClient := &client.Client{
		Credentials: client.NewCredentials(credentials.AccountSID, credentials.AuthToken),
	}

	verifyClient := VerifyV2.NewApiServiceWithClient(defaultClient)
	return twilioProvider{
		verifySID:    credentials.VerifySID,
		verifyClient: verifyClient,
	}
}

func (t twilioProvider) SendVerificationCode(to string) (string, error) {
	params := VerifyV2.CreateVerificationParams{}
	params.SetChannel(channel)
	params.SetTo(to)

	verification, err := t.verifyClient.CreateVerification(t.verifySID, &params)
	if err != nil {
		return "", err
	}

	return *verification.Sid, nil
}

func (t twilioProvider) VerifyCode(verificationId string, code string) (bool, error) {
	params := VerifyV2.CreateVerificationCheckParams{}
	params.SetCode(code)
	params.SetVerificationSid(verificationId)

	res, err := t.verifyClient.CreateVerificationCheck(t.verifySID, &params)

	return *res.Valid, err
}
